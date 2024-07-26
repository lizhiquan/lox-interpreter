package lox

import (
	"fmt"
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
	default:
		return fmt.Errorf("[line %d] Error: Unexpected character: %s", s.line, string(r))
	}

	return nil
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(s.source[s.current:])
	return r
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
	token := Token{Type: tokenType, Lexeme: s.source[s.start:s.current], Literal: literal, Line: s.line}
	s.tokens = append(s.tokens, token)
}
