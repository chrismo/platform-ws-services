package main

import "testing"

// TODO: log calls clutter test output

const (
	//	SlackTestAPIKey = "xoxb-16976869585-tXlGTDPPUJll00rZ0XZ3hSTf"
	SlackTestAPIKey = "xoxp-16973770979-16973770995-16977328082-d9933fe3d4"
)

func TestSlackIntegrationTrigger(t *testing.T) {
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
				ApiKey:  SlackTestAPIKey,
				Channel: "#general",
			},
		},
	}
	res := s.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "Triggered alert posted to Slack.", res.Message)
}

func TestSlackIntegrationResolve(t *testing.T) {
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
				ApiKey:  SlackTestAPIKey,
				Channel: "#general",
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
				ApiKey:  "not-really-a-key",
				Channel: "#general",
			},
		},
	}
	res := s.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "Resolved alert posted to Slack.", res.Message)
}
