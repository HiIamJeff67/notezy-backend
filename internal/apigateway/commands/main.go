package main

import (
	"os"

	apigateway "github.com/HiIamJeff67/notegic-backend/internal/apigateway"
)

func main() {
	application := apigateway.NewApplication()
	if len(os.Args) == 1 {
		shutdown := application.Start()
		defer shutdown()
	}

	Execute()
}
