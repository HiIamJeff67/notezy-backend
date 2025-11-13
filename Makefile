# ============================== Database Shortcut Commands ============================== #
view-dbs:
	docker-compose run --rm notezy-api go run main.go viewDatabases

migrate-build-db:
	docker-compose run --rm notezy-api ./notezy-backend migrateDB
migrate-hotreload-db:
	docker-compose run --rm notezy-api go run main.go migrateDB

clear-build-db:
	docker exec -i notezy-db psql -U jeff -d notezy-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
clear-hotreload-db: # the same as the build version of db
	docker exec -i notezy-db psql -U jeff -d notezy-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

clear-go-cache:
	go clean -modcache
	go mod download

test-auth-e2e:
	docker-compose run --rm notezy-api go test ./test/e2e/auth

# ============================== GraphQL Shortcut Commands ============================== #
gql-generate: # update before generate
	go get github.com/99designs/gqlgen@v0.17.76
	go run github.com/99designs/gqlgen generate

gql-clean:
ifeq ($(OS),Windows_NT)
	@if exist app\graphql\generated\*.* del /q /s app\graphql\generated\*.*
	@if exist app\graphql\models\*.* del /q /s app\graphql\models\*.*
else
	rm -rf app/graphql/generated/*
	rm -rf app/graphql/models/*
endif

gql-regenerate:
	make gql-clean
	make gql-generate

gql-check:
	go run github.com/99designs/gqlgen generate --verbose