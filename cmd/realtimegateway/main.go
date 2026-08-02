package main

import (
	"os"

	realtimegateway "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway"
)

func main() {
	if len(os.Args) == 1 {
		shutdown := realtimegateway.Start()
		defer shutdown()
	}

	Execute()
}
