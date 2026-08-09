package main

import (
	"os"

	notification "github.com/HiIamJeff67/notezy-backend/internal/notification"
)

func main() {
	if len(os.Args) == 1 {
		shutdown := notification.Start()
		defer shutdown()
	}

	Execute()
}
