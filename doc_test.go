package pty

import (
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/term"
)

// Will fill p from reader r
func readBytes(r io.Reader, p []byte) error {
	bytesRead := 0
	for bytesRead < len(p) {
		n, err := r.Read(p[bytesRead:])
		if err != nil {
			return err
		}
		bytesRead = bytesRead + n
	}
	return nil
}

func TestOpen(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestIsTerminal(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	assert.True(t, term.IsTerminal(int(pty.Fd())), "pty is not a terminal")
	assert.True(t, term.IsTerminal(int(tty.Fd())), "tty is not a terminal")

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestName(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	assert.Equal(t, "/dev/ptmx", pty.Name())
	assert.Regexp(t, regexp.MustCompile(`/dev/pts/\d+`), tty.Name())

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestGetsize(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	prows, pcols, err := Getsize(pty)
	assert.NoError(t, err)

	trows, tcols, err := Getsize(tty)
	assert.NoError(t, err)

	assert.Equal(t, prows, trows)
	assert.Equal(t, pcols, tcols)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestGetsizefull(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	psize, err := GetsizeFull(pty)
	assert.NoError(t, err)

	tsize, err := GetsizeFull(tty)
	assert.NoError(t, err)

	assert.Equal(t, psize.X, tsize.X)
	assert.Equal(t, psize.Y, tsize.Y)
	assert.Equal(t, psize.Rows, tsize.Rows)
	assert.Equal(t, psize.Cols, tsize.Cols)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestSetsize(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	psize, err := GetsizeFull(pty)
	assert.NoError(t, err)

	psize.X = psize.X + 1
	psize.Y = psize.Y + 1
	psize.Rows = psize.Rows + 1
	psize.Cols = psize.Cols + 1

	Setsize(tty, psize)

	tsize, err := GetsizeFull(tty)
	assert.NoError(t, err)

	assert.Equal(t, psize.X, tsize.X)
	assert.Equal(t, psize.Y, tsize.Y)
	assert.Equal(t, psize.Rows, tsize.Rows)
	assert.Equal(t, psize.Cols, tsize.Cols)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestReadWriteText(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	// Write to tty, read from pty
	text := []byte("ping")
	n, err := tty.Write(text)
	assert.NoError(t, err)
	assert.Equal(t, len(text), n)

	buffer := make([]byte, 4)
	err = readBytes(pty, buffer)
	assert.NoError(t, err)
	assert.Equal(t, text, buffer)

	// Write to pty, read from tty.
	// We need to send a \n otherwise this will block in the terminal driver.
	text = []byte("pong\n")
	n, err = pty.Write(text)
	assert.NoError(t, err)
	assert.Equal(t, len(text), n)

	buffer = make([]byte, 5)
	err = readBytes(tty, buffer)
	assert.NoError(t, err)
	assert.Equal(t, []byte("pong\n"), buffer)

	// Read the echo back from pty
	buffer = make([]byte, 5)
	err = readBytes(pty, buffer)
	assert.NoError(t, err)
	assert.Equal(t, []byte("pong\r"), buffer)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}

func TestReadWriteControls(t *testing.T) {
	t.Parallel()

	pty, tty, err := Open()
	assert.NoError(t, err)

	// Write the start of a line to pty
	text := []byte("pind")
	n, err := pty.Write(text)
	assert.NoError(t, err)
	assert.Equal(t, len(text), n)

	// Backspace that last char
	n, err = pty.Write([]byte("\b"))
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	// Write the correct char and a CR
	n, err = pty.Write([]byte("g\n"))
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	// Read the line
	buffer := make([]byte, 7)
	err = readBytes(tty, buffer)
	assert.NoError(t, err)
	assert.Equal(t, []byte("pind\bg\n"), buffer)

	// Read the echo back from pty
	buffer = make([]byte, 7)
	err = readBytes(pty, buffer)
	assert.NoError(t, err)
	assert.Equal(t, []byte("pind^Hg"), buffer)

	err = tty.Close()
	assert.NoError(t, err)

	err = pty.Close()
	assert.NoError(t, err)
}
