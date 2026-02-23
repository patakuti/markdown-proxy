package github

import (
	"testing"
)

func TestResolveRawURL(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantOK  bool
	}{
		// GitHub blob URL
		{
			name:   "GitHub basic blob URL",
			path:   "github.com/user/repo/blob/main/README.md",
			want:   "raw.githubusercontent.com/user/repo/main/README.md",
			wantOK: true,
		},
		{
			name:   "GitHub nested path",
			path:   "github.com/user/repo/blob/main/docs/guide/intro.md",
			want:   "raw.githubusercontent.com/user/repo/main/docs/guide/intro.md",
			wantOK: true,
		},
		{
			name:   "GitHub tag ref",
			path:   "github.com/user/repo/blob/v1.2.3/file.md",
			want:   "raw.githubusercontent.com/user/repo/v1.2.3/file.md",
			wantOK: true,
		},
		{
			name:   "GitHub SHA ref",
			path:   "github.com/user/repo/blob/abc1234/file.md",
			want:   "raw.githubusercontent.com/user/repo/abc1234/file.md",
			wantOK: true,
		},
		// GitLab blob URL
		{
			name:   "GitLab basic blob URL",
			path:   "gitlab.com/user/repo/-/blob/main/README.md",
			want:   "gitlab.com/user/repo/-/raw/main/README.md",
			wantOK: true,
		},
		{
			name:   "GitLab subgroup",
			path:   "gitlab.com/group/subgroup/repo/-/blob/main/file.md",
			want:   "gitlab.com/group/subgroup/repo/-/raw/main/file.md",
			wantOK: true,
		},
		{
			name:   "GitLab custom domain",
			path:   "gitlab.example.com/team/project/-/blob/develop/docs/api.md",
			want:   "gitlab.example.com/team/project/-/raw/develop/docs/api.md",
			wantOK: true,
		},
		// Non-matching paths
		{
			name:   "GitHub tree URL (not blob)",
			path:   "github.com/user/repo/tree/main/docs",
			want:   "github.com/user/repo/tree/main/docs",
			wantOK: false,
		},
		{
			name:   "GitHub repo root",
			path:   "github.com/user/repo",
			want:   "github.com/user/repo",
			wantOK: false,
		},
		{
			name:   "empty string",
			path:   "",
			want:   "",
			wantOK: false,
		},
		{
			name:   "random path",
			path:   "example.com/some/page",
			want:   "example.com/some/page",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ResolveRawURL(tt.path)
			if ok != tt.wantOK {
				t.Errorf("ResolveRawURL(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("ResolveRawURL(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestResolveRepoRootURLs(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "GitHub repo root",
			path: "github.com/user/repo",
			want: []string{
				"raw.githubusercontent.com/user/repo/main/README.md",
				"raw.githubusercontent.com/user/repo/master/README.md",
			},
		},
		{
			name: "GitHub repo root with trailing slash",
			path: "github.com/user/repo/",
			want: []string{
				"raw.githubusercontent.com/user/repo/main/README.md",
				"raw.githubusercontent.com/user/repo/master/README.md",
			},
		},
		{
			name: "GitLab repo root",
			path: "gitlab.com/user/repo",
			want: []string{
				"gitlab.com/user/repo/-/raw/main/README.md",
				"gitlab.com/user/repo/-/raw/master/README.md",
			},
		},
		{
			name: "GitLab repo root with trailing slash",
			path: "gitlab.com/user/repo/",
			want: []string{
				"gitlab.com/user/repo/-/raw/main/README.md",
				"gitlab.com/user/repo/-/raw/master/README.md",
			},
		},
		{
			name: "non-matching path",
			path: "github.com/user/repo/blob/main/file.md",
			want: nil,
		},
		{
			name: "empty string",
			path: "",
			want: nil,
		},
		{
			name: "random domain",
			path: "example.com/user/repo",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveRepoRootURLs(tt.path)
			if tt.want == nil {
				if got != nil {
					t.Errorf("ResolveRepoRootURLs(%q) = %v, want nil", tt.path, got)
				}
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ResolveRepoRootURLs(%q) returned %d URLs, want %d", tt.path, len(got), len(tt.want))
				return
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveRepoRootURLs(%q)[%d] = %q, want %q", tt.path, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestHostFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"with path", "github.com/user/repo/blob/main/file.md", "github.com"},
		{"host only", "github.com", "github.com"},
		{"empty string", "", ""},
		{"deep path", "gitlab.example.com/a/b/c/d", "gitlab.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HostFromPath(tt.path); got != tt.want {
				t.Errorf("HostFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestPathFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"with path", "github.com/user/repo/blob/main/file.md", "user/repo/blob/main/file.md"},
		{"host only", "github.com", ""},
		{"empty string", "", ""},
		{"single segment after host", "github.com/user", "user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathFromPath(tt.path); got != tt.want {
				t.Errorf("PathFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
