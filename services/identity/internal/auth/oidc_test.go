package auth

import (
	"testing"
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{"valid", "Bearer mytoken", "mytoken", false},
		{"missing", "", "", true},
		{"invalid format", "Basic something", "", true},
		{"no token", "Bearer", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
