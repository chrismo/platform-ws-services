package main

import "testing"

// TODO: log calls clutter test output

func TestPagerdutyIntegrationTrigger(t *testing.T) {
	pd := Pagerduty{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    1,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
		Settings: Settings{
			PagerdutyKey: "545bf39a778d45b1b4160a7fd782fae9", // free, trial account
		},
	}
	res := pd.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "PagerDuty trigger event posted.", res.Message)
}

func TestPagerdutyIntegrationResolve(t *testing.T) {
	pd := Pagerduty{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    0,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
		Settings: Settings{
			PagerdutyKey: "545bf39a778d45b1b4160a7fd782fae9", // free, trial account
		},
	}
	res := pd.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "PagerDuty resolve event posted.", res.Message)
}

func TestPagerdutyIntegrationNoSetting(t *testing.T) {
	pd := Pagerduty{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    0,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
	}
	res := pd.Transmit(ap)
	equals(t, Skipped, res.Result)
	equals(t, "Deployment/Group has no PagerDuty Setting configured.", res.Message)
}

func TestPagerdutyIntegrationBadSetting(t *testing.T) {
	/* Huh, PagerDuty seems to take any 32 char string as a key and doesn't
	   return any error, I guess as a security precaution, to not divulge
	   difference between good keys and bad keys. Trade-off is: make sure
	   your key is correct, otherwise you may think you're sending PD notices
	   when you're not. */
	pd := Pagerduty{}
	ap := AlertPackage{
		Alert: Alert{
			CapsuleID: "capsule-id",
			Status:    0,
		},
		Check: Check{
			Description: "lookin' for bad stuff",
		},
		Settings: Settings{
			PagerdutyKey: "00000000000000000000000000000000", // probably not a real account
		},
	}
	res := pd.Transmit(ap)
	equals(t, Success, res.Result)
	equals(t, "PagerDuty resolve event posted.", res.Message)
}
