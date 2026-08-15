# Version control conventions

This document defines commit messages, generated development logs, and the
local Git hooks used by the repository.

## Commit message format

Commit headers use this form:

```text
<type>(<optional-scope>)!: <description>
```

Examples:

```text
feat(api-gateway): add api key middleware
fix(core): reject expired api keys
hotfix(client-gateway): restore broken auth cookie handling
refactor: split client gateway runtime
docs: update version control rules
```

The scope is optional and uses lowercase kebab-case. The `!` marker is
optional and indicates a breaking change.

### Types

| Type | Use for |
| --- | --- |
| `feat` | A new capability or user-visible behavior |
| `fix` | A bug or a long-standing design defect |
| `hotfix` | An immediate correction to the previous commit |
| `refactor` | Internal restructuring without intended behavior change |
| `docs` | Documentation-only changes |
| `test` | Test-only changes |
| `chore` | Maintenance that does not fit another type |
| `ci` | CI or automation changes |
| `build` | Build, dependency, or packaging changes |
| `perf` | A performance improvement |
| `revert` | Reverting an earlier change |

Use `hotfix` only when the immediately preceding commit is broken and needs an
urgent correction. Use `fix` for an older bug, a persistent issue, or a design
defect discovered after the original change.

### Description rules

- Write the description in English.
- Use lowercase English words by default.
- Prefer plain-language descriptions over internal mechanism names. For
  example, use `transaction synchronizer` instead of
  `TransactionSynchronizer`.
- If an exact proper noun or mechanism name is unavoidable, wrap it in
  backticks, for example:

  ```text
  refactor(core): replace `TransactionSynchronizer` with transaction synchronizer
  ```

- Keep the subject at or below 72 characters when practical; 120 characters is
  the hard limit enforced by the validator.
- Put implementation detail in the body instead of extending the subject.
- The body is optional and must be separated from the subject by one blank
  line.

### Breaking changes

Use `!` in the header and describe the migration in a final `BREAKING CHANGES:`
block. Each breaking change is a separate bullet:

```text
feat(api)!: change api key authentication

BREAKING CHANGES:
- clients must send the x-api-key header
- cookie authentication is no longer accepted by this route
```

The block is required when `!` is used and may contain multiple bullet points.

Automatically generated merge messages, `Revert "..."` messages, and
`fixup!`/`squash!` autosquash messages are accepted as Git workflow exceptions.

## Development logs

Run the generator before committing a repository change:

```bash
make devlog
git add .
git commit -m "your message"
git push origin main
```

Daily logs are stored by month:

```text
docs/devlogs/YYYY-MM/YYYY-MM-DD.md
```

Each daily file contains:

- `Recent commits`: commits grouped on that date.
- `Changed areas`: top-level repository areas touched by those commits.
- `Regeneration`: the command and hook behavior.

`make devlog` writes today's snapshot and refreshes the root
[`DEVLOG.md`](../../DEVLOG.md) index. It is deterministic and uses local Git
history; it does not generate an AI summary or call GitHub.

## Git hooks

Enable the repository hooks once per local clone:

```bash
make install-hooks
```

The `pre-commit` hook verifies that today's devlog exists in the Git index and
matches a fresh generation. It generates the comparison file in a temporary
directory, compares it with the staged version, and removes the temporary
directory on every exit path. It never overwrites the working tree.

The `commit-msg` hook runs `scripts/validate_commit_message.sh` against Git's
commit message file. A non-zero result blocks the commit and prints the rule
that needs correction.
