.PHONY: test test-race gql-generate gql-clean gql-regenerate

test:
	go test ./...

test-race:
	go test -race ./...

gql-generate:
	go run github.com/99designs/gqlgen@v0.17.76 generate --config core/v1/graphql/gqlgen.yaml

gql-clean:
	rm -rf core/v1/graphql/generated/* core/v1/graphql/models/*

gql-regenerate: gql-clean gql-generate
