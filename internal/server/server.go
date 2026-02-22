package server

import (
	"crypto/subtle"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/handler"
	"github.com/patakuti/markdown-proxy/internal/network"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Run(cfg *config.Config) error {
	mux := http.NewServeMux()

	// In local mode, allow private network access (user is local).
	// In remote mode, block private network access (SSRF prevention).
	client := network.NewSafeClient(!cfg.IsRemoteMode())

	topHandler := handler.NewTopHandler(cfg)
	remoteHandler := handler.NewRemoteHandler(cfg, client)

	mux.HandleFunc("/", topHandler.ServeHTTP)
	mux.HandleFunc("/http/", remoteHandler.ServeHTTP)
	mux.HandleFunc("/https/", remoteHandler.ServeHTTP)

	if cfg.IsRemoteMode() {
		// In remote mode, block local file access and SSE
		forbidden := func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Local file access is not available in remote mode", http.StatusForbidden)
		}
		mux.HandleFunc("/local/", forbidden)
		mux.HandleFunc("/_sse", forbidden)

		// Add login page handler
		loginHandler := handler.NewLoginHandler(cfg)
		mux.HandleFunc("/_login", loginHandler.ServeHTTP)
	} else {
		// In local mode, enable local file access and SSE
		localHandler := handler.NewLocalHandler(cfg)
		sseHandler := handler.NewSSEHandler()
		mux.HandleFunc("/local/", localHandler.ServeHTTP)
		mux.HandleFunc("/_sse", sseHandler.ServeHTTP)
	}

	// Build middleware chain (applied in reverse order)
	// Request flow: verbose -> accessLog -> auth -> mux
	var h http.Handler = mux

	if cfg.AuthToken != "" {
		h = authMiddleware(h, cfg.AuthToken)
	}

	if w := newAccessLogWriter(cfg); w != nil {
		h = accessLogMiddleware(h, w)
	}

	if cfg.Verbose {
		h = verboseMiddleware(h)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Listen, cfg.Port)

	if cfg.IsRemoteMode() {
		log.Printf("WARNING: Starting markdown-proxy in REMOTE mode on %s", addr)
		log.Printf("  Local file access is disabled")
		log.Printf("  Authentication is enabled")
	} else {
		log.Printf("Starting markdown-proxy on %s", addr)
	}

	return http.ListenAndServe(addr, h)
}

// responseWriter wraps http.ResponseWriter to capture status code and
// response size for logging.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// verboseMiddleware logs requests to stderr for debugging (--verbose).
func verboseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

// newAccessLogWriter returns the io.Writer for access log output.
// Returns nil if access logging should be disabled.
func newAccessLogWriter(cfg *config.Config) io.Writer {
	if cfg.AccessLog != "" {
		return &lumberjack.Logger{
			Filename:   cfg.AccessLog,
			MaxSize:    cfg.AccessLogMaxSize,
			MaxBackups: cfg.AccessLogMaxBack,
			MaxAge:     cfg.AccessLogMaxAge,
		}
	}
	if cfg.IsRemoteMode() {
		return os.Stdout
	}
	return nil
}

// authMiddleware checks for a valid authentication token in the request cookie.
// Requests to /_login are allowed through without authentication.
func authMiddleware(next http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/_login") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(handler.CookieName)
		if err == nil && subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(token)) == 1 {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/_login", http.StatusSeeOther)
	})
}

// accessLogMiddleware logs each request in a structured format.
func accessLogMiddleware(next http.Handler, w io.Writer) http.Handler {
	logger := log.New(w, "", 0)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: rw, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if remoteIP == "" {
			remoteIP = r.RemoteAddr
		}

		logger.Printf("%s %s %s %s %d %d %s",
			start.Format(time.RFC3339),
			remoteIP,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			wrapped.size,
			time.Since(start).Round(time.Millisecond),
		)
	})
}
