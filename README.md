# markdown-proxy

[![CI](https://github.com/patakuti/markdown-proxy/actions/workflows/ci.yml/badge.svg)](https://github.com/patakuti/markdown-proxy/actions/workflows/ci.yml)

A Markdown viewer that runs as a local HTTP server and renders files in your browser. It is built mainly for viewing local files, and can open GitHub/GitLab files the same way.

![Demo](docs/demo.gif)

## Motivation

- **Live reload**: Files that are rewritten while you view them — a document you are writing, or a Claude Code conversation log written to a file automatically — are re-rendered in the browser as they change. Scrolling back through earlier conversation and copying a part of it are easy in the rendered view.
- **PlantUML diagrams**: GitHub and GitLab often do not render PlantUML diagrams embedded in Markdown and show them as raw code blocks. markdown-proxy renders them (see `--plantuml-server`).

## Comparison

How markdown-proxy compares to other Markdown viewing tools:

| Feature | markdown-proxy (this tool) | [Markdown Preview Enhanced][mpe] | [grip][grip] | [Madness][madness] | [mdserve][mdserve] ³ |
|---------|:-:|:-:|:-:|:-:|:-:|
| **Overview** | | | | | |
| Type | HTTP server | VS Code extension | HTTP server | HTTP server | HTTP server |
| Runtime dependency | None (single binary) | VS Code | Python | Ruby | None (single binary) |
| Works offline | △ ⁴ | ✅ | ❌ | ✅ | ✅ |
| **Viewing** | | | | | |
| Local file viewing | ✅ | ✅ | ✅ | ✅ | ✅ |
| Live reload | ✅ | ✅ | ✅ | ❌ | ✅ |
| Directory listing | ✅ | ❌ | ❌ | ✅ | ✅ |
| CLI open | ✅ (file/URL) | — | ✅ (file, -b) | ✅ (--open) | ✅ (--open) |
| Line anchor links (`foo.md:12`) | ✅ | Unknown ² | Unknown ² | Unknown ² | Unknown ² |
| Full-text search | ❌ | ❌ | ❌ | ✅ | ❌ |
| **Rendering** | | | | | |
| Code highlighting | ✅ | ✅ | ✅ | ✅ | ✅ |
| Math rendering | ✅ (KaTeX) | ✅ (KaTeX/MathJax) | ❌ | ❌ | ❌ |
| Mermaid diagrams | ✅ | ✅ | ❌ | ✅ | ✅ |
| PlantUML diagrams | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Customization and output** | | | | | |
| CSS themes | 3 built-in + user-defined | 15+ built-in | GitHub only | Customizable | 5 built-in |
| Export (PDF, HTML) | △ ¹ | ✅ (PDF, HTML, Word) | ✅ (HTML) | ❌ | ❌ |
| **Remote files** | | | | | |
| Remote URL fetching | ✅ (GitHub/GitLab) | ❌ | ❌ | ❌ | ❌ |
| Authentication | Token-based | — | — | HTTP Basic | ❌ |

¹ Browser print-to-PDF via toolbar Print link
² Not documented in the project's README or documentation
³ Archived by its author in September 2026 (no longer maintained)
⁴ Markdown, code highlighting and themes work offline; Mermaid and math (KaTeX) load scripts from a CDN and need network access

Based on each project's README, documentation, and package metadata as of October 2026.

[mpe]: https://marketplace.visualstudio.com/items?itemName=shd101wyy.markdown-preview-enhanced
[grip]: https://github.com/joeyespo/grip
[madness]: https://github.com/DannyBen/madness
[mdserve]: https://github.com/jfernandez/mdserve

## Features

### Viewing

- Live reload for local files (auto-refreshes the browser on file changes)
- Directory listing for local files
- Open a file or URL from the command line (starts the server automatically if needed)
- Top page with smart input (auto-detects file path or URL) and recently opened file history (localStorage)
- CSS themes with dropdown switching — 3 built-in themes (GitHub, Simple, Dark) plus user-defined themes
- Toolbar actions
  - Print: opens the browser's print dialog with the page title set to the bare file name (e.g. `foo` instead of `foo.md - markdown-proxy`), so the print header and the default PDF file name are clean; the toolbar is hidden in print output
  - Source: link to original URL on remote server (remote pages only)

### Rendering

- GFM (GitHub Flavored Markdown) with syntax highlighting and a copy button on code blocks
- Math rendering (`$...$` for inline, `$$...$$` for display) via KaTeX
- Code block rendering
  - SVG: inline SVG rendering from ```` ```svg ```` code blocks
  - Mermaid: client-side rendering via mermaid.js from ```` ```mermaid ```` code blocks
  - PlantUML: server-side rendering from ```` ```plantuml ```` code blocks (requires `--plantuml-server`)
    - Disabled by default because diagram content is sent to the specified server
    - When disabled, a hint is shown in place of each PlantUML block
    - To use the public server: `--plantuml-server https://www.plantuml.com/plantuml`
    - Use `--configure` to save the setting permanently
- Text file rendering: `.txt` files are displayed in HTML with line anchors, themes, and live reload

### Navigation

- Table of contents sidebar: toggle `TOC` in the toolbar to open a right-side panel with auto-extracted headings; visibility persists per browser (localStorage)
- Line anchor links: `[text](foo.md:12)` or `<a href="foo.md:12">` links navigate to specific source lines with highlighting (Markdown and text files)
- Link rewriting for seamless proxy navigation (including `file:///` protocol conversion)

### Remote files (GitHub/GitLab)

- Render remote Markdown files; repository root URLs show the README
- Blob URL auto-conversion to raw URL (supports self-hosted GitLab with custom domains)
- Authentication via git credential helper (supports path-based credential matching)
- Redirect-based auth detection for self-hosted GitLab instances

### Server operation

- Two operation modes: local mode and remote mode
- Token-based authentication for remote access
- Access logging with automatic log rotation
- Configuration file (`~/.config/markdown-proxy/config.json`) with interactive setup via `--configure`
- Single binary, no runtime dependencies (Mermaid and math rendering load scripts from a CDN, so they need network access)

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

When a file or URL is given as an argument, markdown-proxy checks if the server is already running; if not, it starts one in the background. Then it opens the file in your default browser. Without arguments, it starts the server in the foreground.

## Use Cases

- **Reading a Claude Code conversation as it grows**: Combine a Claude Code `Stop` hook with live reload to keep the whole conversation open in the browser as rendered Markdown. It is updated after every response whether or not the agent decides to show anything, so you can scroll back through earlier turns and copy parts of them. See [Watching a Claude Code conversation](docs/claude-code-conversation.md).
- **Reading any file that keeps changing**: A document you are writing, a generated report or any other file rewritten by a tool is re-rendered as it changes.
- **Viewing PlantUML diagrams**: Render PlantUML embedded in Markdown that GitHub/GitLab don't display natively.
- **Navigating RAG search reports**: Reports generated by tools like [Local Knowledge RAG MCP Server](https://github.com/patakuti/local-knowledge-rag-mcp) contain links to source documents. In markdown-proxy, open the report in one tab and click through references in new tabs.
- **Browsing private repositories**: Access Markdown files from private GitHub/GitLab repos using your existing git credentials.
- **Sharing a Markdown viewer with your team**: Run in remote mode with token authentication to let team members view documentation through a browser.

## Documentation

- [Usage reference](docs/usage.md): command-line options, URL scheme and line anchor links, configuration file, CSS themes, local/remote modes, access logging, live reload, math rendering
- [Watching a Claude Code conversation](docs/claude-code-conversation.md): a `Stop` hook recipe for reading a conversation as it grows
- [Private repositories](docs/private-repositories.md): reading private GitHub/GitLab repositories with your git credentials

## Security

- **Local mode**: The server binds to `127.0.0.1` only, accepting local connections only
- **Remote mode**: Authentication is enforced via token. Local file access and private network fetching are automatically disabled
- **SSRF protection**: In remote mode, requests to private/internal IP addresses (e.g., `10.x.x.x`, `192.168.x.x`, `127.x.x.x`) are blocked. In local mode, private network access is allowed
- **DNS rebinding prevention**: Resolved IP addresses are used directly for connections, preventing DNS rebinding attacks
- **Constant-time token comparison**: Authentication uses `crypto/subtle.ConstantTimeCompare` to prevent timing attacks
- **HTML sanitizing**: Raw HTML in rendered Markdown is sanitized with an allowlist. Tables, images, `<details>` and inline SVG keep working; `<script>`, `<iframe>`, `on*` handlers and `javascript:` URLs are removed. A document you open cannot run script that reads `/local/...` files
- **Content Security Policy**: Every page is served with a per-request nonce CSP (`script-src 'nonce-...'`, `connect-src 'self'`), so injected script does not run and cannot send data to other hosts
- **Host / Origin checks** (local mode): Requests whose `Host` is not `localhost`, `127.0.0.1` or `[::1]` are rejected (DNS rebinding). Cross-origin `fetch` / `EventSource` requests are rejected in both modes, and no CORS headers are sent, so other websites cannot read your files through the browser
- **Files served as-is**: HTML, SVG and XML files are served with `Content-Security-Policy: sandbox allow-scripts` and `X-Content-Type-Options: nosniff`. Their scripts run in an opaque origin and cannot read other proxy URLs
- **Pinned CDN scripts**: Mermaid and KaTeX are loaded from jsDelivr at exact versions with Subresource Integrity hashes
- **Credential forwarding**: Authorization headers are not forwarded when a remote server redirects to another host

Note that local mode serves any file the user can read under `/local/<absolute path>` (by design, so logs such as `~/.claude/` can be viewed). The protections above prevent other web content from reaching it through your browser.

## Installation

### Download

Download the latest binary from [GitHub Releases](https://github.com/patakuti/markdown-proxy/releases). Prebuilt binaries are provided for Linux and Windows (amd64). On macOS, install with `go install` or build from source (below); the macOS build compiles but has not been tested on a Mac by the maintainer.

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
- **Limited file type support**: Only `.md`, `.markdown`, and `.txt` files are rendered as HTML. Other file types are served as-is (HTML/SVG/XML in a script sandbox, see [Security](#security)).
- **PlantUML disabled by default**: Diagram content is sent to an external server, so it requires explicit opt-in via `--plantuml-server` or `--configure`. When disabled, a hint message is shown in place of PlantUML blocks.
- **GitHub/GitLab branch detection**: When accessing a repository root URL, only `main` and `master` branches are tried for README.md auto-detection.
- **No native PDF export**: Use the toolbar's Print link to export via the browser's print-to-PDF feature. Page breaks are automatically avoided inside tables, code blocks, math expressions, images, blockquotes, and list items; headings are kept together with the following content. Background colors are forced to print via `print-color-adjust: exact`, but the browser's print dialog must also have "Background graphics" enabled (in Chrome: More settings > Background graphics), otherwise the browser suppresses them regardless of this setting.
- **Hidden files excluded**: Files and directories starting with `.` are not shown in directory listings.
- **List marker spacing is normalized**: Under strict CommonMark/GFM, a list marker (e.g. `1.`) followed by 5 or more spaces turns the rest of the line into an indented code block instead of plain text — an easy mistake to make when padding markers for visual alignment. This proxy always normalizes the spacing after a list marker to a single space, so `1.     a` renders the same as `1. a`. This is an intentional deviation from GitHub's rendering. The normalization is applied per line, so deeply nested list items (indented past 3 columns) may not be covered.

## Contributing

Bug reports and feature requests are welcome via [GitHub Issues](https://github.com/patakuti/markdown-proxy/issues). Pull requests are also appreciated — please open an issue first to discuss the change.

```bash
# Build, test and check formatting locally
make build
make test
make fmt-check
```

CI runs `go vet` and `go test` on Linux (with the race detector) and Windows, plus a formatting check on Linux, for every pull request.

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

## History

How the project grew, and where its value turned out to be.

1. **Motivation (Feb 2026)**: The starting point was two pain points — IDE Markdown previews show one document at a time, and GitHub/GitLab do not render PlantUML. (Later research found that Markdown Preview Enhanced can also open multiple previews through a setting, so viewing several documents in browser tabs is a convenience rather than a distinguishing feature.) The first version served local and remote Markdown through a proxy with SVG, Mermaid, and PlantUML rendering. Live reload for local files was added on the first day.
2. **Remote access (Feb 2026)**: Authentication through the git credential helper, self-hosted GitLab support, a remote mode with token authentication, and access logging.
3. **Navigating documents (Mar–Apr 2026)**: Line anchor links (`foo.md:12`), which let reports such as RAG search results link to exact source lines, then text file rendering, a configuration file, opening a file from the command line, and the table of contents sidebar.
4. **Polish (May–Aug 2026)**: Copy buttons, user-defined themes, print-quality fixes, and rendering fixes found through daily use.
5. **A new value: a live local viewer**: In daily use, the most frequent use turned out to be plain local viewing with live reload. For example, writing Claude Code conversation logs to a file automatically and watching them re-render makes it easy to look back at earlier conversation and copy a part of it. The README was reorganized around this use.
6. **Context: agent-side rendering (Sep 2026)**: mdserve, a similar tool built for viewing Markdown written by AI coding agents (with a Claude Code plugin), was archived by its author, who cited agent CLIs shipping their own rendering, such as Artifacts in Claude Code. markdown-proxy has a different focus: it views any Markdown file on your own machine, including files that are written continuously without an agent deciding to show them.

## About this project

This tool was designed and implemented entirely by Claude. The human provided the idea. However, this isn't a one-shot output; the human shaped it through hands-on testing and iterative, detail-oriented feedback.
