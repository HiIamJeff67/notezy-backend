package main

import (
	"os"

	email "github.com/HiIamJeff67/notezy-backend/internal/email"
)

func main() {
	if len(os.Args) == 1 {
		shutdown := email.Start()
		defer shutdown()
	}

	Execute()
}
