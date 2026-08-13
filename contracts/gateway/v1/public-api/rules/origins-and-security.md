# Origin and security rules

- Browser requests must use `credentials: include` so cookies are accepted and returned.
- Cross-origin browser access works only for origins configured in the Gateway allowlist.
- CLI and server clients should omit Origin and Referer rather than forge an allowed browser origin.
- HTTPS is mandatory outside local development.
- Never commit Postman environments after adding credentials.
- Never place account passwords, cookies, CSRF tokens, realtime tickets, or authorization codes in URLs.
- Realtime connection and channel tickets expire after five minutes and are single-use.
- Treat all IDs and permissions returned by the client as untrusted; server authorization remains authoritative.

Public documentation does not itself grant a third-party origin permission to call the service.
