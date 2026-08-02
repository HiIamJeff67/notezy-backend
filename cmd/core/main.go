package main

import (
	"os"

	core "github.com/HiIamJeff67/notezy-backend/internal/services/core"
)

func main() {
	if len(os.Args) == 1 {
		shutdown := core.Start()
		defer shutdown()
	}

	Execute()
}
