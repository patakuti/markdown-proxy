package github

import (
	"regexp"
	"strings"
)

// GitHub blob URL pattern: github.com/user/repo/blob/<ref>/path
var githubBlobRe = regexp.MustCompile(`^github\.com/([^/]+)/([^/]+)/blob/([^/]+)/(.+)$`)

// GitHub repo root pattern: github.com/user/repo (with optional trailing slash)
var githubRepoRe = regexp.MustCompile(`^github\.com/([^/]+)/([^/]+?)/?$`)

// GitLab blob URL pattern: gitlab.com/user/repo/-/blob/<ref>/path
var gitlabBlobRe = regexp.MustCompile(`^gitlab\.com/([^/]+)/([^/]+)/-/blob/([^/]+)/(.+)$`)

// GitLab repo root pattern: gitlab.com/user/repo (with optional trailing slash)
var gitlabRepoRe = regexp.MustCompile(`^gitlab\.com/([^/]+)/([^/]+?)/?$`)

// ResolveRawURL converts a GitHub/GitLab blob URL path to a raw content URL.
// path is the host + path portion (e.g., "github.com/user/repo/blob/main/README.md").
// Returns the raw URL and true if conversion was performed, or the original path and false.
func ResolveRawURL(path string) (string, bool) {
	// GitHub: github.com/user/repo/blob/<ref>/path → raw.githubusercontent.com/user/repo/<ref>/path
	if m := githubBlobRe.FindStringSubmatch(path); m != nil {
		user, repo, ref, filePath := m[1], m[2], m[3], m[4]
		rawURL := "raw.githubusercontent.com/" + user + "/" + repo + "/" + ref + "/" + filePath
		return rawURL, true
	}

	// GitLab: gitlab.com/user/repo/-/blob/<ref>/path → gitlab.com/user/repo/-/raw/<ref>/path
	if m := gitlabBlobRe.FindStringSubmatch(path); m != nil {
		user, repo, ref, filePath := m[1], m[2], m[3], m[4]
		rawURL := "gitlab.com/" + user + "/" + repo + "/-/raw/" + ref + "/" + filePath
		return rawURL, true
	}

	return path, false
}

// ResolveRepoRootURLs returns candidate raw URLs for a repository root.
// For repo root URLs (github.com/user/repo), returns README.md URLs for main and master branches.
// Returns nil if the path is not a repo root.
func ResolveRepoRootURLs(path string) []string {
	// GitHub: github.com/user/repo → try main, then master
	if m := githubRepoRe.FindStringSubmatch(path); m != nil {
		user, repo := m[1], m[2]
		return []string{
			"raw.githubusercontent.com/" + user + "/" + repo + "/main/README.md",
			"raw.githubusercontent.com/" + user + "/" + repo + "/master/README.md",
		}
	}

	// GitLab: gitlab.com/user/repo → try main, then master
	if m := gitlabRepoRe.FindStringSubmatch(path); m != nil {
		user, repo := m[1], m[2]
		return []string{
			"gitlab.com/" + user + "/" + repo + "/-/raw/main/README.md",
			"gitlab.com/" + user + "/" + repo + "/-/raw/master/README.md",
		}
	}

	return nil
}

// HostFromPath extracts the host portion from a URL path.
// e.g., "github.com/user/repo/blob/main/file.md" → "github.com"
func HostFromPath(path string) string {
	idx := strings.Index(path, "/")
	if idx < 0 {
		return path
	}
	return path[:idx]
}
