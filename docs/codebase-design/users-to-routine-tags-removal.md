# UsersToRoutineTags Removal Record

## Decision

`UsersToRoutineTags` was removed from the Station and Routine domain. It
represented shareable tag membership and permission, which conflicted with the
current design: RoutineTag is a user-private aggregate and is owned directly
by its creator.

## Removed model

The backend no longer maintains any of the following:

* a `UsersToRoutineTags` schema, repository, relation, or synchronization
  event;
* tag-level `Owner`, `Admin`, `Write`, or `Read` permissions;
* tag-level invitations, member management, or authorization checks; or
* an accessibility query derived from a user-to-tag join.

## Replacement model

Ownership is derived from the authenticated user during tag creation. Routine
tag links are scoped by `(routineId, routineTagId, userId)`. Routine sharing is
still governed by Station membership; it permits a user to manage their own
tag links when they have the required Routine permission, but never exposes
another user's tags or links.

## Data and read-model effects

`Routine.tagIds` is a caller-scoped projection. It may differ between two
collaborators reading the same Routine. `searchRoutineTags` returns only
caller-owned tags, and `searchRoutines` treats supplied tag IDs as meaningful
only when they belong to the caller.

Removing a user's `UsersToStations` membership deletes that user's tag links
for Routines in the Station. This is an expected cascade, not a tag-permission
state transition.

For the active API and persistence rules, see
[RoutineTag Private Ownership Design](routine-tag-private-ownership.md).
