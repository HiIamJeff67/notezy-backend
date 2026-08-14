# ClientGateway v1 public contract

This directory contains the public, runtime-specific ClientGateway contract for the web/client surface. ClientGateway retains the existing access/refresh JWT and cookie behavior, including account registration/login and client-only routes.

It is separate from the externally advertised APIGateway integration API. That contract is generated at [`contracts/api-gateway/v1/public/`](../../../api-gateway/v1/public/).

The remaining legacy `contracts/gateway/v1/public-api/` artifacts are kept temporarily for migration compatibility and must not be linked from public documentation.
