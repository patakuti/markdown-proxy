package config

import (
	"strings"
	"testing"
)

func TestIsRemoteMode(t *testing.T) {
	tests := []struct {
		name   string
		listen string
		want   bool
	}{
		{"127.0.0.1 is local", "127.0.0.1", false},
		{"localhost is local", "localhost", false},
		{"0.0.0.0 is remote", "0.0.0.0", true},
		{"192.168.1.10 is remote", "192.168.1.10", true},
		{":: is remote", "::", true},
		{"empty string is remote", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{Listen: tt.listen}
			if got := c.IsRemoteMode(); got != tt.want {
				t.Errorf("Config{Listen: %q}.IsRemoteMode() = %v, want %v", tt.listen, got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		listen    string
		authToken string
		wantErr   bool
	}{
		{"local mode without token is OK", "127.0.0.1", "", false},
		{"local mode with token is OK", "127.0.0.1", "secret", false},
		{"remote mode with token is OK", "0.0.0.0", "secret", false},
		{"remote mode without token is error", "0.0.0.0", "", true},
		{"localhost without token is OK", "localhost", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{Listen: tt.listen, AuthToken: tt.authToken}
			err := c.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_ErrorMessage(t *testing.T) {
	c := &Config{Listen: "0.0.0.0", AuthToken: ""}
	err := c.Validate()
	if err == nil {
		t.Fatal("Validate() should return error for remote mode without token")
	}
	if !strings.Contains(err.Error(), "--auth-token") {
		t.Errorf("error message should contain '--auth-token', got: %s", err.Error())
	}
}
