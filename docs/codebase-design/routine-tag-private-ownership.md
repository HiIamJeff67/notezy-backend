# RoutineTag Private Ownership Design

## Domain model

`RoutineTag` is a private user-owned aggregate. It belongs to exactly one
authenticated user and has no permission, membership, invitation, or sharing
model. The server derives ownership from authentication on creation; ownership
is not an input field.

A user may link one of their own tags to a Routine that they can modify. The
link is private to that user, even when the Routine is shared through its
Station. Different collaborators can therefore maintain independent tag links
on the same Routine without observing or modifying one another's tags.

```text
User --owns--> RoutineTag
User --has write access to--> Routine
User --creates private link--> RoutineTag <-> Routine
```

The effective link identity is `routineId + routineTagId + userId`.

## Persistence and authorization

RoutineTag ownership is stored directly on the tag. `UsersToRoutineTags` is
not part of the current schema or authorization model. A caller can read,
update, or hard-delete only their own tags; another user's tag behaves as not
found.

Routine/Station sharing does not grant tag visibility. When a user loses
Station access, the backend removes that user's RoutineTag links for affected
Routines. A tag hard delete removes all links created by its owner.

## REST surface

Routes are rooted at `/api/development/v1/routine-tags`.

| Method | Path | Operation |
| --- | --- | --- |
| `GET` | `/:routineTagId` | Read one owned tag. |
| `GET` | `/` | Read owned tags. |
| `POST` | `/` | Create one tag. |
| `POST` | `/batch` | Create tags in bulk. |
| `PUT` | `/:routineTagId` | Update one tag. |
| `PUT` | `/batch` | Update tags in bulk. |
| `DELETE` | `/:routineTagId/permanently` | Hard-delete one tag. |
| `DELETE` | `/batch/permanently` | Hard-delete tags in bulk. |

Routine link operations remain part of the Routine resource:

| Method | Path | Operation |
| --- | --- | --- |
| `POST` | `/api/development/v1/routines/:routineId/tags/:routineTagId` | Link or unlink one private tag. |
| `POST` | `/api/development/v1/routines/tags` | Link or unlink tags in bulk. |

The Route permission policy requires `Owner`, `Admin`, or `Write` on each
Routine. Services derive user and Station identity from the authenticated
context and existing domain records; they do not accept these identities as
tag-link input.

## Read semantics

Routine `tagIds` is caller-scoped: it contains only tags linked by the current
user, not a global list of collaborators' tags. GraphQL `searchRoutineTags`
returns only caller-owned tags, and `searchRoutines(input)` interprets tag
filters only within the caller's owned-tag scope.
