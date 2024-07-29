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
		s.start = s.current
		if err := s.scanToken(); err != nil {
			s.errs = append(s.errs, err)
		}
	}
	s.tokens = append(s.tokens, Token{Type: EOF, Line: s.line})
	return s.tokens, s.errs
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// skip over comments
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// ignore whitespace
		break
	case '\n':
		s.line++
	case '"':
		return s.string()
	default:
		if unicode.IsDigit(r) {
			return s.number()
		}
		if s.isAlpha(r) {
			s.identifier()
			break
		}
		return fmt.Errorf("[line %d] Error: Unexpected character: %s", s.line, string(r))
	}

	return nil
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

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scanner) number() error {
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
		return fmt.Errorf("[line %d] Error: Failed to parse number: %s", s.line, s.source[s.start:s.current])
	}

	s.addTokenLiteral(NUMBER, value)
	return nil
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return fmt.Errorf("[line %d] Error: Unterminated string.", s.line)
	}

	// closing "
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, value)
	return nil
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

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal any) {
	token := Token{
		Type:    tokenType,
		Lexeme:  s.source[s.start:s.current],
		Literal: NewLiteral(literal),
		Line:    s.line,
	}
	s.tokens = append(s.tokens, token)
}
