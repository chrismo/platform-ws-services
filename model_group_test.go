package main

import "testing"

func TestIdempotentGroup(t *testing.T) {
	// ran into this building out integration.rb tests.
	setupTestDB()
	defer tearDownTestDB()
	c := Group{Id: "3"}
	equals(t, nil, c.Save())
	equals(t, nil, c.Save())
}
