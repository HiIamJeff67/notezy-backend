# Gateway v1 transport contracts

Gateway owns the private transport envelope used when it calls an internal
service. `Request[D]` and `Response[D]` carry version, operation, metadata, and
a service-owned DTO. `Request[D].Tokens` carries authentication credentials as
typed data; Gateway extracts cookies before the call and never forwards the
Cookie header to Core. `ClientRequest[D]` and `ClientResponse[D]` describe the
public Gateway envelope; `ClientRequest[D]` serializes as the DTO body itself,
without internal metadata. Public responses never contain access or refresh
tokens; only non-sensitive refresh metadata may be present. Internal
`Response[D]` may carry a typed `Tokens` envelope for Gateway interception.

Gateway does not own Core domain RequestDto/ResponseDto contracts. Those live
under `contracts/core/v1/api/`; Gateway passes them as the envelope's `Dto`.
