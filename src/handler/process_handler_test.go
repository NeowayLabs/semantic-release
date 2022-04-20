//go:build unit
// +build unit

package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsts(t *testing.T) {
	assert.Equal(t, "CHANGELOG.md", changeLogDefaultFile)
	assert.Equal(t, "setup.py", setupPythonDefaultFile)
	assert.Equal(t, "HOME", homeEnv)
}
