package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/franela/watch-docker/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{org}/{repo}", handlers.GetPulls).Methods("GET")

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.ServeFile(rw, r, "./www/index.html")
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www")))

	n := negroni.Classic()
	n.UseHandler(r)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), n))
}
