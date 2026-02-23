package credential

import (
	"os/exec"
	"testing"
	"time"
)

func TestParseCredentialPath(t *testing.T) {
	tests := []struct {
		name            string
		gitConfigOutput string
		host            string
		remotePath      string
		want            string
	}{
		{
			name:            "single match with .helper suffix",
			gitConfigOutput: "credential.https://github.com/org.helper store\n",
			host:            "github.com",
			remotePath:      "org/repo/blob/main/README.md",
			want:            "org",
		},
		{
			name:            "single match with .username suffix",
			gitConfigOutput: "credential.https://github.com/org.username myuser\n",
			host:            "github.com",
			remotePath:      "org/repo/blob/main/README.md",
			want:            "org",
		},
		{
			name:            "single match with .useHttpPath suffix",
			gitConfigOutput: "credential.https://github.com/org.useHttpPath true\n",
			host:            "github.com",
			remotePath:      "org/repo/blob/main/README.md",
			want:            "org",
		},
		{
			name: "longest match wins",
			gitConfigOutput: "credential.https://github.com/org.helper store\n" +
				"credential.https://github.com/org/repo.helper store\n",
			host:       "github.com",
			remotePath: "org/repo/blob/main/README.md",
			want:       "org/repo",
		},
		{
			name:            "host mismatch",
			gitConfigOutput: "credential.https://gitlab.com/org.helper store\n",
			host:            "github.com",
			remotePath:      "org/repo/blob/main/README.md",
			want:            "",
		},
		{
			name:            "empty git config output",
			gitConfigOutput: "",
			host:            "github.com",
			remotePath:      "org/repo/blob/main/README.md",
			want:            "",
		},
		{
			name:            "empty remote path",
			gitConfigOutput: "credential.https://github.com/org.helper store\n",
			host:            "github.com",
			remotePath:      "",
			want:            "",
		},
		{
			name:            "path boundary check - or does not match org",
			gitConfigOutput: "credential.https://github.com/or.helper store\n",
			host:            "github.com",
			remotePath:      "org/repo/blob/main/README.md",
			want:            "",
		},
		{
			name:            "exact path match",
			gitConfigOutput: "credential.https://github.com/org.helper store\n",
			host:            "github.com",
			remotePath:      "org",
			want:            "org",
		},
		{
			name: "multiple hosts mixed",
			gitConfigOutput: "credential.https://gitlab.com/team.helper store\n" +
				"credential.https://github.com/org.helper store\n" +
				"credential.https://github.com/other.helper store\n",
			host:       "github.com",
			remotePath: "org/repo/blob/main/README.md",
			want:       "org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCredentialPath(tt.gitConfigOutput, tt.host, tt.remotePath)
			if got != tt.want {
				t.Errorf("parseCredentialPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetToken_NoPanic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode (git credential fill may hang)")
	}

	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH, skipping")
	}

	// Run in a goroutine with a short deadline to avoid hanging
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Should not panic even for a non-existent host
		_, _, _ = GetToken("nonexistent.example.invalid", "")
	}()

	select {
	case <-done:
		// completed without panic
	case <-time.After(5 * time.Second):
		t.Skip("skipping: git credential fill timed out (likely waiting for interactive input)")
	}
}
