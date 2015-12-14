// +build seed

package main

import (
	"log"
	"testing"
)

func TestAlerterSetup(t *testing.T) {
	seedAlerts(t)
}

func seedAlerts(t *testing.T) {
	alerter, err := NewAlerter(*redisUrl, *redisPassword)
	if err != nil {
		log.Fatalf("unable to connect to Redis, %s\n", err.Error())
	}
	alerter.pool.Get().Flush()
	defer alerter.pool.Close()
	result := &Alert{
		Name:         "redis0-redis_role",
		CapsuleName:  "redis0",
		Output:       "expected role master, found role slave",
		Status:       0,
		CapsuleID:    "111111",
		DeploymentID: "987654321",
		AccountSlug:  "compose-test",
	}
	alerter.processAlert(result)
}
