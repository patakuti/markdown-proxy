# markdown-proxy

An HTTP proxy server for viewing Markdown files in a browser. Supports both local and remote access modes.

![Demo](docs/demo.gif)

## Motivation

Existing Markdown viewing environments have several pain points:

- **Limited local preview**: Tools like VS Code's Markdown Preview Enhanced are widely used for local Markdown viewing, but they can only display one preview pane at a time, making it difficult to reference multiple documents simultaneously.
- **Incomplete diagram rendering on remote servers**: Platforms such as GitHub and GitLab often do not render PlantUML diagrams embedded in Markdown files, leaving them displayed as raw code blocks.

markdown-proxy solves these problems by rendering Markdown files — including PlantUML diagrams — as HTML in a standard browser, where you can freely open multiple tabs and windows.

## Comparison

How markdown-proxy compares to other Markdown viewing tools:

| Feature | markdown-proxy (this tool) | [Markdown Preview Enhanced][mpe] | [grip][grip] | [Madness][madness] | [mdserve][mdserve] |
|---------|:-:|:-:|:-:|:-:|:-:|
| Type | HTTP server | VS Code extension | HTTP server | HTTP server | HTTP server |
| Local file viewing | ✅ | ✅ | ✅ | ✅ | ✅ |
| Remote URL fetching | ✅ (GitHub/GitLab) | ❌ | ❌ | ❌ | ❌ |
| Multi-tab viewing | ✅ | △ ¹ | ✅ | ✅ | ✅ |
| Directory listing | ✅ | ❌ | ❌ | ✅ | ✅ |
| CLI open | ✅ (file/URL) | — | ✅ (file) | ❌ | ✅ (--open) |
| Live reload | ✅ | ✅ | ✅ | ❌ | ✅ |
| Mermaid diagrams | ✅ | ✅ | ❌ | ✅ | ✅ |
| PlantUML diagrams | ✅ | ✅ | ❌ | ❌ | ❌ |
| Math rendering | ✅ (KaTeX) | ✅ (KaTeX/MathJax) | ❌ | ❌ | ❌ |
| Code highlighting | ✅ | ✅ | ✅ | ✅ | ✅ |
| CSS themes | 3 built-in + user-defined | 15+ built-in | GitHub only | Customizable | 5 built-in |
| Full-text search | ❌ | ❌ | ❌ | ✅ | ❌ |
| Export (PDF, HTML) | △ ² | ✅ (PDF, HTML, Word) | ✅ (HTML) | ❌ | ❌ |
| Authentication | Token-based | — | — | HTTP Basic | ❌ |
| Works offline | ✅ | ✅ | ❌ | ✅ | ✅ |
| Runtime dependency | None (single binary) | VS Code | Python | Ruby | None (single binary) |

¹ One preview pane per editor group
² Browser print-to-PDF via toolbar Print link

[mpe]: https://marketplace.visualstudio.com/items?itemName=shd101wyy.markdown-preview-enhanced
[grip]: https://github.com/joeyespo/grip
[madness]: https://github.com/DannyBen/madness
[mdserve]: https://github.com/jfernandez/mdserve

## Features

- Render local and remote Markdown files as HTML
- Support for GFM (GitHub Flavored Markdown) with syntax highlighting
- Math rendering (`$...$` for inline, `$$...$$` for display) via KaTeX
- Code block rendering
  - SVG: inline SVG rendering from ```` ```svg ```` code blocks
  - Mermaid: client-side rendering via mermaid.js from ```` ```mermaid ```` code blocks
  - PlantUML: server-side rendering from ```` ```plantuml ```` code blocks (requires `--plantuml-server`)
    - Disabled by default because diagram content is sent to the specified server
    - When disabled, a hint is shown in place of each PlantUML block
    - To use the public server: `--plantuml-server https://www.plantuml.com/plantuml`
    - Use `--configure` to save the setting permanently
- GitHub/GitLab integration
  - Blob URL auto-conversion to raw URL (supports self-hosted GitLab with custom domains)
  - Authentication via git credential helper (supports path-based credential matching)
  - Redirect-based auth detection for self-hosted GitLab instances
- Toolbar actions
  - Print: browser print with clean filename (toolbar hidden in print output)
  - Source: link to original URL on remote server (remote pages only)
