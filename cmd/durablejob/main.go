package main

import (
	"os"

	durablejob "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob"
)

func main() {
	if len(os.Args) == 1 {
		shutdown := durablejob.Start()
		defer shutdown()
	}

	Execute()
}
