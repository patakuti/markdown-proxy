package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/patakuti/markdown-proxy/internal/config"
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
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
