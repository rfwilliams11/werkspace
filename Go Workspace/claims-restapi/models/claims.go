package models

import "gopkg.in/mgo.v2/bson"

//Claim model
type Claim struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	Name   string        `bson:"name" json:"name"`
	Source string        `bson:"sourcee" json:"sourcee"`
	Event  string        `bson:"event" json:"event"`
}
