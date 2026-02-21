package credential

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GetToken retrieves an authentication token for the given host using git credential helper.
// Returns username and password/token, or empty strings if not available.
// Uses a timeout to prevent hanging if the credential helper prompts interactively.
func GetToken(host string) (username, password string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "credential", "fill")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("protocol=https\nhost=%s\n\n", host))

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
