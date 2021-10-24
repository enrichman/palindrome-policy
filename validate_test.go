package main

import (
	"encoding/json"
	"testing"

	mapset "github.com/deckarep/golang-set"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
)

func Test_validate_EmptySettingsValidPodApproval(t *testing.T) {
	settings := Settings{}

	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/pod.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func Test_validate_EmptySettingsPalindromePodRejection(t *testing.T) {
	settings := Settings{}

	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/pod-palindrome.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != false {
		t.Error("Unexpected approval")
	}
}

func Test_validate_PalindromePodWhitelistedApproval(t *testing.T) {
	settings := Settings{
		WhitelistedLabels: mapset.NewSetFromSlice([]interface{}{"level", "radar"}),
	}

	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/pod-palindrome.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	t.Log(response.Message)

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func Test_validate_PalindromePodThresholdApproval(t *testing.T) {
	settings := Settings{
		Threshold: 2,
	}

	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/pod-palindrome.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	t.Log(response.Message)

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func Test_validate_PalindromePodWhitelistedAndThresholdApproval(t *testing.T) {
	settings := Settings{
		WhitelistedLabels: mapset.NewSetFromSlice([]interface{}{"level"}),
		Threshold:         1,
	}

	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/pod-palindrome.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	t.Log(response.Message)

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func Test_isPalindrome(t *testing.T) {
	type args struct {
		label string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{label: ""},
			want: true,
		},
		{
			name: "single word",
			args: args{label: "a"},
			want: true,
		},
		{
			name: "no palindrome",
			args: args{label: "foobar"},
			want: false,
		},
		{
			name: "palindrome odd",
			args: args{label: "level"},
			want: true,
		},
		{
			name: "palindrome even",
			args: args{label: "raddar"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPalindrome(tt.args.label); got != tt.want {
				t.Errorf("isPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
