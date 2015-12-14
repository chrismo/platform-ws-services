package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestAlertModelJsonRoundtrip(t *testing.T) {
	input :=
		ClientAlert{Client: "foobar",
			Alert: Alert{
				Name:         "name",
				CapsuleName:  "capsuleName",
				Output:       "output",
				Status:       0.33333333333333, // float64?
				CapsuleID:    "capID",
				DeploymentID: "depID",
				AccountSlug:  "acc",
			},
		}
	jsonString, _ := json.Marshal(input)
	alert, _ := NewAlertFromJSON(bytes.NewReader(jsonString))
	if alert.Name != "name" {
		t.Error("name did not roundtrip")
	}

	if alert.CapsuleName != "capsuleName" {
		t.Error("capsuleName did not roundtrip")
	}

	if alert.Output != "output" {
		t.Error("output did not roundtrip")
	}

	if alert.Status != 0.33333333333333 {
		t.Error("status did not roundtrip")
	}

	if alert.CapsuleID != "capID" {
		t.Error("capsuleID did not roundtrip")
	}

	if alert.DeploymentID != "depID" {
		t.Error("deploymentID did not roundtrip")
	}

	if alert.AccountSlug != "acc" {
		t.Error("accountSlug did not roundtrip")
	}
}

func TestAlertModelFromJsonBadJson(t *testing.T) {
	_, err := NewAlertFromJSON(strings.NewReader(":)"))
	if err == nil {
		t.Error("should have errored")
	}
	if err.Error() != "invalid character ':' looking for beginning of value" {
		t.Error(err)
	}
}
