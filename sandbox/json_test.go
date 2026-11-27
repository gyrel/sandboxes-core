package sandbox

import (
	"bytes"
	"encoding/json"
	"maps"
	"testing"
)

func TestSandboxLabelsJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
	}{
		{name: "nil"},
		{name: "empty", labels: map[string]string{}},
		{name: "set", labels: map[string]string{"team": "platform", "archived": ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&Sandbox{ID: "sb1", Labels: tt.labels})
			if err != nil {
				t.Fatal(err)
			}
			if got := bytes.Contains(data, []byte(`"labels":`)); got != (len(tt.labels) > 0) {
				t.Fatalf("unexpected labels presence in %s", data)
			}
			var got Sandbox
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != "sb1" || !maps.Equal(got.Labels, tt.labels) || (len(tt.labels) == 0 && got.Labels != nil) {
				t.Fatalf("unexpected decoded sandbox ID %q and labels %v", got.ID, got.Labels)
			}
		})
	}
}

func TestSandboxLabelsJSONDecode(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "legacy", data: `{"ID":"sb1","State":"stopped"}`},
		{name: "empty", data: `{"ID":"sb1","State":"stopped","labels":{}}`},
		{name: "null", data: `{"ID":"sb1","State":"stopped","labels":null}`},
		{name: "unknown field", data: `{"labels":{},"extra":1}`, wantErr: true},
		{name: "unknown nested field", data: `{"Meta":{"extra":1}}`, wantErr: true},
		{name: "invalid labels", data: `{"labels":{"authToken":"x"}}`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sandbox{Labels: map[string]string{"old": "value"}}
			err := json.Unmarshal([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, want error %t", err, tt.wantErr)
			}
			if err == nil && (got.ID != "sb1" || got.State != Stopped || got.Labels != nil) {
				t.Fatalf("unexpected decoded sandbox ID %q, state %q and labels %v", got.ID, got.State, got.Labels)
			}
		})
	}
}
