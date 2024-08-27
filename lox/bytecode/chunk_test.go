package bytecode

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	var c Chunk

	constant := c.addConstant(1.2)
	c.write(OP_CONSTANT, 123)
	c.write(byte(constant), 123)

	c.write(OP_RETURN, 123)
	output := captureOutput(func() {
		c.disassemble("test chunk")
	})
	assert.Equal(t, `== test chunk ==
0000  123 OP_CONSTANT         0 '1.2'
0002    | OP_RETURN
`, output)
}

func captureOutput(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}

	r.Close()

	return buf.String()
}
