package lox

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int // the first character in the lexeme being scanned
	current int // the character currently being considered
	line    int // tracks what source line `current` is on
	errs    []error
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, []error) {
	for !s.isAtEnd() {
		token, err := s.scanToken()
		if err != nil {
			s.errs = append(s.errs, err)
		}

		s.tokens = append(s.tokens, token)
	}

	return s.tokens, s.errs
}

func (s *Scanner) scanToken() (Token, error) {
	s.skipWhitespace()
	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(EOF), nil
	}

	r := s.advance()
	switch r {
	case '(':
		return s.makeToken(LEFT_PAREN), nil
	case ')':
		return s.makeToken(RIGHT_PAREN), nil
	case '{':
		return s.makeToken(LEFT_BRACE), nil
	case '}':
		return s.makeToken(RIGHT_BRACE), nil
	case ',':
		return s.makeToken(COMMA), nil
	case '.':
		return s.makeToken(DOT), nil
	case '-':
		return s.makeToken(MINUS), nil
	case '+':
		return s.makeToken(PLUS), nil
	case ';':
		return s.makeToken(SEMICOLON), nil
	case '*':
		return s.makeToken(STAR), nil
	case '!':
		if s.match('=') {
			return s.makeToken(BANG_EQUAL), nil
		} else {
			return s.makeToken(BANG), nil
		}
	case '=':
		if s.match('=') {
			return s.makeToken(EQUAL_EQUAL), nil
		} else {
			return s.makeToken(EQUAL), nil
		}
	case '<':
		if s.match('=') {
			return s.makeToken(LESS_EQUAL), nil
		} else {
			return s.makeToken(LESS), nil
		}
	case '>':
		if s.match('=') {
			return s.makeToken(GREATER_EQUAL), nil
		} else {
			return s.makeToken(GREATER), nil
		}
	case '/':
		return s.makeToken(SLASH), nil
	case '"':
		return s.string()
	default:
		if unicode.IsDigit(r) {
			return s.number()
		}
		if s.isAlpha(r) {
			return s.identifier(), nil
		}

		return Token{}, fmt.Errorf("[line %d] Error: Unexpected character: %s", s.line, string(r))
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) skipWhitespace() {
	for {
		r := s.peek()
		switch r {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.match('/') {
				// A comment goes until the end of the line.
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (s *Scanner) identifier() Token {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	return s.makeToken(tokenType)
}

func (s *Scanner) number() (Token, error) {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// fractional part
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// consume "."
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return Token{}, fmt.Errorf("[line %d] Error: Failed to parse number: %s", s.line, s.source[s.start:s.current])
	}

	token := s.makeTokenLiteral(NUMBER, value)
	return token, nil
}

func (s *Scanner) string() (Token, error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return Token{}, fmt.Errorf("[line %d] Error: Unterminated string.", s.line)
	}

	// closing "
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	token := s.makeTokenLiteral(STRING, value)
	return token, nil
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(s.source[s.current:])
	return r
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(s.source[s.current+1:])
	return r
}

func (s *Scanner) isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return s.isAlpha(r) || unicode.IsDigit(r)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	r, width := utf8.DecodeRuneInString(s.source[s.current:])
	if r != expected {
		return false
	}

	s.current += width
	return true
}

func (s *Scanner) advance() rune {
	r, width := utf8.DecodeRuneInString(s.source[s.current:])
	s.current += width
	return r
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	return s.makeTokenLiteral(tokenType, nil)
}

func (s *Scanner) makeTokenLiteral(tokenType TokenType, literal any) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  s.source[s.start:s.current],
		Literal: NewLiteral(literal),
		Line:    s.line,
	}
}
