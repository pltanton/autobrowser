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
		name           string
		inStr          string
		wantRule       Rule
		wantAssignment Assignment
		wantErr        bool
		err            error
	}{
		{
			name:  "Successfully parse rule",
			inStr: "firefox -command    {}:url.regex='.*example.*'; app.name=slack;app.foo=bar # this is test rule",
			wantRule: Rule{
				Command: []string{"firefox", "-command", "{}"},
				Matchers: map[string]MatcherProps{
					"url": {"regex": ".*example.*"},
					"app": {"name": "slack", "foo": "bar"},
				},
			},
		},
		{
			name:  "Successfully parse assignment",
			inStr: "foo   := biz 'buz -bin {}'",
			wantAssignment: Assignment{
				Key:   "foo",
				Value: []string{"biz", "buz -bin {}"},
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
			got, err := p.ParseInstruction()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseInstruction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (tt.err != nil) && errors.Is(tt.err, err) {
				t.Errorf("Parser.ParseInstruction() error = %v, should be err %v", err, tt.err)
				return
			}
			if !reflect.ValueOf(tt.wantRule).Field(0).IsZero() {
				rule, ok := got.Rule()
				if !ok {
					t.Errorf("Expected to get rule instruction but got %+v", got)
					return
				}
				if !reflect.DeepEqual(rule, tt.wantRule) {
					t.Errorf("Parser.ParseInstruction() = %v, want %v", rule, tt.wantRule)
					return
				}
			}
			if !reflect.ValueOf(tt.wantAssignment).Field(0).IsZero() {
				assignment, ok := got.Assignment()
				if !ok {
					t.Errorf("Expected to get assignment instruction but got %+v", got)
					return
				}
				if !reflect.DeepEqual(assignment, tt.wantAssignment) {
					t.Errorf("Parser.ParseInstruction() = %v, want %v", assignment, tt.wantAssignment)
					return
				}
			}
		})
	}
}
