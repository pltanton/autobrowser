// Package config contains parser to parse configuration rules according to following formal grammar
//
// RULE -> WORDS COLON MATCHER_PROPERTY EQ
// ASSIGNMENT -> WORD ASSIGN WORDS
//
// WORDS -> WORD [WORD]*
// BROWSER_DEF -> VALUE [VALUE]*
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

// ParseInstruction parse single configuration instruction
// it might be neither an assignment or a rule
//
// RULE -> WORDS COLON MATCHER_PROPERTY EQ
//
// ASSIGNMENT -> WORD ASSIGN WORDS
// WORDS -> WORD [WORD]*
func (p *Parser) ParseInstruction() (Instruction, error) {
	p.skipEndls()
	tok := p.scanSkipSpace()
	if tok.Type == EOF {
		return Instruction{}, fmt.Errorf("EOF token reached: %w", io.EOF)
	}
	p.unscan()

	lValue, err := p.parseWordSequence()
	if err != nil {
		return Instruction{}, err
	}

	if len(lValue) == 0 {
		return Instruction{}, fmt.Errorf("expected lValue, but got empty string")
	}

	// It might be an assignment potentially
	if len(lValue) == 1 {
		tok = p.scanSkipSpace()
		if tok.Type == ASSIGN {
			assignment, err := p.parseRestOfAssignment(lValue[0])
			if err != nil {
				return Instruction{}, err
			}
			return FromAssignment(assignment), nil
		}
		p.unscan()
	}

	if tok = p.scanSkipSpace(); tok.Type != COLON {
		return Instruction{}, fmt.Errorf("expected COLON, but got %v", tok)
	}

	rule, err := p.parseRestOfRule(lValue)
	if err != nil {
		return Instruction{}, err
	}

	return FromRule(rule), nil

}

func (p *Parser) parseRestOfAssignment(name string) (Assignment, error) {
	command, err := p.parseWordSequence()
	if err != nil {
		return Assignment{}, fmt.Errorf("expected command, but got err: %w", err)
	}

	if tok := p.scanSkipSpace(); tok.Type != ENDL && tok.Type != EOF {
		return Assignment{}, fmt.Errorf("assignment should end with ENDL or EOF, but got %v", tok)
	}

	return Assignment{
		Key:   name,
		Value: command,
	}, nil
}

func (p *Parser) parseRestOfRule(lValue []string) (Rule, error) {
	var tok Token
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
		Command:  lValue,
		Matchers: matchers,
	}, nil
}

func (p *Parser) parseWordSequence() ([]string, error) {
	result := []string{}

	for tok := p.scanSkipSpace(); tok.Type == WORD; tok = p.scanSkipSpace() {
		result = append(result, tok.Value)
	}
	p.unscan()

	return result, nil
}

// parseMatcherDef parse matcher
//
// MATCHER_DEF -> MATCHER_PROPERTY EQ VALUE
// MATCHER_PROPERTY -> VALUE DOT VALUE
//
// Returns matcher type, property name, property value
func (p *Parser) parseMatcherDef() (string, string, string, error) {
	tok := p.scanSkipSpace()

	if tok.Type != WORD {
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
	if tok.Type != WORD {
		return "", "", "", fmt.Errorf("unexpected token for matcher property name, expected VALUE, but got: %v", tok)
	}
	matcherProp := tok.Value

	tok = p.scan()
	if tok.Type != EQ {
		return "", "", "", fmt.Errorf("unexpected token, expected EQ, but got: %v", tok)
	}

	tok = p.scan()
	if tok.Type != WORD {
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
