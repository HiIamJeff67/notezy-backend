package main

import "fmt"

type ArtifactGenerator struct {
	root   string
	parser *ContractParser
}

func NewArtifactGenerator(root string, parser *ContractParser) *ArtifactGenerator {
	return &ArtifactGenerator{root: root, parser: parser}
}

func (g *ArtifactGenerator) Generate(endpoints []endpoint) {
	if g.parser == nil {
		panic("public API generator requires a contract parser")
	}
	writeGatewayArtifacts(g.root, endpoints)
	writeRealtimeArtifacts(g.root)
	writeDecisionRecord(g.root, endpoints)
	fmt.Printf("generated %d Gateway operations and RealtimeGateway HTTP/WebSocket contracts\n", len(endpoints))
}
