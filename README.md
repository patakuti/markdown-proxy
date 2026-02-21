# markdown-proxy

A local HTTP proxy server for viewing Markdown files in a browser.

## Features

- Render local and remote Markdown files as HTML
- Support for GFM (GitHub Flavored Markdown) with syntax highlighting
- Code block rendering
  - SVG: inline SVG rendering from ```` ```svg ```` code blocks
  - Mermaid: client-side rendering via mermaid.js from ```` ```mermaid ```` code blocks
  - PlantUML: server-side rendering from ```` ```plantuml ```` code blocks
- GitHub/GitLab integration
  - Blob URL auto-conversion to raw URL
  - Authentication via git credential helper (on 401/403 only)
- Multiple CSS themes (GitHub, Simple, Dark) with switching UI
- Live reload for local files (auto-refreshes browser on file changes)
- Directory listing for local files
- Link rewriting for seamless proxy navigation
- Top page with smart input (auto-detects file path or URL)
- Recently opened file history (localStorage)
- Localhost-only binding (127.0.0.1) for security
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

| Flag | Description | Default |
|------|-------------|---------|
| `--port`, `-p` | Listen port | `9080` |
| `--theme` | Default CSS theme (`github`, `simple`, `dark`) | `github` |
| `--plantuml-server` | PlantUML server URL | `https://www.plantuml.com/plantuml` |
| `--allow-private-network` | Allow fetching from private/internal IP addresses | `false` |
| `--verbose`, `-v` | Enable access logging | `false` |

## Live Reload

When viewing local Markdown files or directories (`/local/...`), the browser automatically reloads when the file or directory contents change. This uses Server-Sent Events (SSE) with filesystem notifications (fsnotify).

- **Local files only**: Remote files (`/http/...`, `/https/...`) are not affected
- **No configuration needed**: Works automatically for all local file views
- **Debounced**: Multiple rapid changes are coalesced into a single reload (100ms debounce)

## Security

- **Localhost-only**: The server binds to `127.0.0.1` only. It is not intended for network-facing deployment.
- **SSRF protection**: By default, requests to private/internal IP addresses (e.g., `10.x.x.x`, `192.168.x.x`, `127.x.x.x`) are blocked when fetching remote files. Use `--allow-private-network` to disable this restriction.
- **DNS rebinding prevention**: Resolved IP addresses are used directly for connections, preventing DNS rebinding attacks.

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
  config/              - Command-line flag parsing
  server/              - HTTP server and routing
  handler/             - Request handlers (top, local, remote, SSE)
  network/             - HTTP client with SSRF protection
  markdown/            - Markdown→HTML conversion, link rewriting, code block processing
  credential/          - git credential helper integration
  github/              - GitHub/GitLab URL resolution
  template/            - HTML templates and CSS themes
```
