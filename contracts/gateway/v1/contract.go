package gatewaycontract

type RequestEnvelope interface {
	GetVersion() string
	GetOperation() string
	GetMetadata() RequestMetadata
}
