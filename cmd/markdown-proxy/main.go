package main

import (
	"log"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/server"
)

func main() {
	cfg := config.Parse()
	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
