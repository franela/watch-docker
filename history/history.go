package main

import (
	"context"
	"log"

	mgo "gopkg.in/mgo.v2"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type Pull struct {
	Id int64 `json:id`
}

func fetchHistory() {
	token := ""
	session, err := mgo.Dial("localhost")
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

	/*
		lastPull = Pull{}
		err := c.Find(nil).Sort("-_id").One(&lastPull)
		if err != nil {
			log.Fatal("Could not get the last pull request from DB", err)
		}
	*/

	for {
		opts := &github.PullRequestListOptions{State: "closed"}
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
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
}
