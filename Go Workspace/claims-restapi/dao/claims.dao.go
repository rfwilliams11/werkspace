package dao

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Claim model
type Claim struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	Name   string        `bson:"name" json:"name"`
	Source string        `bson:"sourcee" json:"sourcee"`
	Event  string        `bson:"event" json:"event"`
}

type ClaimsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "claims"
)

func (m *ClaimsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *ClaimsDAO) FindAll() ([]Claim, error) {
	var claims []Claim
	err := db.C(COLLECTION).Find(bson.M{}).All(&claims)
	return claims, err
}

func (m *ClaimsDAO) FindById(id string) (Claim, error) {
	var claim Claim
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&claim)
	return claim, err
}

func (m *ClaimsDAO) Insert(claim Claim) error {
	err := db.C(COLLECTION).Insert(&claim)
	return err
}

func (m *ClaimsDAO) Delete(claim Claim) error {
	err := db.C(COLLECTION).Remove(&claim)
	return err
}

func (m *ClaimsDAO) Update(claim Claim) error {
	err := db.C(COLLECTION).UpdateId(claim.ID, &claim)
	return err
}
