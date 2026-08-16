#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
today="$(date +%F)"
commit_limit="${DEVLOG_COMMITS:-20}"
archive_dir="$root/docs/devlogs"
month="${today%-??}"
archive_dir="${DEVLOG_ARCHIVE_DIR:-$archive_dir}"
archive_file="$archive_dir/$month/$today.md"
index_file="${DEVLOG_INDEX_FILE:-$root/DEVLOG.md}"
temporary_dir="$(mktemp -d "${TMPDIR:-/tmp}/notegic-devlog.XXXXXX")"
trap 'rm -rf "$temporary_dir"' EXIT INT TERM

if [[ ! "$commit_limit" =~ ^[1-9][0-9]*$ ]]; then
	echo "DEVLOG_COMMITS must be a positive integer" >&2
	exit 1
fi

mkdir -p "$archive_dir/$month"

commits="$(git -C "$root" log -n "$commit_limit" --date=short --pretty=format:'- `%h` %ad — %s' HEAD)"
if [[ -z "$commits" ]]; then
	commits='- No commits found.'
fi

areas="$(git -C "$root" log -n "$commit_limit" --name-only --format='' HEAD \
	| awk 'NF {split($0, parts, "/"); print parts[1]}' \
	| sort | uniq -c | sort -nr \
	| awk '{print "- `" $2 "` — " $1 " changed-file references"}')"
if [[ -z "$areas" ]]; then
	areas='- No changed-file information found.'
fi

temporary_file="$temporary_dir/archive.md"
{
	printf '# Development log — %s\n\n' "$today"
	printf 'This snapshot is generated from the latest %s Git commits. It is a deterministic repository summary; detailed intent belongs in commit messages and design documents.\n\n' "$commit_limit"
	printf '## Recent commits\n\n%s\n\n' "$commits"
	printf '## Changed areas\n\n%s\n\n' "$areas"
	cat <<'EOF'
## Regeneration

Run `make devlog` from the repository root. The pre-commit hook validates the staged result when enabled with `make install-hooks`.
EOF
} > "$temporary_file"
mv "$temporary_file" "$archive_file"

root_temporary_file="$temporary_dir/index.md"
recent_snapshots="$(find "$archive_dir" -maxdepth 2 -type f -name '20??-??-??.md' -print \
	| awk -F/ '{key=$NF; sub(/\\.md$/, "", key); print key "\t" $0}' \
	| sort -r \
	| head -n 5 \
	| cut -f2-)"
{
	printf '# Development Log\n\n'
	printf 'This file is an automatically maintained index of recent snapshots.\n\n'
	printf '## Recent snapshots\n\n'
	while IFS= read -r file; do
		[[ -z "$file" ]] && continue
		relative="${file#$root/}"
		name="${relative#docs/devlogs/}"
		name="${name%.md}"
		printf -- '- [%s](%s)\n' "$name" "$relative"
	done <<< "$recent_snapshots"
} > "$root_temporary_file"
mv "$root_temporary_file" "$index_file"

printf 'Generated %s and refreshed DEVLOG.md\n' "${archive_file#$root/}"
