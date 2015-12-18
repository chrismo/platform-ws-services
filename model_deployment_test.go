package main

import "testing"

func TestIdempotentDeployment(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()
	c := Deployment{Id: "123", GroupId: "g123", Type: "type", Name: "name"}
	equals(t, nil, c.Save())
	equals(t, nil, c.Save())
}
