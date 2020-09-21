package main

import (
	"os"

	"github.com/unravela/artisan/artisan"
)

func main() {
	ws, _ := artisan.OpenWorkspace("../artisan-monorepo-demo")
	err := ws.Run("//apps/shop-ui/backend:build")

	if err != nil {
		os.Exit(-1)
	}
}
