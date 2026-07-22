# UsersToRoutineTags Removal Contract

`UsersToRoutineTags` has been removed from the Station and Routine domain.

RoutineTag is now a private property owned directly by its creator. The backend assigns ownership from the authenticated user when a tag is created; ownership is not represented as a shareable relationship and is not sent by frontend clients.

This contract is a frontend migration guide. It complements [RoutineTag Private Ownership Contract](routine-tag-private-ownership-contract.md).

## Removed Client Concepts

Remove all frontend code and local persistence that represents any of the following:

- `UsersToRoutineTags` entity, table, type, schema, or sync event.
- `userId + tagId + permission` tag relationship records.
- RoutineTag permission values: `Owner`, `Admin`, `Write`, and `Read`.
- RoutineTag sharing, invitation, member management, permission editor, or collaborator picker UI.
- Any cache query that resolves accessible tags through a user-to-tag join.

There is no replacement sharing API for RoutineTag.

## Current Ownership Model

```text
Authenticated User --owns--> RoutineTag
Authenticated User --can edit--> Routine
RoutineTag --private link--> Routine
```

A shared Routine can have private tag links from multiple collaborators. Each collaborator sees only the links made with their own tags.

The effective link identity is maintained by the backend:

```text
routineId + routineTagId + authenticated user
```

Do not send `userId`, `ownerId`, or `stationId` in RoutineTag create, update, link, or unlink payloads.

## REST Migration

RoutineTag endpoints keep their existing names and bodies:

```text
GET    /api/development/v1/routineTag/getMyRoutineTagById
GET    /api/development/v1/routineTag/getAllMyRoutineTags
POST   /api/development/v1/routineTag/createRoutineTag
POST   /api/development/v1/routineTag/createRoutineTags
PUT    /api/development/v1/routineTag/updateMyRoutineTagById
PUT    /api/development/v1/routineTag/updateMyRoutineTagsByIds
DELETE /api/development/v1/routineTag/hardDeleteMyRoutineTagById
DELETE /api/development/v1/routineTag/hardDeleteMyRoutineTagsByIds
```

The caller receives only their own tags. A tag ID created by another user behaves as not found. Do not interpret this as a tag permission state or render a permission-recovery UI.

Routine link endpoints are unchanged:

```text
POST /api/development/v1/routine/linkRoutineTagById
POST /api/development/v1/routine/linkRoutineTagsByIds
```

Keep sending only the documented routine/tag IDs and `isUnlink`. `isUnlink: true` removes only the current user's link.

## Routine Response Semantics

`Routine.tagIds` is no longer a global Routine property from the frontend's point of view.

It is the authenticated user's view of that Routine's private tag links:

```json
{
  "id": "routine-uuid",
  "stationId": "station-uuid",
  "tagIds": ["my-routine-tag-uuid"]
}
```

Never assume a missing tag ID means that another collaborator removed a global tag. It only means the current user has no remaining link for that tag and Routine.

When a user's `UsersToStations` membership is removed, the backend deletes that user's RoutineTag links in the affected Station. Refetch or invalidate Routine data after observing access revocation.

## GraphQL Migration

GraphQL schema names remain unchanged:

- `searchRoutineTags(input)` returns only the caller's tags.
- `searchRoutines(input)` returns `tagIds` only for the caller's links.
- `SearchRoutineInput.tagIds` accepts only the caller's tag IDs as meaningful filters.

Passing another user's tag ID to `searchRoutines` returns no matching link. It does not reveal the tag and does not produce a tag-permission workflow.

## Local Database and Cache Migration

1. Drop the client-side `UsersToRoutineTags` table or collection.
2. Remove its foreign keys, indexes, synchronizers, mutations, and optimistic updates.
3. Remove tag permission fields from local RoutineTag types.
4. Scope RoutineTag caches and persisted state to the authenticated account.
5. Invalidate Routine query caches after account switching, RoutineTag hard deletion, tag link/unlink, or Station access revocation.

The frontend may keep a local `RoutinesToTags` relation for UI state, but it must be treated as current-user state. The server remains the authority for link ownership and Station access.
