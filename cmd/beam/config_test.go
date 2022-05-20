package main

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	testAddress := ":5050"
	os.Setenv("ADDRESS", testAddress)
	c := configure()
	if c.Address != testAddress {
		t.Logf("Failed Address Check: Expected %s, but got %s", testAddress, c.Address)
		t.FailNow()
	}
}
