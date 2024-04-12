package config

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		name    string
		inStr   string
		want    Token
		wantErr bool
	}{
		{
			name:    "Lex eq",
			inStr:   "=",
			want:    Token{EQ, "="},
			wantErr: false,
		},
		{
			name:    "Lex dot",
			inStr:   ".",
			want:    Token{DOT, "."},
			wantErr: false,
		},
		{
			name:    "Lex value",
			inStr:   "here_is-va1ue",
			want:    Token{VALUE, "here_is-va1ue"},
			wantErr: false,
		},
		{
			name:    "Lex escaped value",
			inStr:   "'here.is,escaped=value\t*** () ?? {} <> 💀'",
			want:    Token{VALUE, "here.is,escaped=value\t*** () ?? {} <> 💀"},
			wantErr: false,
		},
		{
			name:    "Lex comma",
			inStr:   ",",
			want:    Token{COMMA, ","},
			wantErr: false,
		},
		{
			name:    "Lex single space",
			inStr:   " ",
			want:    Token{SPACE, " "},
			wantErr: false,
		},
		{
			name:    "Lex subsequent spaces",
			inStr:   "  \t\t  ",
			want:    Token{SPACE, "  \t\t  "},
			wantErr: false,
		},
		{
			name:    "Lex colon",
			inStr:   ":",
			want:    Token{COLON, ":"},
			wantErr: false,
		},
		{
			name:    "Lex semicolon",
			inStr:   ";",
			want:    Token{SEMICOLON, ";"},
			wantErr: false,
		},
		{
			name:    "Lex endline",
			inStr:   "\n",
			want:    Token{ENDL, "\n"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input: bufio.NewReader(strings.NewReader(tt.inStr)),
			}
			got, err := l.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_FullSequence(t *testing.T) {
	input := `
firefox:url.regex='.*foo.*';app.class=telegram
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

	for tok, err := l.Next(); tok.Type != EOF; tok, err = l.Next() {
		if err != nil {
			t.Errorf("Unexpected error while tokenizing: %v", err)
			return
		}
		tokens = append(tokens, tok)
	}

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Tokens differs:\n%v\n%v", tokens, expected)
	}
}
