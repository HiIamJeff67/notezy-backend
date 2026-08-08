package main

import (
	"os"

	gateway "github.com/HiIamJeff67/notezy-backend/internal/gateway"
)

func main() {
	if len(os.Args) == 1 {
		shutdown := gateway.Start()
		defer shutdown()
	}

	Execute()
}
