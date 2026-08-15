#!/usr/bin/env bash
set -euo pipefail

message_file="${1:-}"
if [[ -z "$message_file" || ! -f "$message_file" ]]; then
	echo "usage: validate_commit_message.sh <commit-message-file>" >&2
	exit 2
fi

header="$(sed -n '1p' "$message_file" | tr -d '\r')"

case "$header" in
	Merge\ *) exit 0 ;;
	Revert\ \"*\") exit 0 ;;
	fixup\!\ *|squash\!\ *) exit 0 ;;
esac

if [[ ! "$header" =~ ^(feat|fix|hotfix|refactor|docs|test|chore|ci|build|perf|revert)(\([a-z0-9][a-z0-9-]*\))?(!)?:[[:space:]](.+)$ ]]; then
	echo "invalid commit header; expected <type>(<optional-scope>)!: <description>" >&2
	exit 1
fi

if (( ${#header} > 120 )); then
	echo "commit subject must be 120 characters or fewer" >&2
	exit 1
fi

description="${header#*: }"
plain_description="$(printf '%s' "$description" | sed -E 's/`[^`]+`//g')"

if [[ ! "$plain_description" =~ [a-z] ]]; then
	echo "commit description must contain lowercase English words" >&2
	exit 1
fi

if [[ "$plain_description" =~ [^[:space:][:lower:][:digit:][:punct:]] ]]; then
	echo "commit description must use lowercase English; wrap unavoidable proper nouns in backticks" >&2
	exit 1
fi

if [[ "$header" =~ !: ]]; then
	if ! grep -q '^BREAKING CHANGES:$' "$message_file"; then
		echo "breaking commits must include a BREAKING CHANGES: block" >&2
		exit 1
	fi
fi

breaking_line="$(grep -n '^BREAKING CHANGES:$' "$message_file" | head -n 1 | cut -d: -f1 || true)"
if [[ -n "$breaking_line" ]]; then
	if ! tail -n +$((breaking_line + 1)) "$message_file" \
		| awk '/^#/ {next} NF {found=1; if ($0 !~ /^- /) invalid=1} END {exit !(found && !invalid)}'; then
		echo "BREAKING CHANGES: must be followed by one or more bullet points" >&2
		exit 1
	fi
fi
