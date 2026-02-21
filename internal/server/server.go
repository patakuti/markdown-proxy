package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/handler"
	"github.com/patakuti/markdown-proxy/internal/network"
)

func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	client := network.NewSafeClient(cfg.AllowPrivateNetwork)

	topHandler := handler.NewTopHandler(cfg)
	localHandler := handler.NewLocalHandler(cfg)
	remoteHandler := handler.NewRemoteHandler(cfg, client)

	mux.HandleFunc("/", topHandler.ServeHTTP)
	mux.HandleFunc("/local/", localHandler.ServeHTTP)
	mux.HandleFunc("/http/", remoteHandler.ServeHTTP)
	mux.HandleFunc("/https/", remoteHandler.ServeHTTP)

	var h http.Handler = mux
	if cfg.Verbose {
		h = loggingMiddleware(mux)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)
	log.Printf("Starting markdown-proxy on %s", addr)
	return http.ListenAndServe(addr, h)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}
