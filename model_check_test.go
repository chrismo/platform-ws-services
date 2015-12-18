package main

import "testing"

func TestIdempotentCheck(t *testing.T) {
	// ran into this building out integration.rb tests.
	setupTestDB()
	defer tearDownTestDB()
	c := Check{Type: "type", Name: "name", Level: 2}
	equals(t, nil, c.Save())
	equals(t, nil, c.Save())
}
