package main

import (
	"errors"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"

	"fmt"
)

type Settings struct {
	WhitelistedLabels mapset.Set `json:"whitelisted_labels"`
}

// Builds a new Settings instance starting from a validation
// request payload:
// {
//    "request": ...,
//    "settings": {
//       "whitelisted_labels": [...]
//    }
// }
func NewSettingsFromValidationReq(payload []byte) (Settings, error) {
	return newSettings(
		payload,
		"settings.whitelisted_labels")
}

// Builds a new Settings instance starting from a Settings
// payload:
// {
//    "whitelisted_labels": ...
// }
func NewSettingsFromValidateSettingsPayload(payload []byte) (Settings, error) {
	return newSettings(
		payload,
		"whitelisted_labels")
}

func newSettings(payload []byte, whitelistedLabelsPath string) (Settings, error) {
	data := gjson.GetBytes(payload, whitelistedLabelsPath)

	whitelistedLabels := mapset.NewThreadUnsafeSet()
	data.ForEach(func(_, entry gjson.Result) bool {
		whitelistedLabels.Add(entry.String())
		return true
	})

	return Settings{
		WhitelistedLabels: whitelistedLabels,
	}, nil
}

// No special check has to be done
func (s *Settings) Valid() (bool, error) {
	notPalindromeLabels := []string{}

	for _, label := range s.WhitelistedLabels.ToSlice() {
		if !isPalindrome(label.(string)) {
			notPalindromeLabels = append(notPalindromeLabels, label.(string))
		}
	}

	if len(notPalindromeLabels) > 0 {
		errMsg := fmt.Sprintf(
			"The following whitelisted labels are not palindromes: %s",
			strings.Join(notPalindromeLabels, ","),
		)

		return false, errors.New(errMsg)
	}

	return true, nil
}

func validateSettings(payload []byte) ([]byte, error) {
	logger.Info("validating settings")

	settings, err := NewSettingsFromValidateSettingsPayload(payload)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(err.Error()))
	}

	valid, err := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	logger.Warn("rejecting settings")
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}
