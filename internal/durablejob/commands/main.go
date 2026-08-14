package main

import (
	"os"

	durablejob "github.com/HiIamJeff67/notezy-backend/internal/durablejob"
)

func main() {
	application := durablejob.NewApplication()
	if len(os.Args) == 1 {
		shutdown := application.Start()
		defer shutdown()
	}

	Execute()
}
