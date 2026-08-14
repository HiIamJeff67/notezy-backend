# Origin and security rules

- Send the API key in `X-API-Key`; never put it in a URL, query string, or request body.
- Cross-origin browser access works only for origins configured in the APIGateway allowlist.
- CLI and server clients should omit Origin and Referer rather than forge an allowed browser origin.
- HTTPS is mandatory outside local development.
- Never commit Postman environments after adding credentials.
- Never place account passwords, API keys, or authorization codes in URLs.
- Treat all IDs and permissions returned by the client as untrusted; server authorization remains authoritative.

Public documentation does not itself grant a third-party origin permission to call the service.
