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

func TestCurrentChecksMerge(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()
	checkA := Check{Type: "urSql", Name: "storage", Title: "deployment-specific"}
	d := Deployment{Id: "456", GroupId: "g456", Type: "urSql", Name: "bob's database",
		Checks: []Check{checkA},
	}
	equals(t, nil, d.Save())

	checkA.Title = "default"
	equals(t, nil, checkA.Save())

	checkB := Check{Type: "urSql", Name: "cpu", Title: "default"}
	equals(t, nil, checkB.Save())

	d, _ = LookupDeploymentById(d.Id)
	equals(t, "urSql", d.Type)
	equals(t, 1, len(d.Checks))
	equals(t, "deployment-specific", d.Checks[0].Title)

	checks, err := d.DefaultChecks()
	equals(t, nil, err)
	equals(t, 2, len(checks))
	equals(t, "default", checks[0].Title)
	equals(t, "default", checks[1].Title)

	checks, err = d.CurrentChecks()
	equals(t, nil, err)
	equals(t, 2, len(checks))
	equals(t, "deployment-specific", checks[0].Title)
	equals(t, "default", checks[1].Title)

	storageCheck, _ := d.CheckByName("storage")
	equals(t, "deployment-specific", storageCheck.Title)

	cpuCheck, _ := d.CheckByName("cpu")
	equals(t, "default", cpuCheck.Title)
}
