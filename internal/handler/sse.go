package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type SSEHandler struct{}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	watchPath := r.URL.Query().Get("path")
	if watchPath == "" {
		http.Error(w, "missing path parameter", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(watchPath)
	if err != nil {
		http.Error(w, "path not found: "+watchPath, http.StatusNotFound)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		http.Error(w, "failed to create watcher: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer watcher.Close()

	if info.IsDir() {
		err = watcher.Add(watchPath)
	} else {
		// Watch the parent directory to catch file renames/recreations
		err = watcher.Add(filepath.Dir(watchPath))
	}
	if err != nil {
		http.Error(w, "failed to watch path: "+err.Error(), http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	ctx := r.Context()
	debounceCh := make(chan struct{}, 1)
	var debounceTimer *time.Timer

	for {
		select {
		case <-ctx.Done():
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// For files, only react to events for the target file
			if !info.IsDir() && filepath.Clean(event.Name) != filepath.Clean(watchPath) {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}
			// Debounce: reset timer on each event, send after 100ms of quiet
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
				select {
				case debounceCh <- struct{}{}:
				default:
				}
			})
		case <-debounceCh:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}
