# Integrations

This directory contains temporary integration documents used while the
frontend and backend are developing or migrating an API contract together.

These documents are intentionally separate from the permanent architecture
and system-design documentation. They may be revised, replaced, or removed
after the corresponding frontend integration is complete.

- [RoutineTask realtime lifecycle](routine-task-lifecycle-realtime.md):
  temporary WebSocket integration contract for running and completed task hints.
- [ClientGateway / APIGateway split](./client-api-gateway-split.md): Phase 5 boundaries, API key storage, request authentication, rate limiting, and runtime-specific public contract rules.
- [Frontend contract migration](./frontend-contract-migration.md): current
  ClientGateway/APIGateway boundary and RoutineTask monthly execution quota
  changes for the web client.
