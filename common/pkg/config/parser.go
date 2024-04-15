// Package config contains parser to parse configuration rules according to following formal grammar
//
// RULE -> BROWSER_DEF COLON MATCHER_DEF [SEMICOLON RULE_DEF]*
// BROWSER_DEF -> VALUE [SPACE VALUE]*
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
//
// Returns EOF error if sequence is over
func (p *Parser) ParseRule() (Rule, error) {
	p.skipEndls()

	tok := p.scanSkipSpace()
	if tok.Type == EOF {
		return Rule{}, fmt.Errorf("EOF token reached: %w", io.EOF)
	}

	p.unscan()
	browserDef, err := p.parseBrowserCommand()
	if err != nil {
		return Rule{}, err
	}

	matchers := make(map[string]MatcherProps)

	for {
		matcherName, propName, propValue, err := p.parseMatcherDef()
		if err != nil {
			return Rule{}, fmt.Errorf("failed to parse matcher definition: %w", err)
		}

		matcher, ok := matchers[matcherName]
		if !ok {
			matcher = make(MatcherProps)
			matchers[matcherName] = matcher
		}

		matcher[propName] = propValue

		tok = p.scanSkipSpace()
		if tok.Type == ENDL || tok.Type == EOF {
			break
		} else if tok.Type != SEMICOLON {
			return Rule{}, fmt.Errorf("failed to parse patchers definitions, expected SEMICOLON or end of rule, but got %v", tok)
		}
	}

	return Rule{
		Command:  browserDef,
		Matchers: matchers,
	}, nil
}

func (p *Parser) parseBrowserCommand() ([]string, error) {
	result := []string{}

	for tok := p.scanSkipSpace(); tok.Type != COLON; tok = p.scanSkipSpace() {
		if tok.Type != VALUE {
			return nil, fmt.Errorf("expected VALUE or COLON, but got %v", tok)
		}

		result = append(result, tok.Value)
	}

	return result, nil
}

// parseMatcherDef parse matcher
//
// MATCHER_DEF -> MATCHER_PROPERTY EQ VALUE
// MATCHER_PROPERTY -> VALUE DOT VALUE
//
// Returns matcher type, property name, property value
func (p *Parser) parseMatcherDef() (string, string, string, error) {
	tok := p.scan()

	if tok.Type != VALUE {
		return "", "", "", fmt.Errorf("unexpected token for matcher type, expected VALUE, but got: %v", tok)
	}
	matcherType := tok.Value

	tok = p.scan()
	if tok.Type != DOT {
		tok = p.scanSkipSpace()

		// If after dot there is end of rule, just return it as is
		if tok.Type == ENDL || tok.Type == EOF || tok.Type == SEMICOLON {
			return matcherType, "", "", nil
		}

		return "", "", "", fmt.Errorf("unexpected token, expected DOT, ENDL, SEMICOLON or EOF, but got: %v", tok)
	}

	tok = p.scan()
	if tok.Type != VALUE {
		return "", "", "", fmt.Errorf("unexpected token for matcher property name, expected VALUE, but got: %v", tok)
	}
	matcherProp := tok.Value

	tok = p.scan()
	if tok.Type != EQ {
		return "", "", "", fmt.Errorf("unexpected token, expected EQ, but got: %v", tok)
	}

	tok = p.scan()
	if tok.Type != VALUE {
		return "", "", "", fmt.Errorf("unexpected token for matcher property value, expected VALUE, but got: %v", tok)
	}
	matcherPropValue := tok.Value

	return matcherType, matcherProp, matcherPropValue, nil
}

// scan next token by lexer, using 1 wide window buffer
func (p *Parser) scanRaw() Token {
	if p.buf.n != 0 {
		p.buf.n--
		return p.buf.token
	}

	p.buf.token = p.l.Next()

	return p.buf.token
}

// scan skans with skip of comment
func (p *Parser) scan() Token {
	tok := p.scanRaw()
	if tok.Type == COMMENT {
		return p.scanRaw()
	}
	return tok
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) skipEndls() error {
	var tok Token
	for tok = p.scan(); tok.Type == SPACE || tok.Type == ENDL; tok = p.scan() {
		// Just skip it
	}

	// We did at least one scan, so unscan is necessary
	p.unscan()
	return nil
}

// scanSkipSpace scans skipping spaces
func (p *Parser) scanSkipSpace() Token {
	t := p.scan()
	if t.Type == SPACE {
		return p.scan()
	}

	return t
}
