# Gateway v1 transport contracts

Gateway owns the private transport envelope used when it calls an internal
service. `Request[D]` and `Response[D]` carry version, operation, metadata, and
a service-owned DTO. `header.go` defines Gateway-managed private response
headers used when Core refreshes browser credentials.

Gateway does not own Core domain RequestDto/ResponseDto contracts. Those live
under `contracts/core/v1/api/`; Gateway passes them as the envelope's `Dto`.
