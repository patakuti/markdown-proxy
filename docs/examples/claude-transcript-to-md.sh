#!/bin/sh
# Convert a Claude Code transcript (JSONL) to Markdown.
#
# Usage:
#   claude-transcript-to-md.sh TRANSCRIPT.jsonl > conversation.md
#   claude-transcript-to-md.sh --hook OUTPUT.md   # as a Claude Code hook: reads the hook JSON from stdin
#
# Requires jq. Only user prompts and assistant text are kept; thinking,
# tool calls and tool results are omitted.
set -eu

convert() {
  jq -r '
    select(.type == "user" or .type == "assistant")
    | select(.isMeta | not)
    | if .type == "user" then
        (.message.content | select(type == "string") | "## User\n\n" + . + "\n")
      else
        ([.message.content[] | select(.type == "text") | .text] | select(length > 0)
         | "## Claude\n\n" + join("\n\n") + "\n")
      end
  ' "$1"
}

if [ "${1:-}" = "--hook" ]; then
  out="$2"
  transcript=$(jq -r '.transcript_path')
  # Write to a temporary file first so the viewer never sees a half-written file.
  convert "$transcript" > "$out.tmp"
  mv "$out.tmp" "$out"
else
  convert "$1"
fi
