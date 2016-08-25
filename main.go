package main

import (
	"net/http"

	"github.com/takonews/takonews-api/config/routes"
	_ "github.com/takonews/takonews-api/db/migrations"
)

func main() {
	// routing
	mux := http.NewServeMux()
	mux.Handle("/", routes.Router())

	// run server
	if err := http.ListenAndServe(":8000", mux); err != nil {
		panic(err)
	}
}
