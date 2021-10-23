package main

import (
	"testing"
)

func TestParsingSettingsWithAllValuesProvidedFromValidationReq(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {
			"whitelisted_labels": [ "level", "radar" ]
		}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	expected := []string{"level", "radar"}
	for _, exp := range expected {
		if !settings.WhitelistedLabels.Contains(exp) {
			t.Errorf("Missing value %s", exp)
		}
	}
}

func TestParsingSettingsWithNotPalindromeLabelsAreNotValid(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {
			"whitelisted_labels": [ "foo", "bar" ]
		}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	valid, err := settings.Valid()
	if valid {
		t.Errorf("Settings are reported as valid")
	}
	if err == nil {
		t.Errorf("Unexpected missing error")
	}
}

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if settings.WhitelistedLabels.Cardinality() != 0 {
		t.Errorf("Expecpted WhitelistedLabels to be empty")
	}
}

func TestSettingsAreValid(t *testing.T) {
	request := `
	{
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidateSettingsPayload(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}
