package themes

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

var builtinCSS map[string]string

func init() {
	githubHL := scopeChromaTokens(generateSyntaxCSS("github"))
	monokaiHL := scopeChromaTokens(generateSyntaxCSS("monokai"))
	builtinCSS = map[string]string{
		"github": buildGithubCSS(githubHL),
		"simple": buildSimpleCSS(githubHL),
		"dark":   buildDarkCSS(monokaiHL),
	}
}

func generateSyntaxCSS(styleName string) string {
	style := styles.Get(styleName)
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var buf bytes.Buffer
	if err := formatter.WriteCSS(&buf, style); err != nil {
		return ""
	}
	return buf.String()
}

// scopeChromaTokens prefixes chroma token class selectors (e.g. .k, .s) with .chroma
// so they only apply inside code blocks. Lines already starting with .chroma are left alone.
func scopeChromaTokens(css string) string {
	var result strings.Builder
	for _, line := range strings.Split(css, "\n") {
		if line == "" {
			result.WriteString("\n")
			continue
		}
		if strings.Contains(line, "{") {
			idx := strings.Index(line, ".")
			if idx >= 0 && !strings.HasPrefix(line[idx:], ".chroma") {
				result.WriteString(line[:idx])
				result.WriteString(".chroma ")
				result.WriteString(line[idx:])
				result.WriteString("\n")
				continue
			}
		}
		result.WriteString(line)
		result.WriteString("\n")
	}
	return result.String()
}

// DefaultThemesDir returns the path to the user CSS themes directory.
func DefaultThemesDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "markdown-proxy", "themes"), nil
}

// EnsureBuiltinThemes writes built-in theme CSS files to dir if they don't already exist.
func EnsureBuiltinThemes(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating themes directory: %w", err)
	}
	for name, css := range builtinCSS {
		path := filepath.Join(dir, name+".css")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(css), 0644); err != nil {
				return fmt.Errorf("writing theme %s: %w", name, err)
			}
		}
	}
	return nil
}

// ListThemes returns sorted theme names (without .css extension) from dir.
// Falls back to built-in names if dir is empty or unreadable.
func ListThemes(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err == nil {
		var names []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".css") {
				names = append(names, strings.TrimSuffix(e.Name(), ".css"))
			}
		}
		if len(names) > 0 {
			sort.Strings(names)
			return names
		}
	}
	return []string{"dark", "github", "simple"}
}

// BuiltinCSS returns the CSS bytes for a built-in theme, or nil if not found.
// Used as a fallback when theme files cannot be read from disk.
func BuiltinCSS(name string) []byte {
	if css, ok := builtinCSS[name]; ok {
		return []byte(css)
	}
	return nil
}

// ReadTheme returns the CSS content for a named theme from dir.
func ReadTheme(dir, name string) ([]byte, error) {
	if !isValidName(name) {
		return nil, fmt.Errorf("invalid theme name: %q", name)
	}
	return os.ReadFile(filepath.Join(dir, name+".css"))
}

// ValidateTheme returns name if it exists in available, otherwise the first element or "github".
func ValidateTheme(name string, available []string) string {
	for _, t := range available {
		if t == name {
			return name
		}
	}
	if len(available) > 0 {
		return available[0]
	}
	return "github"
}

func isValidName(name string) bool {
	return name != "" && !strings.ContainsAny(name, "/\\.")
}

func buildGithubCSS(highlightCSS string) string {
	return `/* markdown-proxy built-in theme: github */
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: #24292e;
  background: #fff;
}
.toolbar { background: #f6f8fa; border-color: #e1e4e8; }
.home-link { color: #0366d6; }
.toolbar-link { color: #0366d6; }
.markdown-body h1 { padding-bottom: .3em; border-bottom: 1px solid #eaecef; }
.markdown-body h2 { padding-bottom: .3em; border-bottom: 1px solid #eaecef; }
.markdown-body a { color: #0366d6; text-decoration: none; }
.markdown-body a:hover { text-decoration: underline; }
.markdown-body code { background: rgba(27,31,35,.05); padding: .2em .4em; border-radius: 3px; font-size: 85%; }
.markdown-body pre { background: #f6f8fa; padding: 16px; border-radius: 6px; border: 1px solid #e1e4e8; overflow: auto; }
.markdown-body pre code { background: none; padding: 0; font-size: 100%; }
.markdown-body blockquote { color: #6a737d; border-left: .25em solid #dfe2e5; padding: 0 1em; margin: 0; }
.markdown-body img { max-width: 100%; }
.markdown-body table th, .markdown-body table td { border-color: #dfe2e5; }
.markdown-body table tr:nth-child(2n) { background-color: #f6f8fa; }
.plantuml-notice { background: #fff8c5; border-color: #d4a72c; color: #4d3800; }
.plantuml-notice code { background: rgba(0,0,0,.08); }
.line-highlight { background-color: rgba(255,255,0,0.2) !important; border-left-color: #f0c040 !important; }
.copy-btn { background: #fff; color: #444; border-color: #d0d7de; }
.copy-btn:hover { background: #f3f4f6; border-color: #adb5bd; }
.copy-btn.copied { background: #d4edda; color: #155724; border-color: #c3e6cb; }
.toc-panel { background: #f6f8fa; border-left-color: #e1e4e8; }
.toc-header { border-bottom-color: #e1e4e8; }
.toc-list a:hover { background: rgba(0,0,0,0.05); }
.toc-list a.active { border-left-color: #0366d6; background: rgba(3,102,214,0.08); }
.toc-list .toc-caret { color: #6a737d; }
` + highlightCSS
}

