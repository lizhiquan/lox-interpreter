package lox

import (
	"errors"
	"fmt"
)

type VM struct {
	chunk *Chunk
	// instruction pointer points to the next instruction to be executed
	ip int

	stack    [256]Value
	stackTop int

	DebugTraceExecution bool
}

var (
	ErrInterpretCompile = errors.New("interpret: compile error")
	ErrInterpretRuntime = errors.New("interpret: runtime error")
)

func NewVM(chunk *Chunk) *VM {
	return &VM{
		chunk:               chunk,
		ip:                  0,
		DebugTraceExecution: false,
		stackTop:            0,
	}
}

func Interpret(source string) error {
	var chunk Chunk
	if !compile(source, &chunk) {
		return ErrInterpretCompile
	}

	vm := NewVM(&chunk)
	return vm.run()
}

func (vm *VM) run() error {
	for {
		if vm.DebugTraceExecution {
			fmt.Print("          ")
			for i := 0; i < vm.stackTop; i++ {
				fmt.Printf("[ %v ]", vm.stack[i])
			}
			fmt.Println()
			vm.chunk.disassembleInstruction(vm.ip)
		}

		instruction := vm.readByte()
		switch instruction {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)

		case OP_ADD:
			vm.binaryOp("+")

		case OP_SUBTRACT:
			vm.binaryOp("-")

		case OP_MULTIPLY:
			vm.binaryOp("*")

		case OP_DIVIDE:
			vm.binaryOp("/")

		case OP_NEGATE:
			vm.push(-vm.pop())

		case OP_RETURN:
			fmt.Println(vm.pop())
			return nil
		}
	}
}

func (vm *VM) readByte() byte {
	b := vm.chunk.code[vm.ip]
	vm.ip++
	return b
}

func (vm *VM) readConstant() Value {
	return vm.chunk.constants[vm.readByte()]
}

func (vm *VM) push(value Value) {
	if vm.stackTop >= len(vm.stack) {
		panic("stack overflow")
	}

	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() Value {
	if vm.stackTop == 0 {
		panic("stack underflow")
	}

	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) binaryOp(operator string) {
	b := vm.pop()
	a := vm.pop()
	switch operator {
	case "+":
		vm.push(a + b)
	case "-":
		vm.push(a - b)
	case "*":
		vm.push(a * b)
	case "/":
		vm.push(a / b)
	default:
		panic("unknown operator")
	}
}
