package database

import (
	"sync"
	"time"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"

	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	*mgo.Database
}

var (
	instance *mgo.Session
	once     sync.Once
)

// Make sure to call Database.Close for each successful call to GetDatabase
//
// It is recommended to make a new GetDatabase call for each group of database calls.
// i.e. if a handler needs to hit the database 5 times, all calls should be done in one session.
//
// A copy of the session is made for each call. This is generally considered a best practice
// using mgo. Compared to having a single session, the pitfall here would be the overhead of a
// new session for every set of calls. The gain comes from not having a user blocking the db
// access of another user due to a poor performance call.
var GetDatabase = func() (*Database, error) {
	var err error
	once.Do(func() {
		instance, err = newSession()
	})
	if err != nil {
		once = sync.Once{}
		log.Errorf("Error pinging database: %v", err)
		return nil, err
	}
	return &Database{instance.Copy().DB(env.Config().Mongo.DbName)}, err
}

func Exec(collectionName string, f func(*mgo.Collection)) {
	db, err := GetDatabase()
	if err != nil {
		log.Errorf("Error executing database query: %v", err)
	}
	defer db.Close()
	c := db.C(collectionName)
	f(c)
}

func ExecWithTimeout(collectionName string, responseTimeout time.Duration, f func(*mgo.Collection) error) error {
	db, err := GetDatabase()
	if err != nil {
		log.Errorf("Error executing database query: %v", err)
		return err
	}
	db.Session.SetSyncTimeout(responseTimeout)
	db.Session.SetSocketTimeout(responseTimeout)
	defer db.Close()
	c := db.C(collectionName)
	return f(c)
}

func newSession() (*mgo.Session, error) {
	c := env.Config()
	s, err := mgo.DialWithTimeout(c.Mongo.Url, time.Duration(c.Mongo.DialTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return s, err
}

func (db *Database) Ready() (bool, error) {
	_, err := db.Session.DatabaseNames()
	if err != nil {
		log.Errorf("Error pinging database: %v", err)
		return false, err
	}
	return true, nil
}

func (db *Database) Close() {
	db.Session.Close()
}
