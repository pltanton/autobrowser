package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type CharacterClass func(rune) bool

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	EQ
	DOT
	COMMA
	COLON
	SEMICOLON
	VALUE
	SPACE
	ENDL
)

func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case EQ:
		return "EQ"
	case DOT:
		return "DOT"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case VALUE:
		return "VALUE"
	case SPACE:
		return "SPACE"
	case ENDL:
		return "ENDL"
	}
	return "UNKNOWN"
}

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%v:%q", t.Type, t.Value)
}

type Lexer struct {
	input *bufio.Reader
}

func NewLexer(in io.Reader) *Lexer {
	return &Lexer{
		input: bufio.NewReader(in),
	}
}

func (l *Lexer) Next() (Token, error) {
	r, err := l.readRune()
	if err != nil {
		return Token{}, err
	}

	if WhitespaceClass(r) {
		l.unreadRune()
		return l.scanWhitespaces()
	} else if ValueClass(r) {
		l.unreadRune()
		return l.scanValue()
	}

	switch r {
	case rune(0):
		return Token{EOF, ""}, nil
	case '\'':
		l.unreadRune()
		return l.scanEscapedValue()
	case '=':
		return Token{EQ, string(r)}, nil
	case '.':
		return Token{DOT, string(r)}, nil
	case ':':
		return Token{COLON, string(r)}, nil
	case ';':
		return Token{SEMICOLON, string(r)}, nil
	case ',':
		return Token{COMMA, string(r)}, nil
	case '\n':
		return Token{ENDL, string(r)}, nil
	}

	return Token{}, fmt.Errorf("unexpected token %d", r)
}

func (l *Lexer) readRune() (rune, error) {
	r, _, err := l.input.ReadRune()
	if errors.Is(err, io.EOF) {
		return rune(0), nil
	}
	return r, err
}

func (l *Lexer) unreadRune() {
	l.input.UnreadRune()
}

func (l *Lexer) scanCharclassSequence(tokenType TokenType, class CharacterClass) (Token, error) {
	var buf bytes.Buffer

	for {
		r, err := l.readRune()
		if err != nil {
			return Token{}, err
		}

		if isEof(r) {
			break
		} else if !class(r) {
			l.unreadRune()
			break
		} else {
			// We don't check err here, because it is impossible
			buf.WriteRune(r)
		}
	}

	return Token{
		Type:  tokenType,
		Value: buf.String(),
	}, nil
}

func (l *Lexer) scanWhitespaces() (Token, error) {
	return l.scanCharclassSequence(SPACE, WhitespaceClass)
}

func (l *Lexer) scanValue() (Token, error) {
	return l.scanCharclassSequence(VALUE, ValueClass)
}

func (l *Lexer) scanEscapedValue() (Token, error) {
	r, err := l.readRune()
	if err != nil {
		return Token{}, err
	}
	if r != '\'' {
		return Token{}, fmt.Errorf("expected escapesequence started from ' but got %c", r)
	}

	var buf bytes.Buffer
	for {
		r, err := l.readRune()
		if err != nil {
			return Token{}, err
		}

		if isEof(r) {
			break
		} else if r == '\'' {
			break
		} else {
			// We don't check err here, because it is impossible
			buf.WriteRune(r)
		}
	}

	return Token{
		Type:  VALUE,
		Value: buf.String(),
	}, nil
}

var WhitespaceClass CharacterClass = func(r rune) bool {
	return r == ' ' || r == '\t'
}

var ValueClass CharacterClass = func(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '{' || r == '}'
}

func isEof(r rune) bool {
	return r == rune(0)
}
