package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/opener"
	"github.com/patakuti/markdown-proxy/internal/server"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version and exit")
	cfg := config.Parse()
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if cfg.Configure {
		if err := config.RunConfigure(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
		return
	}

	// If a file or URL argument is provided, open it in the browser.
	if args := flag.Args(); len(args) > 0 {
		runServer, err := openFile(cfg, args[0])
		if err != nil {
			log.Fatal(err)
		}
		if !runServer {
			return
		}
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}

// openFile builds the proxy URL for the given argument, ensures the server
// is running, and opens the URL in the default browser.
// With -no-browser it only prints the URL; it returns true when the caller
// should then run the server in the foreground (no server is running yet).
func openFile(cfg *config.Config, arg string) (bool, error) {
	proxyURL, err := opener.BuildURL(arg, cfg.Port)
	if err != nil {
		return false, err
	}

	if cfg.NoBrowser {
		log.Printf("URL: %s", proxyURL)
		return !opener.IsServerRunning(cfg.Port), nil
	}

	if !opener.IsServerRunning(cfg.Port) {
		log.Printf("Starting server on port %d...", cfg.Port)
		serverArgs := buildServerArgs()
		if err := opener.StartServer(cfg.Port, serverArgs); err != nil {
			return false, err
		}
	}

	log.Printf("Opening %s", proxyURL)
	return false, opener.OpenBrowser(proxyURL)
}

// buildServerArgs collects the explicitly specified command-line flags
// to pass to the background server process.
func buildServerArgs() []string {
	return collectServerArgs(flag.Visit)
}

// collectServerArgs extracts flags using the provided visit function.
// This is separated from buildServerArgs to allow testing.
func collectServerArgs(visit func(func(*flag.Flag))) []string {
	var args []string
	visit(func(f *flag.Flag) {
		// Skip flags that are not relevant for the server process.
		if f.Name == "version" || f.Name == "no-browser" {
			return
		}
		args = append(args, fmt.Sprintf("-%s=%s", f.Name, f.Value.String()))
	})
	return args
}
