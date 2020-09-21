package main

import (
	"github.com/unravela/artisan/artisan"
)

func main() {
	ws, err := artisan.OpenWorkspace("../artisan-monorepo-demo")
	if err != nil {
		panic(err)
	}

	err = ws.Run("//apps/shop-ui/frontend:build")
	if err != nil {
		panic(err)
	}
}
