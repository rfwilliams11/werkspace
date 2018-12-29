package mongo

import (
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/database"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func Documents(collectionName string, matching map[string][]string) []bson.M {
	var results []bson.M

	database.Exec(collectionName, func(c *mgo.Collection){
		c.Pipe(Query(Match(matching), nil, "")).All(&results)
	})

	return results
}

func ProjectedDocuments(collectionName string, matching map[string][]string, projecting []string) []bson.M {
	var results []bson.M

	database.Exec(collectionName, func(c *mgo.Collection){
		c.Pipe(Query(Match(matching), projecting, "")).All(&results)
	})

	return results
}

func ProjectedDistinctDocuments(collectionName string, matching map[string][]string, projecting []string, distinctField string) []bson.M {
	var results []bson.M

	database.Exec(collectionName, func(c *mgo.Collection){
		c.Pipe(Query(Match(matching), projecting, distinctField)).All(&results)
	})

	return results
}