# CI/CD and Jenkins Pipeline

NOT-43 defines the division of responsibilities between GitHub Actions and
Jenkins. Both may invoke only commands already exposed by the root `Makefile`
or `internal/cli`; a pipeline must not duplicate a second test workflow.

## GitHub Actions

`.github/workflows/ci.yml` is the quality gate for every pull request, pushes to
`main`/`refactor/**`, and version tags. It includes:

- formatting, `go vet`, unit tests, and `go test -race` for every Go module;
- a clean-tree check after regenerating GraphQL artifacts;
- production container builds for the five Go runtimes and the Yjs Worker;
- publishing runtime images to GitHub Container Registry after a version tag
  passes every gate.

`.github/workflows/integration.yml` runs Docker-backed integration tests only
when manually triggered or scheduled. It first starts the PostgreSQL, Redis,
and Kafka services defined by
`infra/docker/docker-compose.integration.yaml` through the root
`Makefile`'s `compose-integration-up`, runs root integration and
Core/DurableJob broker-flow tests, and finally cleans up with
`compose-integration-down`. Local development can use the same root targets;
when Compose is already running, `make test-integration` or
`make test-integration-kafka` runs tests without taking over container
management.

`.github/workflows/staging.yml` is a promotion workflow requiring approval from
the `staging` environment. On a self-hosted runner labeled `staging`, it runs
`infra/docker/docker-compose.prod.yaml` with the selected GHCR tag, promotes
immutable images through `infra/staging/deploy.sh`, and checks every runtime's
`/startedz` and `/healthz` with `infra/staging/smoke.sh`. Compose accepts
`GATEWAY_IMAGE`, `CORE_IMAGE`, `DURABLE_JOB_IMAGE`, `EMAIL_IMAGE`,
`REALTIME_GATEWAY_IMAGE`, and `YJS_WORKER_IMAGE`, so promotion does not rebuild
images. The staging runner must provide environment settings in
`/etc/notezy/staging.env`; that file is not checked out from the repository and
secrets must not be committed. Compose logs after deployment are uploaded as an
artifact with 14-day retention.

## Jenkins

The root `Jenkinsfile` is the delivery pipeline for a self-hosted agent:

1. checkout, format, vet, unit, race, and generated-contract gates;
2. an optional production container build;
3. promotion of immutable images and the same smoke script on an agent labeled
   `notezy-staging`, controlled by the `DEPLOY_STAGING` and `IMAGE_TAG`
   parameters;
4. optional integration tests controlled by `RUN_INTEGRATION`. Tests share
   `infra/docker/docker-compose.integration.yaml` with GitHub Actions, and the
   Compose stack is cleaned up afterward.

The Jenkins staging agent must provide Docker, the Compose plugin, Git, and a
Jenkins credential that can access GHCR. Staging deployments may promote only
immutable image tags that have passed CI; the staging agent must not recompile
source.

Jenkins does not replace the GitHub Actions PR gate, and this issue does not
perform a production rollout, migration rollback, or disaster recovery.
Image publishing, environment secrets, and deployment permissions must be
managed by Jenkins credentials or GitHub OIDC and must not be committed to the
repository.

## Local Equivalent Commands

```sh
make ci-format
make ci-vet
make ci-unit
make ci-race
make ci-generated
make ci-containers
```

The staging runner's delivery commands are:

```sh
IMAGE_REGISTRY=ghcr.io/ORG/REPO IMAGE_TAG=TAG \
COMPOSE_ENV_FILE=/etc/notezy/staging.env make staging-deploy

make staging-smoke
```

Integration tests use `make test-integration` and
`make test-integration-kafka`; these commands only execute tests and require
Docker-backed dependencies to be available. Use
`make test-integration-managed` for the complete lifecycle, or run
`compose-integration-up`, the test commands, and `compose-integration-down`
separately.
