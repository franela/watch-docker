package main

import (
	"context"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type Pull struct {
	Id int64 `json:id`
}

func main() {
	token := ""

	mongo_url := "mongo"
	if url, exists := os.LookupEnv("MONGO_URL"); exists {
		mongo_url = url
	}
	session, err := mgo.Dial(mongo_url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("github").C("pulls")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.PullRequestListOptions{State: "closed"}
	for {
		prs, resp, err := client.PullRequests.List(ctx, "moby", "moby", opts)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, pr := range prs {
			if pr.MergedAt == nil {
				continue
			}
			pr, _, err := client.PullRequests.Get(ctx, "moby", "moby", *pr.Number)
			if err != nil {
				log.Println(err)
				continue
			}
			err = c.Insert(pr)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Inserted PR [%d]\n", *pr.Number)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
}
