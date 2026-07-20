# RootShelf Sharing Contract

RootShelf 是分享與權限管理的邊界。授權某位使用者後，他對該 RootShelf 下的 SubShelf、Material 與 BlockPack 取得相同的有效存取範圍；BlockPack 不提供獨立分享設定。

## Permission Model

| Permission | Read data | Realtime channel | Manage Read/Write | Manage Admin | Remove Owner |
| --- | --- | --- | --- | --- | --- |
| `Owner` | yes | write | yes | yes | no |
| `Admin` | yes | write | yes | no | no |
| `Write` | yes | write | no | no | no |
| `Read` | yes | read | no | no | no |

Only an Owner can grant, revoke, or downgrade `Admin`. An Admin can create, update, and revoke only `Read` and `Write` permissions. The owner row cannot be modified or deleted through these APIs.

All endpoints use the normal authenticated REST pipeline and the usual `{ success, data, exception }` response envelope. Permission mutations lock the target RootShelf and affected permission rows, then perform their changes atomically.

## Permission APIs

Base path:

```text
/api/development/v1/rootShelf
```

### Upsert One Permission

```text
PUT /:rootShelfId/permissions/:userPublicId
```

```json
{ "permission": "Read" | "Write" | "Admin" }
```

Successful `data`:

```json
{
  "userPublicId": "UUID",
  "permission": "Read",
  "updatedAt": "RFC3339 timestamp",
  "createdAt": "RFC3339 timestamp"
}
```

### Upsert Many Permissions

```text
PUT /:rootShelfId/permissions
```

```json
{
  "permissions": [
    { "userPublicId": "UUID", "permission": "Read" },
    { "userPublicId": "UUID", "permission": "Write" }
  ]
}
```

`permissions` must contain 1 to 1024 unique `userPublicId` values. Successful `data.permissions` is the same per-user object as the single upsert response.

### Delete One Permission

```text
DELETE /:rootShelfId/permissions/:userPublicId
```

Success is HTTP `204 No Content`.

### Delete Many Permissions

```text
DELETE /:rootShelfId/permissions
```

```json
{ "userPublicIds": ["UUID", "UUID"] }
```

`userPublicIds` must contain 1 to 1024 unique values. Success is HTTP `204 No Content`. Missing users, missing permissions, a target owner, or an Admin attempting to manage another Admin are rejected as one atomic request; the API does not partially apply the batch.

## Sharing UI Data

Use GraphQL `searchUsers(input: SearchUserInput!)` to find invite targets. It returns public user data only, including `publicId`; never use an internal user id in a sharing request.

Use GraphQL `searchRootShelves(input: SearchRootShelfInput!)` to render current sharing state. Each `PrivateRootShelf` contains `ownerPublicId`, `sharerPublicIds`, and the caller's effective `permission`. `sharerPublicIds` excludes the owner.

Use the realtime participant endpoint for live presence, not for permission state:

```text
GET /api/development/v1/realtime/blockPacks/:blockPackId/participants
```

Only the RootShelf Owner and Admin can call it. `connectionCount` is the number of active root WebSocket connections for the user in that BlockPack.

## Realtime Lifecycle

Issuing a channel ticket and subscribing both validate the active RootShelf -> SubShelf -> BlockPack -> Yjs document hierarchy. When a RootShelf or SubShelf is soft-deleted, cascading triggers soft-delete affected BlockPacks; a BlockPack trigger mirrors the same `deletedAt` value into its BlockPackYjsDocument. Restoring the parent hierarchy restores the same rows and documents.

An already-open channel is revalidated before each outgoing `yjs-document` mutation. If the user was revoked, the gateway sends `permission_revoked`; if the resource hierarchy is unavailable, it sends `resource_unavailable`. In either case, the logical channel is removed and the frontend must destroy its editor provider. Awareness is ephemeral and is not revalidated per cursor frame.

Immediate cross-Gateway lifecycle fanout is not part of this monolith contract. It is deliberately deferred to the Kafka/microservice architecture work. Until then, an idle connection observes deletion or revocation when it next obtains a ticket, subscribes, or submits a document update.