func buildSimpleCSS(highlightCSS string) string {
	return `/* markdown-proxy built-in theme: simple */
body {
  font-family: Georgia, "Times New Roman", serif;
  color: #333;
  background: #fefefe;
  line-height: 1.8;
}
.toolbar { background: #f6f8fa; border-color: #e1e4e8; }
.home-link { color: #07c; }
.toolbar-link { color: #07c; }
.markdown-body a { color: #07c; }
.markdown-body code { background: #f0f0f0; padding: .15em .3em; border-radius: 2px; font-size: 85%; }
.markdown-body pre { background: #f0f0f0; padding: 14px; border-radius: 4px; border: 1px solid #ddd; overflow: auto; }
.markdown-body pre code { background: none; padding: 0; font-size: 100%; }
.markdown-body blockquote { color: #666; border-left: 3px solid #ccc; padding: 0 1em; margin: 0; }
.markdown-body img { max-width: 100%; }
.markdown-body table th, .markdown-body table td { border-color: #dfe2e5; }
.markdown-body table tr:nth-child(2n) { background-color: #f6f8fa; }
.plantuml-notice { background: #fff8c5; border-color: #d4a72c; color: #4d3800; }
.plantuml-notice code { background: rgba(0,0,0,.08); }
.line-highlight { background-color: rgba(255,255,0,0.2) !important; border-left-color: #f0c040 !important; }
.copy-btn { background: #fff; color: #444; border-color: #d0d7de; }
.copy-btn:hover { background: #f3f4f6; border-color: #adb5bd; }
.copy-btn.copied { background: #d4edda; color: #155724; border-color: #c3e6cb; }
.toc-panel { background: #f5f5f5; border-left-color: #ddd; }
.toc-header { border-bottom-color: #ddd; }
.toc-list a:hover { background: rgba(0,0,0,0.05); }
.toc-list a.active { border-left-color: #07c; background: rgba(0,119,204,0.08); }
.toc-list .toc-caret { color: #666; }
` + highlightCSS
}

func buildDarkCSS(highlightCSS string) string {
	return `/* markdown-proxy built-in theme: dark */
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: #c9d1d9;
  background: #0d1117;
}
.toolbar { background: #1e1e1e; border-color: #444; }
.home-link { color: #58a6ff; }
.toolbar-link { color: #58a6ff; }
.markdown-body h1 { padding-bottom: .3em; border-bottom: 1px solid #21262d; }
.markdown-body h2 { padding-bottom: .3em; border-bottom: 1px solid #21262d; }
.markdown-body a { color: #58a6ff; text-decoration: none; }
.markdown-body a:hover { text-decoration: underline; }
.markdown-body code { background: rgba(110,118,129,.4); padding: .2em .4em; border-radius: 3px; font-size: 85%; }
.markdown-body pre { background: #161b22; padding: 16px; border-radius: 6px; border: 1px solid #30363d; overflow: auto; }
.markdown-body pre code { background: none; padding: 0; font-size: 100%; }
.markdown-body blockquote { color: #8b949e; border-left: .25em solid #30363d; padding: 0 1em; margin: 0; }
.markdown-body img { max-width: 100%; }
.markdown-body table th, .markdown-body table td { border-color: #444; }
.markdown-body table tr:nth-child(2n) { background-color: #2d2d2d; }
.plantuml-notice { background: #2d2a1e; border-color: #966c00; color: #e3b341; }
.plantuml-notice code { background: rgba(255,255,255,.1); }
.line-highlight { background-color: rgba(255,255,0,0.1) !important; border-left-color: #b08820 !important; }
.copy-btn { background: #21262d; color: #c9d1d9; border-color: #30363d; }
.copy-btn:hover { background: #2d333b; border-color: #6e7681; }
.copy-btn.copied { background: #1a3d2b; color: #3fb950; border-color: #2ea043; }
.toc-panel { background: #161b22; border-left-color: #30363d; }
.toc-header { border-bottom-color: #30363d; }
.toc-list a:hover { background: rgba(255,255,255,0.08); }
.toc-list a.active { border-left-color: #58a6ff; background: rgba(88,166,255,0.12); }
.toc-list .toc-caret { color: #8b949e; }
` + highlightCSS
}
