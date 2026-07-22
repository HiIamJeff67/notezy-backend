# RoutineTag Private Ownership Contract

RoutineTag is a private user property. A tag belongs to exactly one authenticated user and is never shared, granted, or administered by another user.

Routine and Station sharing remain unchanged. A user may link one of their own tags to any Routine they can edit. The resulting link is private to that user: another collaborator on the same Routine does not receive that tag ID and cannot inspect, edit, or remove the tag.

All REST responses retain the normal envelope:

```json
{ "success": true, "data": {}, "exception": null }
```

## Frontend Data Rules

- Do not build a RoutineTag sharing, permission, member, or admin UI.
- Do not send `ownerId` when creating or updating a tag. The server always assigns the authenticated user as owner.
- Scope RoutineTag client state by the authenticated user. Clear or refetch this state after an account switch.
- A Routine's `tagIds` contains only links created by the current user. It is not a complete list of all collaborators' tags.
- When a user loses access to a Station, the backend removes that user's links to its Routines. Treat missing tag IDs after a routine refresh as the expected state, not as a client-side conflict.

## REST APIs

Base path:

```text
/api/development/v1/routineTag
```

### Read Tags

```text
GET /getMyRoutineTagById?routineTagId=UUID&isDeleted=false
GET /getAllMyRoutineTags?areDeleted=false
```

`isDeleted` and `areDeleted` should default to `false` in the frontend schema. RoutineTag is hard-deleted, so requesting deleted tags does not return historical tag records.

Each tag has this shape:

```json
{
  "id": "UUID",
  "name": "Computer Science",
  "color": "#2F80ED",
  "icon": "BookOpen",
  "updatedAt": "RFC3339 timestamp",
  "createdAt": "RFC3339 timestamp"
}
```

The API returns only tags owned by the authenticated user. A tag ID belonging to another user behaves as not found.

### Create and Update Tags

```text
POST /createRoutineTag
POST /createRoutineTags
PUT /updateMyRoutineTagById
PUT /updateMyRoutineTagsByIds
DELETE /hardDeleteMyRoutineTagById
DELETE /hardDeleteMyRoutineTagsByIds
```

Request payloads and response shapes are unchanged. There is no owner field in the contract; ownership is derived exclusively from authentication.

`color` remains a six-digit hex string such as `#2F80ED`. `icon` is optional. A hard-deleted tag also removes all of its routine links.

## Linking a Private Tag to a Routine

Base path:

```text
/api/development/v1/routine
```

Single link or unlink:

```text
POST /linkRoutineTagById
```

```json
{
  "routineId": "UUID",
  "routineTagId": "UUID",
  "isUnlink": false
}
```

Bulk link or unlink:

```text
POST /linkRoutineTagsByIds
```

```json
{
  "linkedRoutinesAndTags": [
    { "routineId": "UUID", "routineTagId": "UUID" }
  ],
  "isUnlink": false
}
```

The caller must still have `Owner`, `Admin`, or `Write` access to each Routine. The selected tag must be one returned by that caller's RoutineTag API. `isUnlink: true` removes only the caller's link; it never removes another collaborator's private tag link.

The backend derives `userId` and `stationId` for each link. Frontend clients must not send either value.

## Routine Responses

The existing Routine response remains unchanged:

```json
{
  "id": "UUID",
  "stationId": "UUID",
  "title": "Review distributed systems notes",
  "tagIds": ["UUID"],
  "taskIds": [],
  "itemIds": []
}
```

For all Routine read APIs and GraphQL routine search, `tagIds` are filtered to the authenticated user. Use `getAllMyRoutineTags` to resolve those IDs into presentation metadata.

## GraphQL Search

`searchRoutineTags(input)` returns only the caller's tags. Its GraphQL type and query shape are unchanged.

`searchRoutines(input)` continues to accept `tagIds`. They are interpreted as the caller's own tags. Passing another user's tag ID produces no matching routines rather than exposing tag data or returning a tag-permission workflow.
