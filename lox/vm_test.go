package lox

import "testing"

func TestVM(t *testing.T) {
	var c Chunk
	constant := c.addConstant(1.2)
	c.write(OP_CONSTANT, 123)
	c.write(byte(constant), 123)

	constant = c.addConstant(3.4)
	c.write(OP_CONSTANT, 123)
	c.write(byte(constant), 123)

	c.write(OP_ADD, 123)

	constant = c.addConstant(5.6)
	c.write(OP_CONSTANT, 123)
	c.write(byte(constant), 123)

	c.write(OP_DIVIDE, 123)
	c.write(OP_NEGATE, 123)
	c.write(OP_RETURN, 123)

	vm := NewVM(&c)
	vm.DebugTraceExecution = true
	_ = vm.Run()
}
