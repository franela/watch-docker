package timeline

import (
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongoSession *mgo.Session

type CurationRequest struct {
	Repo    string `json:"repo"`
	Number  int    `json:"number"`
	Curated *bool  `json:"curated"`
}

func init() {
	mongoUrl := "mongo"
	if url, exists := os.LookupEnv("MONGO_URL"); exists {
		mongoUrl = url
	}
	var err error
	for {
		mongoSession, err = mgo.Dial(mongoUrl)
		if err == nil {
			break
		}
		fmt.Printf("MongoDB not accessible, trying again in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	indexes := []mgo.Index{
		mgo.Index{
			Key:        []string{"base.repo.fullname"},
			Background: true,
		},
		mgo.Index{
			Key:        []string{"mergedat"},
			Background: true,
		},
		mgo.Index{
			Key:        []string{"comments"},
			Background: true,
		},
		mgo.Index{
			Key:        []string{"base.repo.fullname", "number"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		},
	}

	c := mongoSession.DB("github").C("pulls")
	for _, index := range indexes {
		c.EnsureIndex(index)
	}
}

func GetProjectTimeline(nameQuery string, size int, skipToken string, curate bool) ([]*github.PullRequest, error) {
	c := mongoSession.DB("github").C("pulls")

	format := "2006-01-02T15:04:05Z"
	var skipTime time.Time
	var err error
	if skipToken != "" {
		skipTime, err = time.Parse(format, skipToken)
		if err != nil {
			skipToken = ""
		}
	}

	prs := []*github.PullRequest{}
	query := bson.M{
		"curated":            true,
		"base.repo.fullname": bson.M{"$regex": fmt.Sprintf(".*%s.*", nameQuery)},
	}
	if curate {
		query["curated"] = nil
	}
	if skipToken != "" {
		query["margedat"] = bson.M{"$lt": skipTime}
	}
	if err = c.Find(query).Sort("-mergedat").Limit(size).All(&prs); err != nil {
		return nil, err
	}

	return prs, nil
}

func CuratePullRequest(r CurationRequest) error {
	c := mongoSession.DB("github").C("pulls")

	query := bson.M{"number": r.Number, "base.repo.fullname": r.Repo}
	change := bson.M{"$set": bson.M{"curated": r.Curated}}

	return c.Update(query, change)
}

func CleanUp() {
	mongoSession.Close()
}
