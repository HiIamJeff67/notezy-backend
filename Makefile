# ============================== Database Shortcut Commands ============================== #
view-hotreload-dbs:
	docker compose exec -T notezy-core go run ./internal/core/commands viewDatabases

view-hotreload-enums:
	docker compose exec -T notezy-core go run ./internal/core/commands viewAllEnums

psql:
	docker exec -it notezy-db psql -U jeff -d notezy-db

# ============================== Migration Commands ============================== #
migrate-build-db:
	docker compose exec -T notezy-core ./core migrateDB
migrate-hotreload-db:
	docker compose exec -T notezy-core go run ./internal/core/commands migrateDB

clear-build-db:
	docker exec -i notezy-db psql -U jeff -d notezy-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
clear-hotreload-db: # the same as the build version of db
	docker exec -i notezy-db psql -U jeff -d notezy-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

remigrate-build-db:
	make clear-build-db
	make migrate-build-db

remigrate-hotreload-db:
	make clear-hotreload-db
	make migrate-hotreload-db

# ============================== Seeding Commands ============================== #
seed-build-db:
	docker compose exec -T notezy-core ./core seedDB
seed-hotreload-db:
	docker compose exec -T notezy-core go run ./internal/core/commands seedDB

clear-go-cache:
	go clean -modcache
	cd contracts && go mod download
	cd shared && go mod download
	cd internal/core && go mod download
	cd internal/gateway && go mod download
	cd internal/durablejob && go mod download
	cd internal/email && go mod download
	cd internal/realtimegateway && go mod download
	cd test && go mod download

test-auth-e2e:
	docker compose exec -T notezy-gateway sh -c 'cd test && go test ./e2e/auth'

test-architecture:
	cd test && go test ./architecture

test-all:
	cd test && go test ./...

# ============================== GraphQL Shortcut Commands ============================== #
gql-generate: # update before generate
	go run github.com/99designs/gqlgen@v0.17.76 generate --config contracts/core/v1/graphql/gqlgen.yaml

gql-clean:
ifeq ($(OS),Windows_NT)
	@if exist contracts\core\v1\graphql\generated\*.* del /q /s contracts\core\v1\graphql\generated\*.*
	@if exist contracts\core\v1\graphql\models\*.* del /q /s contracts\core\v1\graphql\models\*.*
else
	rm -rf contracts/core/v1/graphql/generated/*
	rm -rf contracts/core/v1/graphql/models/*
endif

gql-regenerate:
	make gql-clean
	make gql-generate
