package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/handler"
)

func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	topHandler := handler.NewTopHandler(cfg)
	localHandler := handler.NewLocalHandler(cfg)
	remoteHandler := handler.NewRemoteHandler(cfg)

	mux.HandleFunc("/", topHandler.ServeHTTP)
	mux.HandleFunc("/local/", localHandler.ServeHTTP)
	mux.HandleFunc("/http/", remoteHandler.ServeHTTP)
	mux.HandleFunc("/https/", remoteHandler.ServeHTTP)

	addr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)
	log.Printf("Starting markdown-proxy on %s", addr)
	return http.ListenAndServe(addr, mux)
}
