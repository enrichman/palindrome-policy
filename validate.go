package main

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"
)

func validate(payload []byte) ([]byte, error) {
	if !gjson.ValidBytes(payload) {
		return kubewarden.RejectRequest(
			kubewarden.Message("Not a valid JSON document"),
			kubewarden.Code(400))
	}

	settings, err := NewSettingsFromValidationReq(payload)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	data := gjson.GetBytes(
		payload,
		"request.object.metadata.labels")

	palindromeLabels := mapset.NewThreadUnsafeSet()

	data.ForEach(func(key, value gjson.Result) bool {
		label := key.String()

		if isPalindrome(label) {
			palindromeLabels.Add(label)
		}
		return true
	})

	error_msgs := []string{}

	notWhitelistedPalindromeLabels := palindromeLabels.Difference(settings.WhitelistedLabels)
	if notWhitelistedPalindromeLabels.Cardinality() > settings.Threshold {
		palindromes := []string{}
		for _, v := range notWhitelistedPalindromeLabels.ToSlice() {
			palindromes = append(palindromes, v.(string))
		}

		error_msgs = append(
			error_msgs,
			fmt.Sprintf(
				"Too many palindrome labels that are not-whitelisted: %s. Max allowed [%d]",
				strings.Join(palindromes, ","),
				settings.Threshold,
			))
	}

	if len(error_msgs) > 0 {
		return kubewarden.RejectRequest(
			kubewarden.Message(strings.Join(error_msgs, ". ")),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}

func isPalindrome(label string) bool {
	for head := 0; head < len(label)/2; head++ {
		tail := len(label) - head - 1

		if label[head] != label[tail] {
			return false
		}
	}
	return true
}
