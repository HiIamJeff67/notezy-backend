#!/usr/bin/env sh

set -eu

root_directory=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
compose_file=${COMPOSE_FILE:-"$root_directory/docker-compose.prod.yaml"}
compose_env_file=${COMPOSE_ENV_FILE:-/etc/notezy/staging.env}

compose() {
	docker compose --env-file "$compose_env_file" -f "$compose_file" "$@"
}

check_http() {
	service=$1
	port=$2
	path=$3
	maximum_attempts=${SMOKE_MAXIMUM_ATTEMPTS:-30}
	attempt=1

	while [ "$attempt" -le "$maximum_attempts" ]; do
		if compose exec -T "$service" wget -q -O /dev/null "http://127.0.0.1:${port}${path}"; then
			echo "PASS ${service}${path}"
			return 0
		fi

		sleep "${SMOKE_RETRY_DELAY_SECONDS:-2}"
		attempt=$((attempt + 1))
	done

	echo "FAIL ${service}${path}" >&2
	return 1
}

check_yjs_worker() {
	maximum_attempts=${SMOKE_MAXIMUM_ATTEMPTS:-30}
	attempt=1

	while [ "$attempt" -le "$maximum_attempts" ]; do
		if compose exec -T notezy-yjs-worker node -e \
			"fetch('http://127.0.0.1:' + process.env.YJS_WORKER_PORT + '/healthz').then(response => process.exit(response.ok ? 0 : 1)).catch(() => process.exit(1))"; then
			echo "PASS notezy-yjs-worker/healthz"
			return 0
		fi

		sleep "${SMOKE_RETRY_DELAY_SECONDS:-2}"
		attempt=$((attempt + 1))
	done

	echo "FAIL notezy-yjs-worker/healthz" >&2
	return 1
}

check_http notezy-gateway "${DOCKER_GIN_PORT:-7777}" /startedz
check_http notezy-gateway "${DOCKER_GIN_PORT:-7777}" /healthz
check_http notezy-core 7778 /startedz
check_http notezy-core 7778 /healthz
check_http notezy-realtime-gateway 7779 /startedz
check_http notezy-realtime-gateway 7779 /healthz
check_http notezy-durable-job 8082 /startedz
check_http notezy-durable-job 8082 /healthz
check_http notezy-email 8081 /startedz
check_http notezy-email 8081 /healthz
check_yjs_worker
