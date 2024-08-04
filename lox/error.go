package lox

import "fmt"

type ParseError struct {
	token   Token
	message string
}

func NewParseError(token Token, message string) *ParseError {
	return &ParseError{token: token, message: message}
}

func (e *ParseError) Error() string {
	if e.token.Type == EOF {
		return fmt.Sprintf("[line %d] Error at end: %s", e.token.Line, e.message)
	}

	return fmt.Sprintf("[line %d] Error at '%s': %s", e.token.Line, e.token.Lexeme, e.message)
}

type RuntimeError struct {
	token   Token
	message string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", e.token.Line, e.message)
}

type ReturnError struct {
	Value any
}

func NewReturnError(value any) *ReturnError {
	return &ReturnError{Value: value}
}

func (e *ReturnError) Error() string {
	return fmt.Sprintf("return %s", e.Value)
}
