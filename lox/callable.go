package lox

import (
	"errors"
	"fmt"
	"time"
)

type Callable interface {
	fmt.Stringer
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

type Clock struct{}

var _ Callable = (*Clock)(nil)

func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return time.Now().Unix(), nil
}

func (c Clock) String() string {
	return "<native fn>"
}

type Function struct {
	declaration *FunctionDeclStmt
	closure     *Environment
}

var _ Callable = (*Function)(nil)

func NewFunction(declaration *FunctionDeclStmt, closure *Environment) *Function {
	return &Function{declaration: declaration, closure: closure}
}

func (f *Function) Arity() int {
	return len(f.declaration.Parameters)
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironmentWithEnclosing(f.closure)
	for i, param := range f.declaration.Parameters {
		environment.Define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.Body, environment)
	var returnErr *ReturnError
	if errors.As(err, &returnErr) {
		return returnErr.Value, nil
	}

	return nil, err
}

func (f *Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
