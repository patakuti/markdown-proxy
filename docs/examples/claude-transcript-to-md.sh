#!/bin/bash
# Convert a Claude Code transcript (JSONL) to Markdown so markdown-proxy can
# show a conversation with live reload.
#
# One-shot conversion (prints Markdown to stdout):
#   claude-transcript-to-md.sh TRANSCRIPT.jsonl > conversation.md
#
# As a Claude Code Stop hook (reads the hook JSON from stdin). Choose where
# the file goes:
#
#   --hook FILE                  one fixed file (the latest conversation)
#   --hook --dir DIR --per workspace
#                                DIR/<workspace>.md           one file per workspace
#   --hook --dir DIR --per session
#                                DIR/<workspace>/<session>.md one file per session
#
# <workspace> is the name of the directory the session was started in.
# The file is rewritten from the whole transcript each time, so a missed hook
# never leaves a gap. Only user prompts and assistant text are kept; thinking,
# tool calls and tool results are omitted. Requires jq.
set -eu

convert() {
  jq -r '
    def text_of: if type == "string" then .
                 else map(select(.type == "text") | .text) | join("\n\n") end;
    select(.type == "user" or .type == "assistant")
    | select(.isMeta | not)
    | (.message.content | text_of) as $t
    | select($t | length > 0)
    | "\n---\n\n## " + (if .type == "user" then "User" else "Claude" end) + "\n\n" + $t
  ' "$1"
}

usage() {
  sed -n '2,/^set -eu/p' "$0" | sed '$d; s/^# \{0,1\}//' >&2
  exit 2
}

if [ "${1:-}" != "--hook" ]; then
  [ $# -eq 1 ] || usage
  convert "$1"
  exit 0
fi
shift

out="" dir="" per=""
while [ $# -gt 0 ]; do
  case "$1" in
    --dir) dir="$2"; shift 2 ;;
    --per) per="$2"; shift 2 ;;
    -*) usage ;;
    *) out="$1"; shift ;;
  esac
done

hook=$(cat)
transcript=$(echo "$hook" | jq -r '.transcript_path')
session=$(echo "$hook" | jq -r '.session_id')
cwd=$(echo "$hook" | jq -r '.cwd')

if [ -z "$out" ]; then
  [ -n "$dir" ] && { [ "$per" = workspace ] || [ "$per" = session ]; } || usage

  # $cwd is the directory at the time the hook fires; it can be a
  # subdirectory if the session ran `cd`. The transcript lives in
  # ~/.claude/projects/<start directory with "/" replaced by "-">/, so walk up
  # from $cwd until the encoded name matches to recover the start directory.
  encoded=$(basename "$(dirname "$transcript")")
  candidate="$cwd"
  workspace=""
  while [ -n "$candidate" ] && [ "$candidate" != "/" ]; do
    if [ "$(echo "$candidate" | tr '/' '-')" = "$encoded" ]; then
      workspace=$(basename "$candidate")
      break
    fi
    candidate=$(dirname "$candidate")
  done
  [ -n "$workspace" ] || workspace=$(basename "$cwd")

  if [ "$per" = workspace ]; then
    out="$dir/$workspace.md"
  else
    out="$dir/$workspace/$session.md"
  fi
fi

# The Stop hook can fire before the final response has been written to the
# transcript, so give Claude Code a moment to flush it. Configure the hook with
# "async": true so this wait does not block anything.
sleep 1

# Only --dir creates directories; a mistyped fixed FILE path should fail loudly.
[ -z "$dir" ] || mkdir -p "$(dirname "$out")"
# Write to a temporary file first so the viewer never sees a half-written file.
convert "$transcript" > "$out.tmp"
mv "$out.tmp" "$out"
