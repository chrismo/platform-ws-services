package main

import (
	"io/ioutil"
	"log"
	"testing"
)

// TODO: re-use setup code

const (
	WebHookTestURL = "https://hooks.slack.com/services/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"
)

func TestSlackIntegrationTrigger(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	s := Slack{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    1,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
		Settings: Settings{
			Slack: Slack{
				WebHookURL: WebHookTestURL,
			},
		},
	}
	res := s.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "Triggered alert posted to Slack.", res.Message)
}

func TestSlackIntegrationResolve(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	s := Slack{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    0,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
		Settings: Settings{
			Slack: Slack{
				WebHookURL: WebHookTestURL,
			},
		},
	}
	res := s.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "Resolved alert posted to Slack.", res.Message)
}

func TestSlackIntegrationNoSetting(t *testing.T) {
	s := Slack{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    0,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
	}
	res := s.Transmit(ap)
	equals(t, Skipped, res.Result)
	equals(t, "Deployment/Group has no Slack Setting configured.", res.Message)
}

func TestSlackIntegrationBadSetting(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	s := Slack{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    0,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
		Settings: Settings{
			Slack: Slack{
				WebHookURL: WebHookTestURL,
			},
		},
	}
	res := s.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "Resolved alert posted to Slack.", res.Message)
}
