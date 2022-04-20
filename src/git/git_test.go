//go:build unit
// +build unit

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	gitServiceUnit = New()
)

func TestDefaultSignature(t *testing.T) {

	expected := Signature{
		Name:  "Administrator",
		Email: "amd@example.com",
	}
	actualSignature := setDefaultSignature("Administrator", "amd@example.com")
	assert.NotNil(t, actualSignature)
	assert.EqualValues(t, expected.Name, actualSignature.Name)
	assert.EqualValues(t, expected.Email, actualSignature.Email)
}
