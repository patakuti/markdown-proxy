package handler

import (
	"html/template"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/markdown"
	tmpl "github.com/patakuti/markdown-proxy/internal/template"
)

type LocalHandler struct {
	cfg *config.Config
}

func NewLocalHandler(cfg *config.Config) *LocalHandler {
	return &LocalHandler{cfg: cfg}
}

func (h *LocalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the file path from the URL: /local/path/to/file → /path/to/file
	filePath := strings.TrimPrefix(r.URL.Path, "/local")
	if filePath == "" {
		filePath = "/"
	}

	// Expand ~ to home directory
	if strings.HasPrefix(filePath, "/~/") {
		if home, err := os.UserHomeDir(); err == nil {
			filePath = home + filePath[2:]
		}
	} else if filePath == "/~" {
		if home, err := os.UserHomeDir(); err == nil {
			filePath = home
		}
	}

	// Clean the path to prevent directory traversal
	filePath = filepath.Clean(filePath)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found: "+filePath, http.StatusNotFound)
		} else {
			http.Error(w, "Error accessing file: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		h.serveDirectory(w, filePath)
		return
	}

	h.serveFile(w, filePath)
}

func (h *LocalHandler) serveFile(w http.ResponseWriter, filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	// Non-markdown files: serve as-is
	if ext != ".md" && ext != ".markdown" {
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = http.DetectContentType(data)
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
		return
	}

	// Markdown file: convert to HTML
	htmlContent, err := markdown.Convert(data, h.cfg.PlantUMLServer)
	if err != nil {
		http.Error(w, "Error converting markdown: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Rewrite links
	htmlContent = markdown.RewriteLinks(htmlContent, "local", "")

	page, err := tmpl.RenderMarkdown(&tmpl.PageData{
		Title:   filepath.Base(filePath),
		Content: template.HTML(htmlContent),
		Theme:   h.cfg.Theme,
	})
	if err != nil {
		http.Error(w, "Error rendering page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}

func (h *LocalHandler) serveDirectory(w http.ResponseWriter, dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Error reading directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var dirEntries []tmpl.DirEntry

	// Add parent directory link if not at root
	if dirPath != "/" {
		parent := filepath.Dir(dirPath)
		dirEntries = append(dirEntries, tmpl.DirEntry{
			Name:  "..",
			IsDir: true,
			URL:   "/local" + parent,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		// Directories first, then files
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		url := "/local" + filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			url += "/"
		}
		dirEntries = append(dirEntries, tmpl.DirEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			URL:   url,
		})
	}

	page, err := tmpl.RenderDirectory(&tmpl.DirPageData{
		Title:   dirPath,
		Path:    dirPath,
		Entries: dirEntries,
		Theme:   h.cfg.Theme,
	})
	if err != nil {
		http.Error(w, "Error rendering directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}
