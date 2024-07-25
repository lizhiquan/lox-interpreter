package lox

import "fmt"

//go:generate stringer -type=TokenType
type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE

	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR

	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	if t.Literal == nil {
		t.Literal = "null"
	}
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, t.Literal)
}