- CSS themes with dropdown switching — 3 built-in themes (GitHub, Simple, Dark) plus user-defined themes
- Table of contents sidebar: toggle `TOC` in the toolbar to open a right-side panel with auto-extracted headings; visibility persists per browser (localStorage)
- Live reload for local files (auto-refreshes browser on file changes)
- Directory listing for local files
- Line anchor links: `[text](foo.md:12)` or `<a href="foo.md:12">` links navigate to specific source lines with highlighting (Markdown and text files)
- Text file rendering: `.txt` files are displayed in HTML with line anchors, themes, and live reload
- Link rewriting for seamless proxy navigation (including `file:///` protocol conversion)
- Top page with smart input (auto-detects file path or URL)
- Recently opened file history (localStorage)
- Two operation modes: local mode and remote mode
- Token-based authentication for remote access
- Access logging with automatic log rotation
- Configuration file (`~/.config/markdown-proxy/config.json`) with interactive setup via `--configure`
- Single binary, no runtime dependencies

## Quick Start

```bash
# Install
go install github.com/patakuti/markdown-proxy/cmd/markdown-proxy@latest

# Open a file directly (starts server automatically if needed)
markdown-proxy README.md

# Open a remote URL
markdown-proxy https://github.com/user/repo

# Or start the server manually and browse to http://localhost:9080/
markdown-proxy
```

When a file or URL is given as an argument, markdown-proxy checks if the server is already running; if not, it starts one in the background. Then it opens the file in your default browser. Without arguments, it starts the server in the foreground as before.

## Use Cases

