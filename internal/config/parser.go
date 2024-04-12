// Package config contains parser to parse configuration rules according to following formal grammar
//
// RULE -> BROWSER_DEF COLON MATCHER_DEF [SEMICOLON RULE_DEF]*
// BROWSER_DEF -> VALUE
// MATCHER_DEF -> MATCHER_PROPERTY EQ VALUE
// MATCHER_PROPERTY -> VALUE DOT VALUE
//
// Example of single rule:
//
//	`firefox:url.regex='.*exapmle\.com.*';app.name=telegram`
package config

import (
	"fmt"
	"io"
)

// Parses configuration
type Parser struct {
	l *Lexer

	// Current window size is naive with 1 element wide, could be refactod to use wider window
	buf struct {
		token Token
		n     int
	}
}

func NewParser(in io.Reader) *Parser {
	return &Parser{
		l: NewLexer(in),
	}
}

// Parse rule parses single rule
//
// RULE -> RULE_NAME COLON MATCHER_PROPERTY EQ
func (p *Parser) ParseRule() (Rule, bool, error) {
	p.skipEndls()

	tok, err := p.scanSkipSpace()
	if err != nil {
		return Rule{}, false, err
	}

	if tok.Type == EOF {
		return Rule{}, true, nil
	}

	if tok.Type != VALUE {
		return Rule{}, false, fmt.Errorf("unexpected token, expected VALUE, but got: %v", tok)
	}
	browserDef := tok.Value

	if tok, err = p.scan(); err != nil {
		return Rule{}, false, err
	}
	if tok.Type != COLON {
		return Rule{}, false, fmt.Errorf("unexpected token, expected COLON, but got: %v", tok)
	}

	matchers := make(map[string]MatcherProps)

	for {
		matcherName, propName, propValue, err := p.parseMatcherDef()
		if err != nil {
			return Rule{}, false, fmt.Errorf("failed to parse matcher definition: %w", err)
		}

		matcher, ok := matchers[matcherName]
		if !ok {
			matcher = make(MatcherProps)
			matchers[matcherName] = matcher
		}

		matcher[propName] = propValue

		if tok, err = p.scanSkipSpace(); err != nil {
			return Rule{}, false, err
		}
		if tok.Type == ENDL || tok.Type == EOF {
			break
		} else if tok.Type != SEMICOLON {
			return Rule{}, false, fmt.Errorf("failed to parse patchers definitions, expected SEMICOLON or end of rule, but got %v", tok)
		}
	}

	return Rule{
		Target:   browserDef,
		Matchers: matchers,
	}, false, nil
}

// parseMatcherDef parse matcher
//
// MATCHER_DEF -> MATCHER_PROPERTY EQ VALUE
// MATCHER_PROPERTY -> VALUE DOT VALUE
func (p *Parser) parseMatcherDef() (string, string, string, error) {
	tok, err := p.scan()
	if err != nil {
		return "", "", "", err
	}

	if tok.Type != VALUE {
		return "", "", "", fmt.Errorf("unexpected token for matcher type, expected VALUE, but got: %v", tok)
	}
	matcherType := tok.Value

	if tok, err = p.scan(); err != nil {
		return "", "", "", err
	}
	if tok.Type != DOT {
		return "", "", "", fmt.Errorf("unexpected token, expected DOT, but got: %v", tok)
	}

	if tok, err = p.scan(); err != nil {
		return "", "", "", err
	}
	if tok.Type != VALUE {
		return "", "", "", fmt.Errorf("unexpected token for matcher property name, expected VALUE, but got: %v", tok)
	}
	matcherProp := tok.Value

	if tok, err = p.scan(); err != nil {
		return "", "", "", err
	}
	if tok.Type != EQ {
		return "", "", "", fmt.Errorf("unexpected token, expected EQ, but got: %v", tok)
	}

	if tok, err = p.scan(); err != nil {
		return "", "", "", err
	}
	if tok.Type != VALUE {
		return "", "", "", fmt.Errorf("unexpected token for matcher property value, expected VALUE, but got: %v", tok)
	}
	matcherPropValue := tok.Value

	return matcherType, matcherProp, matcherPropValue, err
}

func (p *Parser) scan() (Token, error) {
	if p.buf.n != 0 {
		p.buf.n--
		return p.buf.token, nil
	}

	var err error
	p.buf.token, err = p.l.Next()
	if err != nil {
		return p.buf.token, err
	}

	return p.buf.token, nil
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) skipEndls() error {
	var tok Token
	var err error
	for tok, err = p.scan(); (tok.Type == SPACE || tok.Type == ENDL) && err == nil; tok, err = p.scan() {
	}

	if err != nil {
		return err
	}

	p.unscan()
	return nil
}

func (p *Parser) scanSkipSpace() (Token, error) {
	t, err := p.scan()
	if t.Type == SPACE {
		return p.scan()
	}

	return t, err
}
