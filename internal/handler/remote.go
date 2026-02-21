package handler

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/patakuti/markdown-proxy/internal/config"
	"github.com/patakuti/markdown-proxy/internal/credential"
	ghub "github.com/patakuti/markdown-proxy/internal/github"
	"github.com/patakuti/markdown-proxy/internal/markdown"
	tmpl "github.com/patakuti/markdown-proxy/internal/template"
)

type RemoteHandler struct {
	cfg    *config.Config
	client *http.Client
}

func NewRemoteHandler(cfg *config.Config, client *http.Client) *RemoteHandler {
	return &RemoteHandler{cfg: cfg, client: client}
}

func (h *RemoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Determine scheme from URL path prefix
	var scheme string
	var remotePath string
	if strings.HasPrefix(r.URL.Path, "/https/") {
		scheme = "https"
		remotePath = strings.TrimPrefix(r.URL.Path, "/https/")
	} else if strings.HasPrefix(r.URL.Path, "/http/") {
		scheme = "http"
		remotePath = strings.TrimPrefix(r.URL.Path, "/http/")
	} else {
		http.Error(w, "Invalid URL scheme", http.StatusBadRequest)
		return
	}

	if remotePath == "" {
		http.Error(w, "No remote path specified", http.StatusBadRequest)
		return
	}

	// Check if this is a repo root URL (e.g., github.com/user/repo)
	if candidates := ghub.ResolveRepoRootURLs(remotePath); candidates != nil {
		for _, candidate := range candidates {
			candidateURL := scheme + "://" + candidate
			body, contentType, err := h.fetchRemote(candidateURL, remotePath)
			if err == nil {
				h.renderMarkdownResponse(w, body, contentType, remotePath, scheme)
				return
			}
		}
		http.Error(w, "Could not find README.md in repository", http.StatusNotFound)
		return
	}

	// Resolve GitHub/GitLab blob URLs to raw URLs
	fetchPath := remotePath
	if resolved, ok := ghub.ResolveRawURL(remotePath); ok {
		fetchPath = resolved
	}

	// Build the remote URL
	remoteURL := scheme + "://" + fetchPath

	// Fetch the remote content
	// Pass remotePath for credential lookup (use original host, not resolved raw host)
	body, contentType, err := h.fetchRemote(remoteURL, remotePath)
	if err != nil {
		http.Error(w, "Error fetching remote file: "+err.Error(), http.StatusBadGateway)
		return
	}

	h.renderResponse(w, body, contentType, fetchPath, remotePath, scheme)
}

// renderMarkdownResponse renders body as Markdown HTML (used for repo root README.md).
func (h *RemoteHandler) renderMarkdownResponse(w http.ResponseWriter, body []byte, contentType, remotePath, scheme string) {
	htmlContent, err := markdown.Convert(body, h.cfg.PlantUMLServer)
	if err != nil {
		http.Error(w, "Error converting markdown: "+err.Error(), http.StatusInternalServerError)
		return
	}

	server := ghub.HostFromPath(remotePath)
	htmlContent = markdown.RewriteLinks(htmlContent, scheme, server)

	page, err := tmpl.RenderMarkdown(&tmpl.PageData{
		Title:   path.Base(remotePath) + " - README.md",
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

// renderResponse renders body based on file extension.
func (h *RemoteHandler) renderResponse(w http.ResponseWriter, body []byte, contentType, fetchPath, remotePath, scheme string) {
	ext := strings.ToLower(path.Ext(fetchPath))

	// Non-markdown files: pass through
	if ext != ".md" && ext != ".markdown" {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.Write(body)
		return
	}

	// Markdown file: convert to HTML
	htmlContent, err := markdown.Convert(body, h.cfg.PlantUMLServer)
	if err != nil {
		http.Error(w, "Error converting markdown: "+err.Error(), http.StatusInternalServerError)
		return
	}

	server := ghub.HostFromPath(remotePath)
	htmlContent = markdown.RewriteLinks(htmlContent, scheme, server)

	page, err := tmpl.RenderMarkdown(&tmpl.PageData{
		Title:   path.Base(remotePath),
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

func (h *RemoteHandler) fetchRemote(remoteURL, remotePath string) ([]byte, string, error) {
	// First attempt: without authentication
	body, contentType, err := h.doFetch(remoteURL, "", "")
	if err == nil {
		return body, contentType, nil
	}

	// If 401/403/404, retry with git credential helper
	// Note: GitHub returns 404 for private repos when unauthenticated
	// Use original remotePath host (e.g. github.com) for credential lookup,
	// not the resolved raw host (e.g. raw.githubusercontent.com)
	if httpErr, ok := err.(*httpError); ok && (httpErr.StatusCode == 401 || httpErr.StatusCode == 403 || httpErr.StatusCode == 404) {
		host := ghub.HostFromPath(remotePath)
		log.Printf("Got %d for %s, trying git credential for host=%s", httpErr.StatusCode, remoteURL, host)
		username, password, credErr := credential.GetToken(host)
		if credErr != nil {
			log.Printf("Warning: git credential failed for %s: %v", host, credErr)
			return nil, "", err // return original HTTP error
		}
		if password != "" {
			if username == "" {
				username = "oauth2"
			}
			log.Printf("Retrying with credential (user=%s) for %s", username, remoteURL)
			return h.doFetch(remoteURL, username, password)
		}
		log.Printf("Warning: git credential returned empty password for %s", host)
	}

	return nil, "", err
}

func (h *RemoteHandler) doFetch(remoteURL, username, password string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", remoteURL, nil)
	if err != nil {
		return nil, "", err
	}

	if password != "" {
		// GitHub/GitLab tokens work with "token <PAT>" authorization header
		req.Header.Set("Authorization", "token "+password)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", &httpError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return body, resp.Header.Get("Content-Type"), nil
}

type httpError struct {
	StatusCode int
	Status     string
}

func (e *httpError) Error() string {
	return "remote server returned " + e.Status
}
