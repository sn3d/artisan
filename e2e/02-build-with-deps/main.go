package main

import (
	"os"

	"github.com/unravela/artisan/workspace"
)

func main() {
	ws, _ := workspace.Open("../artisan-monorepo-demo")
	err := ws.Run("//apps/shop-ui/backend:build")

	if err != nil {
		os.Exit(-1)
	}
}
