# ============================== Database Shortcut Commands ============================== #
view-hotreload-dbs:
	docker compose exec -T notezy-api go run ./cmd/api viewDatabases

view-hotreload-enums:
	docker compose exec -T notezy-api go run ./cmd/api viewAllEnums

psql:
	docker exec -it notezy-db psql -U jeff -d notezy-db

# ============================== Migration Commands ============================== #
migrate-build-db:
	docker compose exec -T notezy-api ./api migrateDB
migrate-hotreload-db:
	docker compose exec -T notezy-api go run ./cmd/api migrateDB

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
	docker compose exec -T notezy-api ./api seedDB
seed-hotreload-db:
	docker compose exec -T notezy-api go run ./cmd/api seedDB

clear-go-cache:
	go clean -modcache
	go mod download

test-auth-e2e:
	docker compose exec -T notezy-api go test ./test/e2e/auth

test-architecture:
	go test ./test/architecture

# ============================== GraphQL Shortcut Commands ============================== #
gql-generate: # update before generate
	go get github.com/99designs/gqlgen@v0.17.76
	go run github.com/99designs/gqlgen generate --config infra/graphql/gqlgen.yaml

gql-clean:
ifeq ($(OS),Windows_NT)
	@if exist internal\platform\graphql\generated\*.* del /q /s internal\platform\graphql\generated\*.*
	@if exist internal\platform\graphql\models\*.* del /q /s internal\platform\graphql\models\*.*
else
	rm -rf internal/platform/graphql/generated/*
	rm -rf internal/platform/graphql/models/*
endif

gql-regenerate:
	make gql-clean
	make gql-generate
