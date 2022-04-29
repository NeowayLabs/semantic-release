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

func AssertEqualValues(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Not equal: \n"+
			"expected: %v\n"+
			"actual  : %v", expected, actual)
	}
}
