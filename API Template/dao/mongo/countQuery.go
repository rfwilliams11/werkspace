package mongo

import (
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/database"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
)

func CountDocuments(collectionName string, matching map[string][]string) int {
	count := -1

	database.Exec(collectionName, func(c *mgo.Collection){
		mGoCount, err := c.Find(Match(matching)).Count()
		if err == nil {
			count = mGoCount
		} else {
			log.Warn(err)
		}
	})

	return count
}