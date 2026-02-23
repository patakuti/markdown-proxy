# markdown-proxy

An HTTP proxy server for viewing Markdown files in a browser. Supports both local and remote access modes.

## Motivation

Existing Markdown viewing environments have several pain points:

- **Limited local preview**: Tools like VS Code's Markdown Preview Enhanced are widely used for local Markdown viewing, but they can only display one preview pane at a time, making it difficult to reference multiple documents simultaneously.
- **Incomplete diagram rendering on remote servers**: Platforms such as GitHub and GitLab often do not render PlantUML diagrams embedded in Markdown files, leaving them displayed as raw code blocks.

markdown-proxy solves these problems by rendering Markdown files — including PlantUML diagrams — as HTML in a standard browser, where you can freely open multiple tabs and windows.

## Features

- Render local and remote Markdown files as HTML
- Support for GFM (GitHub Flavored Markdown) with syntax highlighting
- Math rendering (`$...$` for inline, `$$...$$` for display) via KaTeX
- Code block rendering
  - SVG: inline SVG rendering from ```` ```svg ```` code blocks
  - Mermaid: client-side rendering via mermaid.js from ```` ```mermaid ```` code blocks
  - PlantUML: server-side rendering from ```` ```plantuml ```` code blocks (requires `--plantuml-server`)
    - Disabled by default because diagram content is sent to the specified server
    - To use the public server: `--plantuml-server https://www.plantuml.com/plantuml`
- GitHub/GitLab integration
  - Blob URL auto-conversion to raw URL (supports self-hosted GitLab with custom domains)
  - Authentication via git credential helper (supports path-based credential matching)
  - Redirect-based auth detection for self-hosted GitLab instances
- Multiple CSS themes (GitHub, Simple, Dark) with switching UI
- Live reload for local files (auto-refreshes browser on file changes)
- Directory listing for local files
- Link rewriting for seamless proxy navigation
- Top page with smart input (auto-detects file path or URL)
- Recently opened file history (localStorage)
- Two operation modes: local mode and remote mode
- Token-based authentication for remote access
- Access logging with automatic log rotation
- Single binary, no runtime dependencies

## URL Scheme

| Type | URL Format |
|------|-----------|
| Top page | `http://localhost:9080/` |
| Local file | `http://localhost:9080/local/path/to/file.md` |
| Local directory | `http://localhost:9080/local/path/to/dir/` |
| Remote (HTTP) | `http://localhost:9080/http/server/path/to/file.md` |
| Remote (HTTPS) | `http://localhost:9080/https/server/path/to/file.md` |
| GitHub repo | `http://localhost:9080/https/github.com/user/repo/blob/main/README.md` |

## Usage

```bash
markdown-proxy [options]
```

### Options

> **Note:** Options use a single dash (e.g., `-port`) following Go's `flag` package convention. Double dashes (`--port`) also work.

| Flag | Description | Default |
|------|-------------|---------|
| `-port`, `-p` | Listen port | `9080` |
| `-listen` | Bind address (`127.0.0.1` for local, `0.0.0.0` for remote) | `127.0.0.1` |
| `-theme` | Default CSS theme (`github`, `simple`, `dark`) | `github` |
| `-plantuml-server` | PlantUML server URL | (disabled) |
| `-auth-token` | Authentication token (required in remote mode) | |
| `-auth-cookie-max-age` | Authentication cookie max age in days | `30` |
| `-access-log` | Access log file path | |
| `-access-log-max-size` | Max log file size in MB before rotation | `100` |
| `-access-log-max-backups` | Max number of old log files to retain | `3` |
| `-access-log-max-age` | Max days to retain old log files | `28` |
| `-verbose`, `-v` | Enable debug logging to stderr | `false` |

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
  template/            - HTML templates and CSS themes
```
