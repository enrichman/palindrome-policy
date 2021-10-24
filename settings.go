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
	Threshold         int        `json:"threshold"`
}

// Builds a new Settings instance starting from a validation
// request payload:
// {
//    "request": ...,
//    "settings": {
//       "whitelisted_labels": [...],
//       "threshold": 3
//    }
// }
func NewSettingsFromValidationReq(payload []byte) (Settings, error) {
	return newSettings(
		payload,
		"settings.whitelisted_labels",
		"settings.threshold")
}

// Builds a new Settings instance starting from a Settings
// payload:
// {
//    "whitelisted_labels": [...],
//    "threshold": 3
// }
func NewSettingsFromValidateSettingsPayload(payload []byte) (Settings, error) {
	return newSettings(
		payload,
		"whitelisted_labels",
		"threshold")
}

func newSettings(payload []byte, paths ...string) (Settings, error) {
	data := gjson.GetManyBytes(payload, paths...)

	whitelistedLabels := mapset.NewThreadUnsafeSet()
	data[0].ForEach(func(_, entry gjson.Result) bool {
		whitelistedLabels.Add(entry.String())
		return true
	})

	threshold := data[1].Int()

	return Settings{
		WhitelistedLabels: whitelistedLabels,
		Threshold:         int(threshold),
	}, nil
}

// Valid check if the received settings are valid.
// WhitelistedLabels cannot contains any not palindrome and the Threshold field cannot be negative
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
			strings.Join(notPalindromeLabels, ","))

		return false, errors.New(errMsg)
	}

	if s.Threshold < 0 {
		errMsg := fmt.Sprintf(
			"Threshold cannot be negative: %d",
			s.Threshold)

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
