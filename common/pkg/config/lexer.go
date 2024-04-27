package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"unicode"
)

type CharacterClass func(rune) bool

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	EQ
	ASSIGN
	DOT
	COMMA
	COLON
	SEMICOLON
	WORD
	SPACE
	COMMENT
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
	case ASSIGN:
		return "ASSIGN"
	case DOT:
		return "DOT"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case WORD:
		return "WORD"
	case SPACE:
		return "SPACE"
	case ENDL:
		return "ENDL"
	case COMMENT:
		return "COMMENT"
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

func (l *Lexer) Next() Token {
	r := l.readRune()

	if WhitespaceClass(r) {
		l.unreadRune()
		return l.scanWhitespaces()
	} else if ValueClass(r) {
		l.unreadRune()
		return l.scanValue()
	}

	switch r {
	case rune(0):
		return Token{EOF, ""}
	case '\'':
		l.unreadRune()
		return l.scanEscapedValue()
	case '#':
		l.unreadRune()
		return l.scanComment()
	case ':':
		if l.readRune() == '=' {
			return Token{ASSIGN, ":="}
		} else {
			l.unreadRune()
			return Token{COLON, string(r)}
		}
	case '=':
		return Token{EQ, string(r)}
	case '.':
		return Token{DOT, string(r)}
	case ';':
		return Token{SEMICOLON, string(r)}
	case ',':
		return Token{COMMA, string(r)}
	case '\r':
		if l.readRune() == '\n' { // newline is \r\n on Windows
			return Token{ENDL, "\r\n"}
		} else {
			l.unreadRune()
			return Token{ILLEGAL, string(r)}
		}
	case '\n':
		return Token{ENDL, string(r)}
	}

	return Token{ILLEGAL, string(r)}
}

func (l *Lexer) readRune() rune {
	r, _, err := l.input.ReadRune()
	if errors.Is(io.EOF, err) {
		return rune(0)
	} else if err != nil {
		slog.Warn("Unexpected error occured while reading configuration", "err", err)
		return rune(0)
	}
	return r
}

func (l *Lexer) unreadRune() {
	l.input.UnreadRune()
}

func (l *Lexer) scanCharclassSequence(tokenType TokenType, class CharacterClass) Token {
	var buf bytes.Buffer

	for {
		r := l.readRune()
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
	}
}

func (l *Lexer) scanComment() Token {
	return l.scanCharclassSequence(COMMENT, func(r rune) bool {
		return r != '\n' && !isEof(r)
	})
}

func (l *Lexer) scanWhitespaces() Token {
	return l.scanCharclassSequence(SPACE, WhitespaceClass)
}

func (l *Lexer) scanValue() Token {
	return l.scanCharclassSequence(WORD, ValueClass)
}

func (l *Lexer) scanEscapedValue() Token {
	// Skip 1st quuote
	l.readRune()
	var buf bytes.Buffer
	for {
		r := l.readRune()

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
		Type:  WORD,
		Value: buf.String(),
	}
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
