package network

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		// IPv4 private
		{"loopback 127.0.0.1", "127.0.0.1", true},
		{"loopback 127.255.255.255", "127.255.255.255", true},
		{"10.0.0.0", "10.0.0.0", true},
		{"10.255.255.255", "10.255.255.255", true},
		{"172.16.0.0", "172.16.0.0", true},
		{"172.31.255.255", "172.31.255.255", true},
		{"192.168.0.0", "192.168.0.0", true},
		{"192.168.255.255", "192.168.255.255", true},
		{"link-local 169.254.1.1", "169.254.1.1", true},

		// IPv6 private
		{"IPv6 loopback", "::1", true},
		{"IPv6 unique local fc00::", "fc00::1", true},
		{"IPv6 unique local fd00::", "fd00::1", true},
		{"IPv6 link-local fe80::", "fe80::1", true},

		// Public IPs
		{"Google DNS 8.8.8.8", "8.8.8.8", false},
		{"Cloudflare DNS 1.1.1.1", "1.1.1.1", false},
		{"public 203.0.113.1", "203.0.113.1", false},

		// Boundary values
		{"172.15.255.255 is public", "172.15.255.255", false},
		{"172.32.0.0 is public", "172.32.0.0", false},
		{"9.255.255.255 is public", "9.255.255.255", false},
		{"11.0.0.0 is public", "11.0.0.0", false},
		{"192.167.255.255 is public", "192.167.255.255", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP: %s", tt.ip)
			}
			if got := isPrivateIP(ip); got != tt.want {
				t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestNewSafeClient_AllowPrivateTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewSafeClient(true)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("expected successful connection to localhost, got error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNewSafeClient_AllowPrivateFalse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewSafeClient(false)
	_, err := client.Get(ts.URL)
	if err == nil {
		t.Fatal("expected error when connecting to localhost with allowPrivate=false, got nil")
	}
}

func TestNewSafeClient_Timeout(t *testing.T) {
	client := NewSafeClient(true)
	if client.Timeout != 30*time.Second {
		t.Errorf("expected client timeout 30s, got %v", client.Timeout)
	}
}
