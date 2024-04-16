package config

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseRule(t *testing.T) {
	tests := []struct {
		name    string
		inStr   string
		want    Rule
		wantErr bool
		err     error
	}{
		{
			name:  "Successfully parse",
			inStr: "firefox -command {}:url.regex='.*example.*';app.name=slack;app.foo=bar",
			want: Rule{
				Command: []string{"firefox", "-command", "{}"},
				Matchers: map[string]MatcherProps{
					"url": {"regex": ".*example.*"},
					"app": {"name": "slack", "foo": "bar"},
				},
			},
		},
		{
			name:    "Bad start of rule",
			inStr:   ";firefox:url.regex='.*example.*';app.name=slack;app.foo=bar",
			wantErr: true,
		},
		{
			name:    "Bad matcherdef",
			inStr:   "firefox:url.regex=.*example.*;app.name=slack;app.foo=bar",
			wantErr: true,
		},
		{
			name:    "Should be over",
			inStr:   "   \n\n\n   \n",
			wantErr: true,
			err:     io.EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.inStr))
			got, err := p.ParseRule()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (tt.err != nil) && errors.Is(tt.err, err) {
				t.Errorf("Parser.ParseRule() error = %v, should be err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.ParseRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
