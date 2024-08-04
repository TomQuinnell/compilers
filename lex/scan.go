package lex

import (
	"errors"
	d "example/compilers/domain"
	"example/compilers/util"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type IScanner interface {
	Scan() ([]d.Token, error)
}

type ErrScan struct {
	errs []error
}

func (e ErrScan) Error() string {
	var sb strings.Builder
	sb.WriteString("scan error(s):")
	for _, e := range e.errs {
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}

type Scanner struct {
	source string
	tokens []*d.Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) Scan() ([]*d.Token, error) {
	errs := make([]error, 0)
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return nil, ErrScan{errs: errs}
	}

	s.tokens = append(s.tokens, d.NewToken(d.EOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	// log.Trace().Str("current_char", fmt.Sprintf("%c", c)).Msg("Current char")

	switch c {
	// Single characters
	case '(':
		s.addToken(d.LEFT_PAREN)
		return nil
	case ')':
		s.addToken(d.RIGHT_PAREN)
		return nil
	case '{':
		s.addToken(d.LEFT_BRACE)
		return nil
	case '}':
		s.addToken(d.RIGHT_BRACE)
		return nil
	case ',':
		s.addToken(d.COMMA)
		return nil
	case '.':
		s.addToken(d.DOT)
		return nil
	case '-':
		s.addToken(d.MINUS)
		return nil
	case '+':
		s.addToken(d.PLUS)
		return nil
	case ';':
		s.addToken(d.SEMICOLON)
		return nil
	case '*':
		s.addToken(d.STAR)
		return nil

	// Operators
	case '!':
		if s.matches('=') {
			s.addToken(d.BANG_EQUAL)
		} else {
			s.addToken(d.BANG)
		}
		return nil
	case '=':
		if s.matches('=') {
			s.addToken(d.EQUAL_EQUAL)
		} else {
			s.addToken(d.EQUAL)
		}
		return nil
	case '<':
		if s.matches('=') {
			s.addToken(d.LESS_EQUAL)
		} else {
			s.addToken(d.LESS)
		}
		return nil
	case '>':
		if s.matches('=') {
			s.addToken(d.GREATER_EQUAL)
		} else {
			s.addToken(d.GREATER)
		}
		return nil

	// Ignore whitespaces
	case ' ', '\r', '\t':
		s.start++
		return nil
	case '\n':
		s.start++
		s.line++
		return nil

	// String literals
	case '"':
		for s.peek() != '"' && !s.isAtEnd() {
			if s.peek() == '\n' {
				s.line++
			}
			s.advance()
		}

		if s.isAtEnd() {
			return errors.New("unterminated string")
		}

		s.advance()

		value, err := s.substring(s.start+1, s.current-1)
		if err != nil {
			log.Err(err).Msg("Failed to parse string literal")
			return errors.New("failed to parse string literal")
		}
		s.addTokenWithLiteral(d.STRING, value)
		return nil

	// Longer lexemes
	case '/':
		if s.matches('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(d.SLASH)
		}
		return nil
	default:
		if util.IsDigit(c) {
			for util.IsDigit(s.peek()) {
				s.advance()
			}

			// Fraction
			if s.peek() == '.' && util.IsDigit(s.peekNext()) {
				s.advance()

				for util.IsDigit(s.peek()) {
					s.advance()
				}
			}

			text, err := s.currentSlice()
			if err != nil {
				log.Err(err).Msg("Failed to parse float literal string")
				return errors.New("invalid string literal")
			}
			value, err := strconv.ParseFloat(text, 64)
			if err != nil {
				log.Err(err).Msg("Failed to parse float literal")
				return errors.New("invalid float literal")
			}

			s.addTokenWithLiteral(d.NUMBER, value)
			return nil
		} else if util.IsAlpha(c) {
			for util.IsAlphaNumeric(s.peek()) {
				s.advance()
			}

			// Check if keyword matched
			text, err := s.currentSlice()
			if err != nil {
				log.Err(err).Msg("Failed to parse keyword")
				return errors.New("invalid keyword")
			}
			kind, ok := d.Keywords[text]
			if ok {
				s.addToken(kind)
			} else {
				s.addToken(d.IDENTIFIER)
			}
			return nil
		} else {
			return fmt.Errorf("unexpected character: '%c'", c)
		}
	}
}

func (s *Scanner) advance() rune {
	c := s.currentChar()
	s.current++
	return c
}

func (s *Scanner) matches(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.currentChar() != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}

	return s.currentChar()
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return util.RuneAt(s.source, s.current+1)
}

func (s *Scanner) currentChar() rune {
	if s.current >= len(s.source) {
		return '\x00'
	}
	return util.RuneAt(s.source, s.current)
}

func (s *Scanner) addToken(kind d.TokenType) {
	s.addTokenWithLiteral(kind, nil)
}

func (s *Scanner) addTokenWithLiteral(kind d.TokenType, literal interface{}) {
	text, err := s.currentSlice()
	if err != nil {
		log.Err(err).Msg("Failed to parse token literal")
	}
	newToken := d.NewToken(kind, text, literal, s.line)
	s.tokens = append(s.tokens, newToken)
}

func (s *Scanner) currentSlice() (string, error) {
	return s.substring(s.start, s.current)
}

func (s *Scanner) substring(st int, end int) (string, error) {
	if st < 0 || end >= len(s.source) {
		return "", fmt.Errorf("out of range access: %d %d", st, end)
	}
	return s.source[st:end], nil
}
