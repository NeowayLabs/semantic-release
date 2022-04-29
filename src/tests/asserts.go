package tests

import (
	"reflect"
	"testing"
)

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error("Should not return an error")
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Should return an error")
	}
}

func AssertEqualValues(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Not equal: \n"+
			"expected: %v\n"+
			"actual  : %v", expected, actual)
	}
}

func AssertEmpty(t *testing.T, actual interface{}) {
	t.Helper()
	actualType := reflect.TypeOf(actual).Name()

	switch actualType {
	case "string":
		if "" != actual {
			t.Error("Should return an empty string")
		}
	}
}
