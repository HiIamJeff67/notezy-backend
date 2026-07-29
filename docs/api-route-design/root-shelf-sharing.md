# RootShelf Sharing API Design

## Scope

RootShelf is the collaboration and permission boundary for its descendant
SubShelves, Materials, BlockPacks, Blocks, and Yjs documents. `BlockPack` does
not expose an independent sharing model.

## Effective permission model

| Permission | Read scope | Realtime channel | Manage Read/Write | Manage Admin | Transfer ownership |
| --- | --- | --- | --- | --- | --- |
| `Owner` | yes | write | yes | yes | yes |
| `Admin` | yes | write | yes | no | no |
| `Write` | yes | write | no | no | no |
| `Read` | yes | read | no | no | no |

Only the Owner may grant, revoke, or downgrade `Admin`. An Admin may manage
only `Read` and `Write`. Generic permission CRUD never modifies or deletes the
Owner row; ownership transfer is a separate Owner-only operation.

`UsersToShelves` is the persistence authority for membership and effective
permission. RootShelf ownership, membership rows, and owner-scoped quota are
updated atomically by the ownership-transfer workflow.

## REST surface

All routes are rooted at `/api/development/v1/root-shelves` and use the
standard authenticated response/error pipeline.

| Method | Path | Operation |
| --- | --- | --- |
| `GET` | `/:rootShelfId/permissions/:userPublicId` | Read one membership permission. |
| `POST` | `/:rootShelfId/permissions/:userPublicId` | Create one membership permission. |
| `PUT` | `/:rootShelfId/permissions/:userPublicId` | Upsert one membership permission. |
| `PUT` | `/:rootShelfId/permissions` | Upsert multiple membership permissions. |
| `PATCH` | `/:rootShelfId/permissions/:userPublicId` | Update one membership permission. |
| `DELETE` | `/:rootShelfId/permissions/:userPublicId` | Remove one membership permission. |
| `DELETE` | `/:rootShelfId/permissions` | Remove multiple membership permissions atomically. |
| `POST` | `/:rootShelfId/ownership` | Transfer ownership. |
| `DELETE` | `/:rootShelfId/memberships/me` | Leave one RootShelf. |
| `DELETE` | `/memberships/me` | Leave multiple RootShelves. |

Bulk permission requests require unique public IDs and are all-or-nothing. URI
identifiers are authoritative; binders map them into request DTOs rather than
trusting duplicated body identifiers.

## Transactional invariants

Permission mutations lock the target RootShelf and affected membership rows.
Ownership transfer additionally locks the old and new owner membership rows
and both `UserAccount` rows. The target must already hold a non-owner
membership (at least `Read`); on success, the old owner becomes `Admin`, the
target becomes the unique owner, and quota accounting remains balanced in the
same transaction.

Self-leave must remove only the authenticated user's relationship and follows
the ownership-transfer guard: an Owner cannot leave while still being the
owner. Ownership must first be transferred successfully.

## Realtime interaction

Channel ticket issuance, subscription, and every Yjs document mutation verify
the active descendant hierarchy and the required effective permission. A
revoked permission yields `permission_revoked`; an inactive/deleted hierarchy
yields `resource_unavailable`. Both detach the affected channel.

Cross-Gateway lifecycle fanout is intentionally deferred to the Kafka event
architecture. Until then, an idle channel is revalidated when it subscribes or
submits a document mutation.

## Read models

GraphQL `searchRootShelves` is the root-sharing read model. It exposes the
owner public ID, non-owner sharer public IDs, and the caller's effective
permission. Presence remains a separate, ephemeral Realtime concern.
