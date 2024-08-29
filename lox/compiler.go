package lox

import (
	"fmt"
	"math"
	"os"
)

type compiler struct {
	scanner        *Scanner
	compilingChunk *Chunk
	current        Token
	previous       Token
	hadError       bool
	rules          map[TokenType]parseRule
}

func compile(source string, chunk *Chunk) bool {
	c := &compiler{
		scanner:        NewScanner(source),
		compilingChunk: chunk,
	}

	c.rules = map[TokenType]parseRule{
		LEFT_PAREN:    {c.grouping, nil, PREC_NONE},
		RIGHT_PAREN:   {nil, nil, PREC_NONE},
		LEFT_BRACE:    {nil, nil, PREC_NONE},
		RIGHT_BRACE:   {nil, nil, PREC_NONE},
		COMMA:         {nil, nil, PREC_NONE},
		DOT:           {nil, nil, PREC_NONE},
		MINUS:         {c.unary, c.binary, PREC_TERM},
		PLUS:          {nil, c.binary, PREC_TERM},
		SEMICOLON:     {nil, nil, PREC_NONE},
		SLASH:         {nil, c.binary, PREC_FACTOR},
		STAR:          {nil, c.binary, PREC_FACTOR},
		BANG:          {nil, nil, PREC_NONE},
		BANG_EQUAL:    {nil, nil, PREC_NONE},
		EQUAL:         {nil, nil, PREC_NONE},
		EQUAL_EQUAL:   {nil, nil, PREC_NONE},
		GREATER:       {nil, nil, PREC_NONE},
		GREATER_EQUAL: {nil, nil, PREC_NONE},
		LESS:          {nil, nil, PREC_NONE},
		LESS_EQUAL:    {nil, nil, PREC_NONE},
		IDENTIFIER:    {nil, nil, PREC_NONE},
		STRING:        {nil, nil, PREC_NONE},
		NUMBER:        {c.number, nil, PREC_NONE},
		AND:           {nil, nil, PREC_NONE},
		CLASS:         {nil, nil, PREC_NONE},
		ELSE:          {nil, nil, PREC_NONE},
		FALSE:         {nil, nil, PREC_NONE},
		FOR:           {nil, nil, PREC_NONE},
		FUN:           {nil, nil, PREC_NONE},
		IF:            {nil, nil, PREC_NONE},
		NIL:           {nil, nil, PREC_NONE},
		OR:            {nil, nil, PREC_NONE},
		PRINT:         {nil, nil, PREC_NONE},
		RETURN:        {nil, nil, PREC_NONE},
		SUPER:         {nil, nil, PREC_NONE},
		THIS:          {nil, nil, PREC_NONE},
		TRUE:          {nil, nil, PREC_NONE},
		VAR:           {nil, nil, PREC_NONE},
		WHILE:         {nil, nil, PREC_NONE},
		EOF:           {nil, nil, PREC_NONE},
	}

	c.advance()
	c.expression()
	c.consume(EOF, "Expect end of expression.")
	c.end()
	return !c.hadError
}

func (c *compiler) advance() {
	c.previous = c.current

	for {
		var err error
		c.current, err = c.scanner.scanToken()
		if err == nil {
			break
		}

		fmt.Fprintln(os.Stderr, err)
	}
}

func (c *compiler) expression() {
	c.parsePrecedence(PREC_ASSIGNMENT)
}

type precedence int

const (
	PREC_NONE       precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type parseRule struct {
	prefix     func()
	infix      func()
	precedence precedence
}

func (c *compiler) parsePrecedence(prec precedence) {
	c.advance()
	prefixRule := c.getRule(c.previous.Type).prefix
	if prefixRule == nil {
		c.error("Expect expression.")
		return
	}

	prefixRule()

	for prec <= c.getRule(c.current.Type).precedence {
		c.advance()
		infixRule := c.getRule(c.previous.Type).infix
		infixRule()
	}
}

func (c *compiler) getRule(tokenType TokenType) parseRule {
	return c.rules[tokenType]
}

func (c *compiler) number() {
	c.emitConstant(Value(c.previous.Literal.Value.(float64)))
}

func (c *compiler) grouping() {
	c.expression()
	c.consume(RIGHT_PAREN, "Expect ')' after expression.")
}

func (c *compiler) unary() {
	operatorType := c.previous.Type

	// compile the operand
	c.parsePrecedence(PREC_UNARY)

	switch operatorType {
	case MINUS:
		c.emitByte(OP_NEGATE)
	default:
		return
	}
}

func (c *compiler) binary() {
	operatorType := c.previous.Type
	rule := c.getRule(operatorType)
	c.parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case PLUS:
		c.emitByte(OP_ADD)
	case MINUS:
		c.emitByte(OP_SUBTRACT)
	case STAR:
		c.emitByte(OP_MULTIPLY)
	case SLASH:
		c.emitByte(OP_DIVIDE)
	}
}

func (c *compiler) emitConstant(value Value) {
	c.emitBytes(OP_CONSTANT, c.makeConstant(value))
}

func (c *compiler) makeConstant(value Value) byte {
	constant := c.currentChunk().addConstant(value)
	if constant > math.MaxUint8 {
		c.error("Too many constants in one chunk.")
		return 0
	}

	return byte(constant)
}

func (c *compiler) errorAtCurrent(message string) {
	c.errorAt(c.current, message)
}

func (c *compiler) errorAt(token Token, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line)

	if token.Type == EOF {
		fmt.Fprint(os.Stderr, " at end")
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", token.Lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
	c.hadError = true
}

func (c *compiler) error(message string) {
	c.errorAt(c.previous, message)
}

func (c *compiler) consume(expectedType TokenType, message string) {
	if c.current.Type == expectedType {
		c.advance()
		return
	}

	c.errorAtCurrent(message)
}

func (c *compiler) end() {
	c.emitReturn()

	if os.Getenv("DEBUG_PRINT_CODE") == "1" && !c.hadError {
		c.currentChunk().disassemble("code")
	}
}

func (c *compiler) currentChunk() *Chunk {
	return c.compilingChunk
}

func (c *compiler) emitByte(b byte) {
	c.currentChunk().write(b, c.previous.Line)
}

func (c *compiler) emitBytes(b1 byte, b2 byte) {
	c.emitByte(b1)
	c.emitByte(b2)
}

func (c *compiler) emitReturn() {
	c.emitByte(OP_RETURN)
}