- **Reviewing multiple documents side by side**: Open several Markdown files in separate browser tabs — no single-pane limitation like IDE preview plugins.
- **Viewing PlantUML/Mermaid diagrams**: Render diagrams embedded in Markdown that GitHub/GitLab don't display natively.
- **Browsing private repositories**: Access Markdown files from private GitHub/GitLab repos using your existing git credentials.
- **Navigating RAG search reports**: Reports generated by tools like [Local Knowledge RAG MCP Server](https://github.com/patakuti/local-knowledge-rag-mcp) contain links to source documents. With IDE preview plugins (single-pane), switching between the report and referenced documents is tedious. In markdown-proxy, open the report in one tab and click through references in new tabs — navigate freely between them.
- **Sharing a Markdown viewer with your team**: Run in remote mode with token authentication to let team members view documentation through a browser.

## URL Scheme

| Type | URL Format |
|------|-----------|
| Top page | `http://localhost:9080/` |
| Local file | `http://localhost:9080/local/path/to/file.md` |
| Local directory | `http://localhost:9080/local/path/to/dir/` |
| Remote (HTTP) | `http://localhost:9080/http/server/path/to/file.md` |
| Remote (HTTPS) | `http://localhost:9080/https/server/path/to/file.md` |
| GitHub repo | `http://localhost:9080/https/github.com/user/repo/blob/main/README.md` |

### Line Anchor Links

Markdown and text files support line-level linking using the `file:line` syntax:

- `[text](foo.md:12)` — links to line 12 of `foo.md`
- `[text](foo.md:12-34)` — links to lines 12–34 of `foo.md`
- `[text](foo.txt:12)` — links to line 12 of `foo.txt`
- `<a href="foo.md:12">text</a>` — same, using raw HTML

When navigating to a line anchor (`#L12` or `#L12-L34`), the page scrolls to the target line and highlights the surrounding content. For text files, individual lines are highlighted; for Markdown files, the containing block element is highlighted. Highlighting is hidden in print output.

## Usage

```bash
markdown-proxy [options] [file-or-url]
```

When `file-or-url` is provided:
- **Local file**: Opens the file via the proxy (relative paths are resolved automatically)
  - `markdown-proxy README.md`
  - `markdown-proxy ../docs/design.md`
  - `markdown-proxy /absolute/path/to/file.md`
- **Remote URL**: Opens the URL via the proxy
  - `markdown-proxy https://github.com/user/repo`
- If the server is not already running, it is started automatically in the background
- The file is opened in the default browser (`xdg-open` on Linux, `start` on Windows, `open` on macOS)

### Options

> **Note:** Options use a single dash (e.g., `-port`) following Go's `flag` package convention. Double dashes (`--port`) also work.

| Flag | Description | Default |
|------|-------------|---------|
| `-port`, `-p` | Listen port | `9080` |
| `-listen` | Bind address (`127.0.0.1` for local, `0.0.0.0` for remote) | `127.0.0.1` |
| `-theme` | Default CSS theme name (any file in the themes directory) | `github` |
| `-plantuml-server` | PlantUML server URL | (disabled) |
| `-auth-token` | Authentication token (required in remote mode) | |
| `-auth-cookie-max-age` | Authentication cookie max age in days | `30` |
| `-access-log` | Access log file path | |
| `-access-log-max-size` | Max log file size in MB before rotation | `100` |
| `-access-log-max-backups` | Max number of old log files to retain | `3` |
| `-access-log-max-age` | Max days to retain old log files | `28` |
| `-verbose`, `-v` | Enable debug logging to stderr | `false` |
| `-configure` | Interactively create configuration file | |
| `-version` | Show version and exit | |

## Configuration File

Settings can be saved to a configuration file so you don't need to pass flags every time.

### Interactive Setup

```bash
markdown-proxy --configure
```

This walks you through each setting interactively and saves the result.

### File Location

| Platform | Path |
|----------|------|
| Linux | `~/.config/markdown-proxy/config.json` |
| Windows | `%APPDATA%/markdown-proxy/config.json` |

### Supported Settings

The configuration file stores these settings as JSON:

```json
{
  "plantuml-server": "https://www.plantuml.com/plantuml",
  "theme": "dark",
  "port": 9080,
  "listen": "127.0.0.1"
}
```

Command-line flags override configuration file values. Security-sensitive settings (`-auth-token`, `-access-log`, etc.) are not stored in the config file.

## CSS Themes

markdown-proxy supports user-defined CSS themes in addition to the three built-in themes.

### Theme Directory

| Platform | Path |
|----------|------|
| Linux | `~/.config/markdown-proxy/themes/` |
| Windows | `%APPDATA%/markdown-proxy/themes/` |

On first launch, the three built-in themes are written here as editable CSS files:

- `github.css` — GitHub-style light theme
- `simple.css` — Serif light theme
- `dark.css` — Dark theme

### Adding a Custom Theme

1. Create a `.css` file in the themes directory, e.g. `~/.config/markdown-proxy/themes/mycolor.css`
2. Write standard CSS (no special scoping required — the file is loaded as the sole active stylesheet)
3. Restart markdown-proxy (or simply open a new page — the theme list is read at startup)
4. Select `mycolor` from the Theme dropdown

**Example:**

```css
body { font-family: "Source Serif Pro", serif; color: #1a1a1a; background: #fffff8; }
.toolbar { background: #f0ece0; border-color: #c8bfa0; }
.home-link, .toolbar-link { color: #8b4513; }
.markdown-body a { color: #8b4513; }
.markdown-body pre { background: #f5f0e8; border: 1px solid #c8bfa0; overflow: auto; }
.markdown-body pre code { background: none; padding: 0; font-size: 100%; }
.copy-btn { background: #fff; color: #444; border-color: #c8bfa0; }
.copy-btn:hover { background: #f0ece0; }
.copy-btn.copied { background: #d4edda; color: #155724; border-color: #c3e6cb; }
.toc-panel { background: #f5f0e8; border-left-color: #c8bfa0; }
.toc-header { border-bottom-color: #c8bfa0; }
.toc-list a.active { border-left-color: #8b4513; background: rgba(139,69,19,0.08); }
```

You can also edit the built-in themes directly — changes take effect on the next request (the file is read from disk on each theme CSS request).

## Operation Modes

### Local Mode (default)

```bash
markdown-proxy
```

The server binds to `127.0.0.1` and all features are available:
- Local file access (`/local/...`)
- Remote file access (`/http/...`, `/https/...`)
- Live reload via SSE (`/_sse`)
- Private network access is allowed for remote file fetching

### Remote Mode

```bash
markdown-proxy -listen 0.0.0.0 -auth-token my-secret-token
```

The server binds to the specified address for network access. For security:
- **Local file access is disabled**: `/local/` and `/_sse` return 403 Forbidden
- **Authentication is required**: `-auth-token` must be specified
- **Private network access is blocked**: SSRF protection prevents fetching from internal IPs
- **Top page** shows URL input only (no local file path input)

Users must authenticate via a login page (`/_login`) by entering the access token. The token is stored in an HttpOnly cookie for the configured duration.

## Access Logging

Access logs record each request in the following format:

```
2026-02-22T15:04:05+09:00 192.168.1.10 GET /https/github.com/user/repo 200 1234 150ms
```

- `-access-log /var/log/mdproxy/access.log`: Log to a file with automatic rotation
- In remote mode without `-access-log`: Logs to stdout by default
- In local mode without `-access-log`: No access logging

Log rotation is handled automatically using configurable size, count, and age limits.

## Live Reload

When viewing local Markdown files or directories (`/local/...`), the browser automatically reloads when the file or directory contents change. This uses Server-Sent Events (SSE) with filesystem notifications (fsnotify).

- **Local files only**: Remote files (`/http/...`, `/https/...`) are not affected
- **Local mode only**: Not available in remote mode
- **No configuration needed**: Works automatically for all local file views
- **Debounced**: Multiple rapid changes are coalesced into a single reload (100ms debounce)

## Math Rendering

Mathematical expressions are rendered using [KaTeX](https://katex.org/). Use standard LaTeX syntax:

- **Inline math**: `$E = mc^2$` renders inline within text
- **Display math**: `$$\int_0^\infty e^{-x} dx = 1$$` renders as a centered block

No configuration needed. Math expressions are automatically detected and rendered.

## Security

- **Local mode**: The server binds to `127.0.0.1` only, accepting local connections only
- **Remote mode**: Authentication is enforced via token. Local file access and private network fetching are automatically disabled
- **SSRF protection**: In remote mode, requests to private/internal IP addresses (e.g., `10.x.x.x`, `192.168.x.x`, `127.x.x.x`) are blocked. In local mode, private network access is allowed
- **DNS rebinding prevention**: Resolved IP addresses are used directly for connections, preventing DNS rebinding attacks
- **Constant-time token comparison**: Authentication uses `crypto/subtle.ConstantTimeCompare` to prevent timing attacks

## Private Repository Access

markdown-proxy can access private repositories on GitHub and GitLab using your existing git credentials. When authentication is needed (e.g., 401, 403, or 404 from a private repo), markdown-proxy automatically invokes `git credential fill` to retrieve stored credentials.

If credentials are not configured, an error page with setup guidance is displayed instead of a raw error message.

### GitHub

**Option A: GitHub CLI (recommended)**

```bash
gh auth login
```

This configures the git credential helper automatically.

**Option B: Personal Access Token**

1. Create a token at https://github.com/settings/tokens (classic) or https://github.com/settings/tokens?type=beta (fine-grained)
2. Required scope: `repo` (classic) or repository read access (fine-grained)
3. Store it via git credential helper:
   ```bash
   echo -e "protocol=https\nhost=github.com\n" | git credential fill
   ```
   If no credential is returned, configure a credential helper (e.g., `git config --global credential.helper store`) and save the token.

### GitLab

**gitlab.com:**

```bash
glab auth login
```

**Self-hosted GitLab:**

`glab auth login --hostname` may not work on some self-hosted instances. If authentication fails, use a Personal Access Token (PAT) instead:

1. Create a PAT at Settings > Access Tokens with `read_repository` scope
2. Register it with glab:
   ```bash
   glab auth login --hostname gitlab.example.com --token <YOUR_TOKEN>
   ```

This stores the PAT in git's credential helper, which markdown-proxy uses automatically.

### Organization-specific Credentials

If you need different credentials for specific organizations, use path-based credential configuration:

```gitconfig
[credential "https://github.com/my-org"]
    helper = !gh auth git-credential
    useHttpPath = true
```

### Verifying Credentials

To verify that git can provide credentials for a host:

```bash
echo -e "protocol=https\nhost=github.com\n" | git credential fill
```

This should output `username` and `password` fields. If nothing is returned, credentials are not configured for that host.

## Installation

### Download

Download the latest binary from [GitHub Releases](https://github.com/patakuti/markdown-proxy/releases).

### go install

```bash
go install github.com/patakuti/markdown-proxy/cmd/markdown-proxy@latest
```

## Build

```bash
make build
```

### Cross-compile

```bash
# Linux
make linux

# Windows
make windows
```

### Manual build

```bash
go build -o markdown-proxy ./cmd/markdown-proxy
```

## Known Limitations

- **Read-only viewer**: No editing capabilities; this is a rendering-only tool.
- **Limited file type support**: Only `.md`, `.markdown`, and `.txt` files are rendered as HTML. Other file types are served as-is.
- **PlantUML disabled by default**: Diagram content is sent to an external server, so it requires explicit opt-in via `--plantuml-server` or `--configure`. When disabled, a hint message is shown in place of PlantUML blocks.
- **GitHub/GitLab branch detection**: When accessing a repository root URL, only `main` and `master` branches are tried for README.md auto-detection.
- **No native PDF export**: Use the toolbar's Print link to export via the browser's print-to-PDF feature. Page breaks are automatically avoided inside tables, code blocks, math expressions, images, blockquotes, and list items; headings are kept together with the following content.
- **Hidden files excluded**: Files and directories starting with `.` are not shown in directory listings.

## Contributing

Bug reports and feature requests are welcome via [GitHub Issues](https://github.com/patakuti/markdown-proxy/issues). Pull requests are also appreciated — please open an issue first to discuss the change.

```bash
# Build and test locally
make build
```

## Project Structure

```
cmd/markdown-proxy/    - Entry point
internal/
  config/              - Command-line flag parsing, mode detection
  server/              - HTTP server, routing, middleware (auth, access log)
  handler/             - Request handlers (top, local, remote, SSE, login)
  network/             - HTTP client with SSRF protection
  markdown/            - Markdown→HTML conversion, link rewriting, code block processing
  credential/          - git credential helper integration
  github/              - GitHub/GitLab URL resolution
  template/            - HTML templates and structural CSS
  themes/              - Theme management (built-in CSS generation, file I/O)
```
