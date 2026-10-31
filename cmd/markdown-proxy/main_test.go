package main

import (
	"flag"
	"net"
	"reflect"
	"testing"

	"github.com/patakuti/markdown-proxy/internal/config"
)

func TestCollectServerArgs_Basic(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("port", "9080", "")
	fs.String("theme", "github", "")
	fs.Parse([]string{"-port=8080", "-theme=dark"})

	args := collectServerArgs(fs.Visit)
	want := []string{"-port=8080", "-theme=dark"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("got %v, want %v", args, want)
	}
}

func TestCollectServerArgs_SkipsVersion(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("version", false, "")
	fs.String("port", "9080", "")
	fs.Parse([]string{"-version", "-port=8080"})

	args := collectServerArgs(fs.Visit)
	want := []string{"-port=8080"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("got %v, want %v", args, want)
	}
}

func TestCollectServerArgs_NoFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("port", "9080", "")
	fs.Parse([]string{})

	args := collectServerArgs(fs.Visit)
	if args != nil {
		t.Errorf("got %v, want nil", args)
	}
}

func TestCollectServerArgs_SkipsNoBrowser(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("no-browser", false, "")
	fs.String("port", "9080", "")
	fs.Parse([]string{"-no-browser", "-port=8080"})

	args := collectServerArgs(fs.Visit)
	want := []string{"-port=8080"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("got %v, want %v", args, want)
	}
}

func TestOpenFile_NoBrowser(t *testing.T) {
	// Not running: the caller must start the server in the foreground.
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()

	cfg := &config.Config{Port: port, NoBrowser: true}
	runServer, err := openFile(cfg, "README.md")
	if err != nil || !runServer {
		t.Errorf("not running: got (%v, %v), want (true, nil)", runServer, err)
	}

	// Already running: nothing to start, just print the URL.
	l, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	cfg.Port = l.Addr().(*net.TCPAddr).Port
	runServer, err = openFile(cfg, "README.md")
	if err != nil || runServer {
		t.Errorf("running: got (%v, %v), want (false, nil)", runServer, err)
	}
}
