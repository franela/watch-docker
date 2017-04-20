package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/franela/watch-docker/timeline"
)

func GetPulls(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	fmt.Println(query)

	prs, err := timeline.GetProjectTimeline("moby", 20, 10, "2017-04-20T12:55:46Z")
	if err != nil {
		return
	}

	rw.Header().Add("content-type", "application/json")
	json.NewEncoder(rw).Encode(prs)
}
