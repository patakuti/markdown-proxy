package credential

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// GetToken retrieves an authentication token for the given host using git credential helper.
// Returns username and password/token, or empty strings if not available.
// Uses a timeout to prevent hanging if the credential helper prompts interactively.
// If path is non-empty, it searches git config for a path-based credential section
// (e.g., [credential "https://github.com/org"]) and uses the matching path for lookup.
func GetToken(host, path string) (username, password string, err error) {
	// Find the matching credential path from git config
	credPath := findCredentialPath(host, path)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	args := []string{"credential", "fill"}
	input := fmt.Sprintf("protocol=https\nhost=%s\n", host)
	if credPath != "" {
		// Enable useHttpPath so git matches [credential "https://host/path"] sections
		args = []string{"-c", "credential.useHttpPath=true", "credential", "fill"}
		input += fmt.Sprintf("path=%s\n", credPath)
	}
	input += "\n"
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdin = strings.NewReader(input)

	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("git credential fill failed: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if k, v, ok := strings.Cut(line, "="); ok {
			switch k {
			case "username":
				username = v
			case "password":
				password = v
			}
		}
	}

	return username, password, nil
}

// findCredentialPath searches git config for path-based credential sections
// matching the given host and returns the best (longest) matching credential path.
// e.g., for host="github.com" and remotePath="gn-nhc/azores/blob/main/README.md",
// if [credential "https://github.com/gn-nhc"] exists, returns "gn-nhc".
// Returns empty string if no path-specific credential config is found.
func findCredentialPath(host, remotePath string) string {
	if remotePath == "" {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "config", "--get-regexp", `^credential\.`)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	bestMatch := parseCredentialPath(string(out), host, remotePath)

	if bestMatch != "" {
		log.Printf("[credential] found matching credential path: %s (from config for https://%s/%s)", bestMatch, host, bestMatch)
	}

	return bestMatch
}

// parseCredentialPath parses git config output for path-based credential sections
// matching the given host and returns the best (longest) matching credential path.
// gitConfigOutput is the raw output from "git config --get-regexp ^credential\.".
func parseCredentialPath(gitConfigOutput, host, remotePath string) string {
	prefix := "credential.https://" + host + "/"
	knownSuffixes := []string{".helper", ".username", ".useHttpPath"}
	bestMatch := ""

	scanner := bufio.NewScanner(strings.NewReader(gitConfigOutput))
	for scanner.Scan() {
		line := scanner.Text()
		key, _, _ := strings.Cut(line, " ")
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// key = "credential.https://github.com/gn-nhc.helper"
		// remainder = "gn-nhc.helper"
		remainder := strings.TrimPrefix(key, prefix)

		// Extract the credential path by removing the known config key suffix
		var credPath string
		for _, suffix := range knownSuffixes {
			if strings.HasSuffix(remainder, suffix) {
				credPath = strings.TrimSuffix(remainder, suffix)
				break
			}
		}
		if credPath == "" {
			continue
		}

		// Check if credPath matches remotePath (exact or prefix with path boundary)
		if remotePath == credPath || strings.HasPrefix(remotePath, credPath+"/") {
			if len(credPath) > len(bestMatch) {
				bestMatch = credPath
			}
		}
	}

	return bestMatch
}
