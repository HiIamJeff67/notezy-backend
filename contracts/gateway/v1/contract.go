package gatewaycontract

const Version = "v1"

type RequestEnvelope interface {
	GetVersion() string
	GetOperation() string
	GetMetadata() RequestMetadata
}
