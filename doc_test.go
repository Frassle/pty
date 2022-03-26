package pty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {

	pty, tty, err := Open()
	assert.NoError(t, err)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}
