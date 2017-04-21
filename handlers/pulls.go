package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/franela/watch-docker/timeline"
)

func GetPulls(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	skip := r.URL.Query().Get("skip")

	prs, err := timeline.GetProjectTimeline(query, 21, 20, skip)
	if err != nil {
		return
	}

	rw.Header().Add("content-type", "application/json")
	json.NewEncoder(rw).Encode(prs)
}
