package main

import "testing"

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
	d := Deployment{Id: "234", GroupId: "g234", Type: "type", Name: "name",
		Checks: []Check{c},
	}
	checks, err := d.CurrentChecks()
	equals(t, nil, err)
	equals(t, 1, len(checks))
}

func TestCurrentChecksWithManyChecks(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()
	checkA := Check{Type: "type", Name: "nameA"}
	checkB := Check{Type: "type", Name: "nameB"}
	equals(t, nil, checkA.Save())
	equals(t, nil, checkB.Save())
	d := Deployment{Id: "345", GroupId: "g345", Type: "type", Name: "name",
		Checks: []Check{checkA, checkB},
	}
	checks, err := d.CurrentChecks()
	equals(t, nil, err)
	equals(t, 2, len(checks))
}
