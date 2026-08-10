module github.com/HiIamJeff67/notezy-backend/test

go 1.26.0

require (
	github.com/HiIamJeff67/notezy-backend/contracts v0.0.0
	github.com/HiIamJeff67/notezy-backend/shared v0.0.0
	github.com/google/uuid v1.6.0
	github.com/twmb/franz-go v1.21.5
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/comail/colog v0.0.0-20160416085026-fba8e7b1f46c // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.13.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.42.0 // indirect
	go.opentelemetry.io/otel/log v0.18.0 // indirect
	go.opentelemetry.io/otel/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
)

replace github.com/HiIamJeff67/notezy-backend/contracts => ../contracts

replace github.com/HiIamJeff67/notezy-backend/shared => ../shared
