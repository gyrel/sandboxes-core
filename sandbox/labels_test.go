package sandbox

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateLabels(t *testing.T) {
	full := make(map[string]string, 16)
	for i := range 16 {
		full[fmt.Sprintf("label%d", i)] = ""
	}
	overflow := make(map[string]string, 17)
	for key, value := range full {
		overflow[key] = value
	}
	overflow["extra"] = ""
	tests := []struct {
		name    string
		labels  map[string]string
		wantErr bool
	}{
		{name: "nil"},
		{name: "empty", labels: map[string]string{}},
		{name: "valid grammar", labels: map[string]string{"A": "", "9.team_A-1": " ~!"}},
		{name: "max lengths", labels: map[string]string{strings.Repeat("x", 63): strings.Repeat("~", 255)}},
		{name: "16 pairs", labels: full},
		{name: "17 pairs", labels: overflow, wantErr: true},
		{name: "empty key", labels: map[string]string{"": "x"}, wantErr: true},
		{name: "leading hyphen", labels: map[string]string{"-name": "x"}, wantErr: true},
		{name: "long key", labels: map[string]string{strings.Repeat("x", 64): "x"}, wantErr: true},
		{name: "Unicode key", labels: map[string]string{"tëam": "x"}, wantErr: true},
		{name: "space in key", labels: map[string]string{"team name": "x"}, wantErr: true},
		{name: "long value", labels: map[string]string{"name": strings.Repeat("x", 256)}, wantErr: true},
		{name: "newline", labels: map[string]string{"name": "a\nb"}, wantErr: true},
		{name: "carriage return", labels: map[string]string{"name": "a\rb"}, wantErr: true},
		{name: "NUL", labels: map[string]string{"name": "a\x00b"}, wantErr: true},
		{name: "DEL", labels: map[string]string{"name": "\x7f"}, wantErr: true},
		{name: "Unicode value", labels: map[string]string{"name": "é"}, wantErr: true},
		{name: "token", labels: map[string]string{"myToKeNs": "x"}, wantErr: true},
		{name: "key", labels: map[string]string{"myKEYs": "x"}, wantErr: true},
		{name: "secret", labels: map[string]string{"nonSECRETflag": "x"}, wantErr: true},
		{name: "password", labels: map[string]string{"myPasswordHint": "x"}, wantErr: true},
		{name: "api", labels: map[string]string{"internalAPIv2": "x"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateLabels(tt.labels); (err != nil) != tt.wantErr {
				t.Fatalf("ValidateLabels() error = %v, want error %t", err, tt.wantErr)
			}
		})
	}
}
