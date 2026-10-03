[← Back to README](../README.md)

# Watching a Claude Code conversation

Claude Code stores each conversation as a JSONL transcript. A Claude Code `Stop` hook (it runs each time Claude finishes responding) can convert it to Markdown, and markdown-proxy re-renders the file as it changes.

1. Copy [`examples/claude-transcript-to-md.sh`](examples/claude-transcript-to-md.sh) somewhere on your machine (it needs `jq`). It keeps your prompts and Claude's text replies; thinking, tool calls and tool results are left out.
2. Add the hook to `~/.claude/settings.json`. Choose where the output goes:

   | Output | Option | File |
   |---|---|---|
   | One fixed file | `--hook FILE` | `FILE` (always the latest conversation) |
   | One file per workspace | `--hook --dir DIR --per workspace` | `DIR/<workspace>.md` |
   | One file per session | `--hook --dir DIR --per session` | `DIR/<workspace>/<session id>.md` |

   `<workspace>` is the name of the directory the session was started in. Example with a fixed file:

   ```json
   {
     "hooks": {
       "Stop": [
         {
           "hooks": [
             {
               "type": "command",
               "command": "/path/to/claude-transcript-to-md.sh --hook ~/claude-conversation.md",
               "async": true
             }
           ]
         }
       ]
     }
   }
   ```

   `"async": true` matters: the `Stop` hook can fire before Claude Code has written the final response to the transcript, so the script waits one second before reading it. Asynchronous execution keeps that wait from blocking anything.

3. Open the output file once; it refreshes after each response:

   ```bash
   markdown-proxy ~/claude-conversation.md
   ```

The file is rewritten from the whole transcript every time, so it never has gaps. Directories are created for `--dir`, but not for a fixed `FILE`, so a mistyped path fails with an error. To convert an existing transcript once, run `claude-transcript-to-md.sh ~/.claude/projects/<project>/<session>.jsonl > conversation.md`.
