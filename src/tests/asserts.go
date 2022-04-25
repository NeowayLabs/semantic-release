package tests

import "testing"

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error("Should not return an error")
	}
}

func AssertError(t *testing.T, err error) {
	if err == nil {
		t.Error("Should return an error")
	}
}
