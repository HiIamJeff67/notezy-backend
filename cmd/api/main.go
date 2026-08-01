package main

import (
	"os"

	commands "github.com/HiIamJeff67/notezy-backend/cmd/api/commands"
	gateway "github.com/HiIamJeff67/notezy-backend/internal/gateway"
	core "github.com/HiIamJeff67/notezy-backend/internal/services/core"
)

func main() {
	if len(os.Args) > 1 {
		commands.Execute()
		return
	}

	go core.Start()
	core.WaitUntilReady()
	gateway.Start()
}
