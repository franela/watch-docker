package timeline

import (
	"fmt"
	"time"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/google/go-github/github"
)

type Pull struct {
	Id int64 `json:id`
}


func GetProjectTimeline(nameQuery string, size int, skipToken string) ([]*github.PullRequest, error) {
	mongo_url := "mongo"
	if url, exists := os.LookupEnv("MONGO_URL"); exists {
		mongo_url = url
	}
	session, err := mgo.Dial(mongo_url)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("github").C("pulls")

	format := "2006-01-02T15:04:05.000Z"
	var skipTime time.Time
	if skipToken != "" {
		skipTime, err = time.Parse(format, skipToken)
		if err != nil {
			skipToken = ""
		}
	}

	prs := []*github.PullRequest{}
	var query bson.M
	if skipToken != "" {
		query = bson.M{"mergedat": bson.M{"$lt": skipTime}, "base.repo.fullname": bson.M{"$regex": fmt.Sprintf(".*%s.*", nameQuery)}}
	} else {
		query = bson.M{"base.repo.fullname": bson.M{"$regex": fmt.Sprintf(".*%s.*", nameQuery)}}
	}
	
	if err = c.Find(query).Sort("-mergedat").Limit(size).All(&prs); err != nil {
		return nil, err
	}
	return prs, nil
}
