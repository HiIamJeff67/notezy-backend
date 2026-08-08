.PHONY: test test-race test-e2e view-databases view-enums psql migrate seed clear-db remigrate-db

test:
	go test ./...

test-race:
	go test -race ./...

test-e2e:
	go test ./test/e2e/...

view-databases:
	docker compose -f ../../docker-compose.yaml exec -T notezy-core go -C internal/core run ./commands viewDatabases

view-enums:
	docker compose -f ../../docker-compose.yaml exec -T notezy-core go -C internal/core run ./commands viewAllEnums

psql:
	docker exec -it notezy-db psql -U jeff -d notezy-db

migrate:
	docker compose -f ../../docker-compose.yaml exec -T notezy-core go -C internal/core run ./commands migrateDB

seed:
	docker compose -f ../../docker-compose.yaml exec -T notezy-core go -C internal/core run ./commands seedDB

clear-db:
	docker exec -i notezy-db psql -U jeff -d notezy-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

remigrate-db: clear-db migrate
