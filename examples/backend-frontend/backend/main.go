package main

import (
	"log"
	"net/http"

	_ "example.com/backend/statik"
	"github.com/rakyll/statik/fs"
)

func main() {

	statikFS, _ := fs.New()

	fileServer := http.FileServer(statikFS)
	http.Handle("/", fileServer)

	log.Println("Running on http://localhost:8089")
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		log.Fatal(err)
	}
}
