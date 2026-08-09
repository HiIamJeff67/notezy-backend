# ============================== Workspace Commands ============================== #

WORKSPACE_MODULES := contracts shared internal/cli internal/core internal/gateway internal/durablejob internal/email internal/notification internal/realtimegateway test

.PHONY: ci-format ci-vet ci-unit ci-race ci-generated ci-containers staging-deploy staging-smoke kafka-topics

ci-format:
	@files="$$(find contracts shared internal test -type f -name '*.go' -not -path '*/vendor/*' -print0 | xargs -0 gofmt -l)"; \
	if [ -n "$$files" ]; then echo "Unformatted Go files:"; echo "$$files"; exit 1; fi

ci-vet:
	@set -e; \
	for module in $(WORKSPACE_MODULES); do \
		echo "go vet $$module"; \
		(cd "$$module" && GOWORK=off go vet ./...); \
	done

ci-unit:
	$(MAKE) test-all

ci-race:
	@set -e; \
	for module in $(WORKSPACE_MODULES); do \
		$(MAKE) test-race MODULE="$$module"; \
	done

ci-generated:
	$(MAKE) -C contracts gql-generate
	@git diff --exit-code -- contracts/core/v1/graphql/generated contracts/core/v1/graphql/models

ci-containers:
	@set -e; \
	for runtime in gateway core durablejob email notification realtimegateway yjsworker; do \
		echo "docker build $$runtime"; \
		target=production; \
		if [ "$$runtime" = yjsworker ]; then target=runtime; fi; \
		docker build --target "$$target" --file "internal/$$runtime/Dockerfile" --tag "notezy-ci-$$runtime" .; \
	done

staging-deploy:
	infra/staging/deploy.sh

staging-smoke:
	infra/staging/smoke.sh

CLI_RUN := go -C internal/cli run .

test-all:
	$(CLI_RUN) test-all

test-module:
	@if [ -z "$(MODULE)" ]; then echo "usage: make test-module MODULE=<module>"; exit 1; fi
	$(CLI_RUN) test-module $(MODULE)

test-race:
	@if [ -z "$(MODULE)" ]; then echo "usage: make test-race MODULE=<module>"; exit 1; fi
	$(CLI_RUN) test-race $(MODULE)

test-contracts:
	$(MAKE) -C contracts test

test-shared:
	$(MAKE) -C shared test

test-core:
	$(MAKE) -C internal/core test

test-gateway:
	$(MAKE) -C internal/gateway test

test-durable-job:
	$(MAKE) -C internal/durablejob test

test-email:
	$(MAKE) -C internal/email test

test-realtime-gateway:
	$(MAKE) -C internal/realtimegateway test

test-notification:
	$(MAKE) -C internal/notification test

test-architecture:
	$(MAKE) -C test test-architecture

test-integration:
	$(MAKE) -C test test-integration

test-integration-kafka:
	$(MAKE) -C test test-integration-kafka

test-load-websocket:
	$(MAKE) -C test test-load-websocket

test-soak-websocket:
	$(MAKE) -C test test-soak-websocket

test-load-kafka-lag:
	$(MAKE) -C test test-load-kafka-lag

kafka-topics:
	$(CLI_RUN) kafka topics ensure
