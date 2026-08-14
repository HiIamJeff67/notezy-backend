package main

import (
	"os"

	notification "github.com/HiIamJeff67/notezy-backend/internal/notification"
)

func main() {
	application := notification.NewApplication()
	if len(os.Args) == 1 {
		shutdown := application.Start()
		defer shutdown()
	}

	Execute()
}
