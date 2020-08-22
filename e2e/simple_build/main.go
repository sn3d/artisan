package main

import (
	"github.com/unravela/delvin/workspace"
	"os"
)

func main() {
	ws, _ := workspace.Open("./examples/backend-frontend")
	err := ws.Run("//frontend:build")

	if err != nil {
		os.Exit(-1)
	}
}
