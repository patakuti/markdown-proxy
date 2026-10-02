package opener

import (
	"net"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildURL_LocalAbsolute(t *testing.T) {
	// Use a real absolute path: on Windows "/home/user" has no drive letter
	// and filepath.Abs would prepend the current drive.
	abs := filepath.Join(t.TempDir(), "doc.md")
	url, err := BuildURL(abs, 9080)
	if err != nil {
		t.Fatal(err)
	}
	want := "http://localhost:9080/local/" + strings.TrimPrefix(filepath.ToSlash(abs), "/")
	if url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestBuildURL_LocalRelative(t *testing.T) {
	url, err := BuildURL("file.md", 9080)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "http://localhost:9080/local/") {
		t.Errorf("expected /local/ prefix, got %q", url)
	}
	if !strings.HasSuffix(url, "/file.md") {
		t.Errorf("expected to end with /file.md, got %q", url)
	}
}

func TestBuildURL_HTTPS(t *testing.T) {
	url, err := BuildURL("https://github.com/user/repo", 9080)
	if err != nil {
		t.Fatal(err)
	}
	want := "http://localhost:9080/https/github.com/user/repo"
	if url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestBuildURL_HTTP(t *testing.T) {
	url, err := BuildURL("http://example.com/doc.md", 9080)
	if err != nil {
		t.Fatal(err)
	}
	want := "http://localhost:9080/http/example.com/doc.md"
	if url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestBuildURL_CustomPort(t *testing.T) {
	url, err := BuildURL("file.md", 8080)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "http://localhost:8080/local/") {
		t.Errorf("expected port 8080, got %q", url)
	}
}

func TestBuildURL_RelativeParentDir(t *testing.T) {
	url, err := BuildURL("../other/file.md", 9080)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "http://localhost:9080/local/") {
		t.Errorf("expected /local/ prefix, got %q", url)
	}
	// Should not contain ".." after resolution
	if strings.Contains(url, "..") {
		t.Errorf("expected resolved path without '..', got %q", url)
	}
	if !strings.HasSuffix(url, "/other/file.md") {
		t.Errorf("expected to end with /other/file.md, got %q", url)
	}
}

func TestBuildURL_HTTPSWithPath(t *testing.T) {
	url, err := BuildURL("https://github.com/user/repo/blob/main/README.md", 9080)
	if err != nil {
		t.Fatal(err)
	}
	want := "http://localhost:9080/https/github.com/user/repo/blob/main/README.md"
	if url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestIsServerRunning_NoServer(t *testing.T) {
	// Use a port that is very unlikely to be in use.
	if IsServerRunning(19999) {
		t.Error("expected no server running on port 19999")
	}
}

func TestIsServerRunning_WithServer(t *testing.T) {
	// Start a temporary TCP listener.
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	// Extract the port from the listener address.
	port := ln.Addr().(*net.TCPAddr).Port
	if !IsServerRunning(port) {
		t.Errorf("expected server running on port %d", port)
	}
}
