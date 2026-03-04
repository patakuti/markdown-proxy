package handler

import (
	"errors"
	"fmt"
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
	var remoteEscaped string
	if strings.HasPrefix(r.URL.Path, "/https/") {
		scheme = "https"
		remotePath = strings.TrimPrefix(r.URL.Path, "/https/")
		remoteEscaped = strings.TrimPrefix(r.URL.EscapedPath(), "/https/")
	} else if strings.HasPrefix(r.URL.Path, "/http/") {
		scheme = "http"
		remotePath = strings.TrimPrefix(r.URL.Path, "/http/")
		remoteEscaped = strings.TrimPrefix(r.URL.EscapedPath(), "/http/")
	} else {
		http.Error(w, "Invalid URL scheme", http.StatusBadRequest)
		return
	}

	if remotePath == "" {
		http.Error(w, "No remote path specified", http.StatusBadRequest)
		return
	}

	// Check if this is a repo root URL (e.g., github.com/user/repo)
	if candidates := ghub.ResolveRepoRootURLs(remoteEscaped); candidates != nil {
		var lastAuthErr *authError
		for _, candidate := range candidates {
			candidateURL := scheme + "://" + candidate
			body, contentType, err := h.fetchRemote(candidateURL, remotePath)
			if err == nil {
				h.renderMarkdownResponse(w, body, contentType, remotePath, scheme)
				return
			}
			var ae *authError
			if errors.As(err, &ae) {
				lastAuthErr = ae
			}
		}
		if lastAuthErr != nil {
			h.renderAuthError(w, lastAuthErr)
			return
		}
		http.Error(w, "Could not find README.md in repository", http.StatusNotFound)
		return
	}

	// Resolve GitHub/GitLab blob URLs to raw URLs
	fetchPath := remoteEscaped
	if resolved, ok := ghub.ResolveRawURL(remoteEscaped); ok {
		fetchPath = resolved
	}

	// Build the remote URL
	remoteURL := scheme + "://" + fetchPath

	// Fetch the remote content
	// Pass remotePath for credential lookup (use original host, not resolved raw host)
	body, contentType, err := h.fetchRemote(remoteURL, remotePath)
	if err != nil {
		var ae *authError
		if errors.As(err, &ae) {
			h.renderAuthError(w, ae)
			return
		}
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
	// Use remotePath (original blob URL) for extension detection, not fetchPath.
	// fetchPath may be a GitLab API URL (ending in /raw?ref=...) which has no file extension.
	ext := strings.ToLower(path.Ext(remotePath))

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
	// First attempt: without authentication, do not follow redirects.
	// Self-hosted GitLab redirects to login page (302) for unauthenticated access,
	// and Go's http.Client follows redirects by default, returning 200 (login HTML).
	// By disabling redirects, we get the 3xx status code and can trigger auth retry.
	body, contentType, err := h.doFetch(remoteURL, "", "", false)
	if err == nil {
		return body, contentType, nil
	}

	// If 3xx/401/403/404, retry with git credential helper
	// Note: GitHub returns 404 for private repos when unauthenticated
	// Note: Self-hosted GitLab returns 302 redirect to login for unauthenticated access
	// Use original remotePath host (e.g. github.com) for credential lookup,
	// not the resolved raw host (e.g. raw.githubusercontent.com)
	if httpErr, ok := err.(*httpError); ok && (httpErr.StatusCode == 401 || httpErr.StatusCode == 403 || httpErr.StatusCode == 404 || (httpErr.StatusCode >= 300 && httpErr.StatusCode <= 399)) {
		host := ghub.HostFromPath(remotePath)
		credPath := ghub.PathFromPath(remotePath)
		log.Printf("Got %d for %s, trying git credential for host=%s path=%s", httpErr.StatusCode, remoteURL, host, credPath)
		username, password, credErr := credential.GetToken(host, credPath)
		if credErr != nil {
			log.Printf("Warning: git credential failed for %s: %v", host, credErr)
			return nil, "", &authError{
				StatusCode: httpErr.StatusCode,
				Host:       host,
				Reason:     "credential_failed",
				Err:        credErr,
			}
		}
		if password != "" {
			if username == "" {
				username = "oauth2"
			}
			log.Printf("Retrying with credential (user=%s) for %s", username, remoteURL)
			// Authenticated retry follows redirects normally (for legitimate redirects)
			body, contentType, retryErr := h.doFetch(remoteURL, username, password, true)
			if retryErr != nil {
				return nil, "", &authError{
					StatusCode: httpErr.StatusCode,
					Host:       host,
					Reason:     "auth_rejected",
					Err:        retryErr,
				}
			}
			return body, contentType, nil
		}
		log.Printf("Warning: git credential returned empty password for %s", host)
		return nil, "", &authError{
			StatusCode: httpErr.StatusCode,
			Host:       host,
			Reason:     "credential_empty",
			Err:        err,
		}
	}

	return nil, "", err
}

func (h *RemoteHandler) doFetch(remoteURL, username, password string, followRedirects bool) ([]byte, string, error) {
	req, err := http.NewRequest("GET", remoteURL, nil)
	if err != nil {
		return nil, "", err
	}

	if password != "" {
		// Use "Bearer" format for OAuth2 tokens (glab auth login) and PATs.
		// GitHub also accepts "Bearer" as an alias for "token".
		req.Header.Set("Authorization", "Bearer "+password)
	}

	client := h.client
	if !followRedirects {
		// Shallow copy the client and disable redirect following.
		// This allows detecting 3xx responses (e.g., self-hosted GitLab
		// redirecting to login page) as auth-required signals.
		c := *h.client
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		client = &c
	}

	resp, err := client.Do(req)
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

// authError represents an authentication-related error when fetching remote files.
type authError struct {
	StatusCode int
	Host       string
	Reason     string // "credential_failed", "credential_empty", "auth_rejected"
	Err        error
}

func (e *authError) Error() string {
	return fmt.Sprintf("authentication failed for %s: %s", e.Host, e.Reason)
}

func (e *authError) Unwrap() error {
	return e.Err
}

// buildAuthHints generates user-facing hint messages based on the auth failure reason and host.
func buildAuthHints(host, reason string) []template.HTML {
	safeHost := template.HTMLEscapeString(host)

	var hints []template.HTML

	// Reason-specific message
	switch reason {
	case "credential_failed":
		hints = append(hints, "The git credential helper could not provide credentials. Make sure a credential helper is configured.")
	case "credential_empty":
		hints = append(hints, "The git credential helper returned an empty password. You may need to re-authenticate.")
	case "auth_rejected":
		hints = append(hints, "Credentials were provided but the server rejected them. The token may have expired or lack the required scope.")
	}

	// Host-specific setup hints
	switch {
	case host == "github.com" || strings.HasSuffix(host, ".github.com"):
		hints = append(hints,
			template.HTML("Run <code>gh auth login</code> to authenticate with GitHub CLI, or configure a Personal Access Token."),
		)
	case host == "gitlab.com" || strings.HasSuffix(host, ".gitlab.com"):
		hints = append(hints,
			template.HTML("Run <code>glab auth login</code> to authenticate with GitLab CLI, or configure a Personal Access Token with <code>read_repository</code> scope."),
		)
	default:
		hints = append(hints,
			template.HTML("Configure a git credential helper for <code>"+safeHost+"</code>. For self-hosted GitLab, create a Personal Access Token (<code>read_repository</code> scope) and run <code>glab auth login --hostname "+safeHost+" --token &lt;YOUR_TOKEN&gt;</code>."),
		)
	}

	// Common documentation link
	hints = append(hints,
		template.HTML(`See <a href="https://git-scm.com/docs/gitcredentials" target="_blank" rel="noopener">git-credentials documentation</a> for more details.`),
	)

	return hints
}

// renderAuthError renders an HTML error page for authentication failures.
func (h *RemoteHandler) renderAuthError(w http.ResponseWriter, ae *authError) {
	statusCode := ae.StatusCode
	// Map to user-friendly status: use 403 for auth errors regardless of the original status
	// (GitHub returns 404 for private repos, but the user-facing message should indicate access denial)
	if statusCode == 404 {
		statusCode = 403
	}

	message := fmt.Sprintf(
		"Could not access %s. This may be a private repository that requires authentication, or the resource may not exist.",
		template.HTMLEscapeString(ae.Host),
	)

	hints := buildAuthHints(ae.Host, ae.Reason)

	page, err := tmpl.RenderError(&tmpl.ErrorPageData{
		Title:   "Access Denied",
		Theme:   h.cfg.Theme,
		Status:  statusCode,
		Message: message,
		Hints:   hints,
	})
	if err != nil {
		log.Printf("Error rendering auth error page: %v", err)
		http.Error(w, "Error fetching remote file: "+ae.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(page)
}
