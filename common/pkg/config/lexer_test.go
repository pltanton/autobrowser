package config

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		name  string
		inStr string
		want  Token
	}{
		{
			name:  "Lex eq",
			inStr: "=",
			want:  Token{EQ, "="},
		},
		{
			name:  "Lex dot",
			inStr: ".",
			want:  Token{DOT, "."},
		},
		{
			name:  "Lex value",
			inStr: "here_is-va1ue",
			want:  Token{VALUE, "here_is-va1ue"},
		},
		{
			name:  "Lex escaped value",
			inStr: "'here.is,escaped=value\t*** () ?? {} <> 💀'",
			want:  Token{VALUE, "here.is,escaped=value\t*** () ?? {} <> 💀"},
		},
		{
			name:  "Lex comma",
			inStr: ",",
			want:  Token{COMMA, ","},
		},
		{
			name:  "Lex single space",
			inStr: " ",
			want:  Token{SPACE, " "},
		},
		{
			name:  "Lex subsequent spaces",
			inStr: "  \t\t  ",
			want:  Token{SPACE, "  \t\t  "},
		},
		{
			name:  "Lex colon",
			inStr: ":",
			want:  Token{COLON, ":"},
		},
		{
			name:  "Lex semicolon",
			inStr: ";",
			want:  Token{SEMICOLON, ";"},
		},
		{
			name:  "Lex endline",
			inStr: "\n",
			want:  Token{ENDL, "\n"},
		},

		{
			name:  "Lex comment",
			inStr: "# Hello.={} ;:",
			want:  Token{COMMENT, "# Hello.={} ;:"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input: bufio.NewReader(strings.NewReader(tt.inStr)),
			}
			got := l.Next()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_FullSequence(t *testing.T) {
	input := `
firefox:url.regex='.*foo.*';app.class=telegram # Commentary with row description
'firefox -p work':url.host='github.com'`

	expected := []Token{
		{ENDL, "\n"},

		{VALUE, "firefox"},
		{COLON, ":"},
		{VALUE, "url"},
		{DOT, "."},
		{VALUE, "regex"},
		{EQ, "="},
		{VALUE, ".*foo.*"},
		{SEMICOLON, ";"},
		{VALUE, "app"},
		{DOT, "."},
		{VALUE, "class"},
		{EQ, "="},
		{VALUE, "telegram"},

		{SPACE, " "},
		{COMMENT, "# Commentary with row description"},
		{ENDL, "\n"},

		{VALUE, "firefox -p work"},
		{COLON, ":"},
		{VALUE, "url"},
		{DOT, "."},
		{VALUE, "host"},
		{EQ, "="},
		{VALUE, "github.com"},
	}

	l := &Lexer{
		input: bufio.NewReader(strings.NewReader(input)),
	}

	tokens := make([]Token, 0)

	for tok := l.Next(); tok.Type != EOF; tok = l.Next() {
		tokens = append(tokens, tok)
	}

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Tokens differs:\n%v\n%v", tokens, expected)
	}
}
