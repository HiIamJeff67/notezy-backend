package main

import "os"

func main() {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	parser := NewContractParser(root)
	endpoints := parser.Parse()
	NewEndpointValidator(parser).Validate(endpoints)
	NewArtifactGenerator(root, parser).Generate(endpoints)
}
