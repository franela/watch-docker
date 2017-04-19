package main

import (
	"log"
	"net/http"

	"github.com/franela/watch-docker/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{org}/{repo}", handlers.GetPulls).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(r)

	log.Fatal(http.ListenAndServe(":3000", n))
}
