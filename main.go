package main

import (
	"net/http"
	"os"

	"github.com/takonews/takonews-api/config/routes"
	_ "github.com/takonews/takonews-api/db/migrations"
)

func main() {
	// routing
	mux := http.NewServeMux()
	mux.Handle("/", routes.Router())

	// run server
	if err := http.ListenAndServeTLS(":8000", os.Getenv("HOME")+"/.ssh/serverpub.key", os.Getenv("HOME")+"/.ssh/serverpriv.key", mux); err != nil {
		panic(err)
	}
}
