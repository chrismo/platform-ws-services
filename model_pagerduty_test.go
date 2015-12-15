package main

import "testing"

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
	}
	pd.Transmit(ap)
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
	}
	pd.Transmit(ap)
}
