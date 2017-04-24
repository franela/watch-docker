package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/franela/watch-docker/timeline"
)

func GetPulls(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	skip := r.URL.Query().Get("skip")
	curate := r.URL.Query().Get("curate") == "true"

	prs, err := timeline.GetProjectTimeline(query, 21, skip, curate)
	if err != nil {
		return
	}

	rw.Header().Add("content-type", "application/json")
	json.NewEncoder(rw).Encode(prs)
}

func CuratePull(rw http.ResponseWriter, r *http.Request) {
	var curationReq timeline.CurationRequest

	if err := json.NewDecoder(r.Body).Decode(&curationReq); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := timeline.CuratePullRequest(curationReq); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
