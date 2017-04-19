package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPulls(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	org := vars["org"]

	pull := r.URL.Query().Get("pull")
	fmt.Println(org, repo, pull)
}
