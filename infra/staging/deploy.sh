#!/usr/bin/env sh

set -eu

: "${IMAGE_TAG:?IMAGE_TAG is required}"
: "${IMAGE_REGISTRY:?IMAGE_REGISTRY is required}"

root_directory=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
compose_file=${COMPOSE_FILE:-"$root_directory/docker-compose.prod.yaml"}
compose_env_file=${COMPOSE_ENV_FILE:-/etc/notezy/staging.env}

export GATEWAY_IMAGE="$IMAGE_REGISTRY/notezy-gateway:$IMAGE_TAG"
export CORE_IMAGE="$IMAGE_REGISTRY/notezy-core:$IMAGE_TAG"
export DURABLE_JOB_IMAGE="$IMAGE_REGISTRY/notezy-durablejob:$IMAGE_TAG"
export EMAIL_IMAGE="$IMAGE_REGISTRY/notezy-email:$IMAGE_TAG"
export REALTIME_GATEWAY_IMAGE="$IMAGE_REGISTRY/notezy-realtimegateway:$IMAGE_TAG"
export YJS_WORKER_IMAGE="$IMAGE_REGISTRY/notezy-yjsworker:$IMAGE_TAG"

docker compose --env-file "$compose_env_file" -f "$compose_file" pull
docker compose --env-file "$compose_env_file" -f "$compose_file" up -d --no-build --remove-orphans
