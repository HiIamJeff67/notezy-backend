# Versioning and Release Tags

This document defines the Notezy Backend versioning, Git tag, GitHub Release,
and deployment version management conventions.

## Tag and Release Naming

The Git tag name and the displayed Release name are separate fields:

- Git tag name: `v0.1.0-beta.1`
- Annotated tag message: `Release v0.1.0-beta.1`
- GitHub Release title: `Release v0.1.0-beta.1`

The current GitHub Actions workflow listens for tags beginning with `v` (`v*`).
Therefore, do not use `Release v0.1.0-beta.1` as the Git tag name because it
does not match the current release trigger.

`Release v...` is the human-readable title and annotation. It is not part of
the Git tag name.

## Tag Format

Stable releases use Semantic Versioning:

```text
v<major>.<minor>.<patch>
```

Examples:

```text
v1.0.0
v1.1.0
v1.1.1
```

Pre-release versions use:

```text
v<major>.<minor>.<patch>-<stage>.<number>
```

The supported pre-release stages are:

```text
alpha
beta
rc
```

Examples:

```text
v0.1.0-alpha.1
v0.1.0-beta.1
v0.1.0-beta.2
v1.0.0-rc.1
```

The recommended first Beta release is:

```text
v0.1.0-beta.1
```

Version increments follow these rules:

| Change type | Example |
| --- | --- |
| Backward-compatible bug fix | `v0.1.1` |
| New backward-compatible feature | `v0.2.0` |
| Breaking change | `v1.0.0` |
| Beta revision | `v0.1.0-beta.2` |
| Release Candidate revision | `v1.0.0-rc.2` |

## Branch, Commit, and Tag Responsibilities

Each naming layer has a separate purpose:

```text
Branch: feature/notification-inbox
Commit: feat(notification): add notification inbox endpoint
Tag:    v0.1.0-beta.1
```

- A branch describes a group of work in progress.
- A commit describes one change and follows the project's Conventional Commit format.
- A tag identifies a confirmed, deployable version.

The following must not be used as release tags:

```text
feature-notification
refactor-microservices
chore-ci
```

Those meanings belong in branch names or commit messages.

## Release Workflow

Create a release tag only on a commit that is already merged into `main` and
has passed CI:

```bash
git switch main
git pull --ff-only origin main

git tag -a v0.1.0-beta.1 -m "Release v0.1.0-beta.1"
git push origin v0.1.0-beta.1
```

After the tag is pushed:

1. GitHub Actions runs formatting, unit, race, generated-contract, and container checks.
2. When the `v*` tag workflow and all required checks pass, runtime images are published to GHCR.
3. The Staging workflow is manually run with the selected tag and executes the staging smoke test.

Docker images must use the same version tag:

```text
notezy-core:v0.1.0-beta.1
notezy-client-gateway:v0.1.0-beta.1
notezy-notification:v0.1.0-beta.1
```

## Tag Management Rules

- Use annotated tags, not lightweight tags.
- Treat pushed tags as immutable; never point an existing tag at another commit.
- Do not overwrite a published version or reuse a deleted version number.
- Roll back by deploying an older version tag; do not move the current tag.
- Use `Release v...` as the GitHub Release title.
- Release notes should record major features, fixes, breaking changes, migrations, and deployment notes.
- Production deployments must use an explicit version tag instead of `latest`.

## Common Commands

```bash
# List local tags
git tag

# List remote tags
git ls-remote --tags origin

# Show the commit and annotation for a tag
git show v0.1.0-beta.1

# Push one tag
git push origin v0.1.0-beta.1
```

Git itself does not enforce the tag format; the format above is the Notezy
project convention. The GitHub Actions `v*` trigger is an automation setting.
If the project later requires stricter validation, the workflow must validate
the tag format explicitly as well.
