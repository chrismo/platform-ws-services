package main

import (
	"log"
	"testing"
)

func TestIdempotentDeployment(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()
	d := Deployment{Id: "123", GroupId: "g123", Type: "type", Name: "name"}
	equals(t, nil, d.Save())
	equals(t, nil, d.Save())
}

func TestCurrentChecksWithNoChecks(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()
	d := Deployment{Id: "123", GroupId: "g123", Type: "type", Name: "name"}
	_, err := d.CurrentChecks()
	equals(t, nil, err)
}

func TestCurrentChecksWithOneCheck(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()
	c := Check{Type: "type", Name: "name"}
	equals(t, nil, c.Save())
	d := Deployment{Id: "123", GroupId: "g123", Type: "type", Name: "name",
		Checks: []Check{c},
	}
	_, err := d.CurrentChecks()
	log.Print(err)
	equals(t, nil, err)
}
