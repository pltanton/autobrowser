package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseRule(t *testing.T) {
	tests := []struct {
		name     string
		inStr    string
		want     Rule
		wantOver bool
		wantErr  bool
	}{
		{
			name:  "Successfully parse",
			inStr: "firefox:url.regex='.*example.*';app.name=slack;app.foo=bar",
			want: Rule{
				Target: "firefox",
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
			name:     "Should be over",
			inStr:    "   \n\n\n   \n",
			wantOver: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.inStr))
			got, over, err := p.ParseRule()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (over != true) == tt.wantOver {
				t.Errorf("Parser.ParseRule() over = %v, wantOver %v", over, tt.wantOver)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.ParseRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
