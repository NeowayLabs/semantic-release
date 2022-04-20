package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVars(t *testing.T) {
	assert.EqualValues(t, "\033[0m", Reset)
	assert.EqualValues(t, "\033[31m", Red)
	assert.EqualValues(t, "\033[32m", Green)
	assert.EqualValues(t, "\033[33m", Yellow)
	assert.EqualValues(t, "\033[34m", Blue)
	assert.EqualValues(t, "\033[35m", Purple)
	assert.EqualValues(t, "\033[36m", Cyan)
	assert.EqualValues(t, "\033[37m", Gray)
	assert.EqualValues(t, "\033[97m", White)
	assert.EqualValues(t, "\033[40;1;37m", BGBlack)
	assert.EqualValues(t, "\033[41;1;37m", BGRed)
	assert.EqualValues(t, "\033[42;1;37m", BGGreen)
	assert.EqualValues(t, "\033[44;1;37m", BGBlue)
	assert.EqualValues(t, "\033[45;1;37m", BGPurple)
	assert.EqualValues(t, "\033[46;1;37m", BGCyan)
	assert.EqualValues(t, "\033[47;1;37m", BGGray)
	assert.EqualValues(t, "\033[47m", BGWhite)
	assert.EqualValues(t, "\033[5;30m", Blink)
}
