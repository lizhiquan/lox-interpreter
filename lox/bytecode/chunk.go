package bytecode

import "fmt"

const (
	OP_CONSTANT byte = iota
	OP_RETURN
)

type Value float64

type Chunk struct {
	code      []byte
	lines     []int
	constants []Value
}

func (c *Chunk) write(byte byte, line int) {
	c.code = append(c.code, byte)
	c.lines = append(c.lines, line)
}

func (c *Chunk) addConstant(value Value) int {
	c.constants = append(c.constants, value)
	return len(c.constants) - 1
}

func (c *Chunk) disassemble(name string) {
	fmt.Printf("== %s ==\n", name)

	offset := 0
	for offset < len(c.code) {
		offset = c.disassembleInstruction(offset)
	}
}

func (c *Chunk) disassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && c.lines[offset] == c.lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", c.lines[offset])
	}

	instruction := c.code[offset]
	switch instruction {
	case OP_CONSTANT:
		return c.constantInstruction("OP_CONSTANT", offset)
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func (c *Chunk) constantInstruction(name string, offset int) int {
	constant := c.code[offset+1]
	fmt.Printf("%-16s %4d '%v'\n", name, constant, c.constants[constant])
	return offset + 2
}

func simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}
