[← Back to README](../README.md)

# Private Repository Access

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
