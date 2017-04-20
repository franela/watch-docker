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
	}
	c := mongoSession.DB("github").C("pulls")
	for _, index := range indexes {
		c.EnsureIndex(index)
	}
}

func GetProjectTimeline(nameQuery string, size int, importance int, skipToken string) ([]*github.PullRequest, error) {
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
	var query bson.M
	if skipToken != "" {
		query = bson.M{
			"comments":           bson.M{"$gt": importance},
			"mergedat":           bson.M{"$lt": skipTime},
			"base.repo.fullname": bson.M{"$regex": fmt.Sprintf(".*%s.*", nameQuery)},
		}
	} else {
		query = bson.M{
			"comments":           bson.M{"$gt": importance},
			"base.repo.fullname": bson.M{"$regex": fmt.Sprintf(".*%s.*", nameQuery)},
		}
	}
	if err = c.Find(query).Sort("-mergedat").All(&prs); err != nil {
		return nil, err
	}

	return prs, nil
}

func CleanUp() {
	mongoSession.Close()
}
