//go:build unit
// +build unit

package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringInSlice(t *testing.T) {
	t.Run("Not Found", func(t *testing.T) {
		assert.False(t, IsStringInSlice("foo", []string{"bar", "whatever"}))
	})

	t.Run("Found", func(t *testing.T) {
		assert.True(t, IsStringInSlice("foo", []string{"bar", "foo"}))
	})
}
