package main

import (
	"os"

	email "github.com/HiIamJeff67/notezy-backend/internal/email"
)

func main() {
	application := email.NewApplication()
	if len(os.Args) == 1 {
		shutdown := application.Start()
		defer shutdown()
	}

	Execute()
}
