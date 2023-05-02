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

func AssertNotEmpty(t *testing.T, actual interface{}) {
	t.Helper()
	actualType := reflect.TypeOf(actual).Name()

	switch actualType {
	case "string":
		if "" == actual {
			t.Error("String value should not be empty")
		}
	}
}

func AssertTrue(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Errorf("Should be true")
	}
}

func AssertFalse(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Errorf("Should be false")
	}
}

func AssertNil(t *testing.T, actual interface{}) {
	t.Helper()
	if actual != nil {
		t.Errorf("Expected nil but was %v", actual)
	}
}

func AssertNotNil(t *testing.T, actual interface{}) {
	t.Helper()
	if actual == nil {
		t.Error("Expected not to be nil")
	}
}
