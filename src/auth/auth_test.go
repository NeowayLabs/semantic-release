//go:build unit
// +build unit

package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicKey(t *testing.T) {
	sshKey := os.Getenv("SSH_INTEGRATION_SEMANTIC")
	sshKey = FormatSSHKey(sshKey, "#")
	key, err := GetPublicKey([]byte(sshKey))
	assert.NotNil(t, key)
	assert.NoError(t, err)
}

func TestFormatSSHKey(t *testing.T) {
	expected := "hi\nhello"
	assert.EqualValues(t, expected, FormatSSHKey("hi#hello", "#"))
}
