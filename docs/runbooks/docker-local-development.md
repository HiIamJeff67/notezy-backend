# Docker local development

This runbook describes the supported local Compose workflows. The root
`docker-compose.yaml` is the development stack. The production-like stack is
maintained at
`infra/docker/docker-compose.prod.yaml` and uses the same runtime boundaries,
health checks, and dependency ordering as the deployment configuration.

## Development stack

Start the development stack in the background and wait for its health checks:

```sh
make compose-up
```

`make compose-up` decrypts `secrets/envs/.env.enc` with SOPS into a mode
`0600` temporary file, passes that file to Docker Compose, and removes it when
Compose exits. Raw `docker compose up` does not invoke SOPS; it only reads the
root `.env` file and bypasses the encrypted environment flow.

This is preferred for local development and CI-style shell sessions:

- `--build` rebuilds images when source or Dockerfile changes.
- `-d` detaches from the long-running service logs.
- `--wait` waits for services to become running or healthy and returns a
  non-zero status when they cannot reach that state.

Inspect status and logs with:

```sh
docker compose ps
docker compose logs -f notezy-client-gateway
docker compose logs -f notezy-api-gateway
```

### Hot reload

The Go runtimes use Air in development. ClientGateway, Core, DurableJob, Email,
Notification, and RealtimeGateway run `air -c .air.toml` from their own module
directory; APIGateway currently runs `go run ./commands` from its module. Source
changes are rebuilt and restarted automatically;
changes to dependencies, environment variables, Dockerfiles, or Compose
configuration still require an image rebuild and container recreation.

The Yjs Worker uses `pnpm dev` (`tsx watch`) instead of Air. Production images
for every runtime use the compiled entrypoint only and do not run a watcher.

Stop the development stack with:

```sh
make compose-down
```

## Production-like local validation

The production Compose file can be exercised locally; deployment is not
required to validate its Compose wiring. Use a separate project name to isolate
its Compose network and lifecycle. The file currently sets explicit
`container_name` values, so stop the development stack first or provide
distinct `DOCKER_*_SERVICE_NAME` values if both stacks must run simultaneously:

```sh
COMPOSE_PROJECT_NAME=notezy-prod-local \
COMPOSE_FILE=infra/docker/docker-compose.prod.yaml \
COMPOSE_ENCRYPTED_ENV_FILE=secrets/envs/.env.production.enc \
make compose-up

docker compose \
  --project-name notezy-prod-local \
  --project-directory . \
  --env-file .env \
  -f infra/docker/docker-compose.prod.yaml \
  ps
```

The production-like file currently wires the database, Redis, Kafka-dependent
runtimes, Yjs worker, Gateway, Realtime Gateway, and Nginx. The health checks
and `depends_on` conditions ensure that a dependent service is not started
until its required runtime reports healthy. Application runtime health checks
run every 60 seconds; dependency checks for databases, Redis, and Kafka retain
their shorter intervals so startup ordering remains responsive.

Run the staging smoke checks against this local stack when the local
environment contains all required settings:

```sh
COMPOSE_FILE=infra/docker/docker-compose.prod.yaml \
COMPOSE_PROJECT_NAME=notezy-prod-local \
COMPOSE_ENCRYPTED_ENV_FILE=secrets/envs/.env.production.enc \
SOPS_CONFIG_FILE=.sops.yaml \
make staging-smoke
```

Clean up the production-like stack with the same project name and file:

```sh
COMPOSE_PROJECT_NAME=notezy-prod-local \
COMPOSE_FILE=infra/docker/docker-compose.prod.yaml \
COMPOSE_ENCRYPTED_ENV_FILE=secrets/envs/.env.production.enc \
make compose-down
```

## What local validation proves

Local production-like Compose validation covers:

- Compose interpolation and service names.
- Image builds and runtime entrypoints.
- Dependency ordering and `startedz`/`healthz` checks.
- Runtime-to-runtime DNS and internal port wiring.
- Basic startup and smoke behavior.

It cannot prove production-only concerns such as registry pulls, external
managed databases or Kafka, secret rotation, TLS certificates, DNS/load
balancers, resource limits, network policies, or multi-host failure behavior.
Those require a staging deployment before production rollout.
