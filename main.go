package main

import (
	"os"

	app "github.com/HiIamJeff67/notezy-backend/app"
)

func main() {
	// Check `os.Args` to distinguish the HTTP server from internal Cobra commands.
	app.StartApplication(len(os.Args) > 1)
}
