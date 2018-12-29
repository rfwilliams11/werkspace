package testapi

import (
	"os"
	"github.com/globalsign/mgo/dbtest"
	"github.com/globalsign/mgo"
	"io/ioutil"
)

const (
	dbTempDir = "rtr-advancement-api-mongodb"
	dbName = "rtr-advancement-api-test"
)

type TestDatabase struct {
	Server  *dbtest.DBServer
	Session *mgo.Session
}

func (t *TestDatabase) GetDatabase() *mgo.Database {
	return t.Session.Copy().DB(dbName)
}

func (t *TestDatabase) Insert(collection string, doc interface{}) {
	err := t.Session.DB(dbName).C(collection).Insert(doc)
	if err != nil {
		panic(err)
	}
}

func (t *TestDatabase) Drop() {
	t.Session.DB(dbName).DropDatabase()
}

func InitMongo() *TestDatabase {
	tempDir, _ := ioutil.TempDir(os.TempDir(), dbTempDir)
	server := dbtest.DBServer{}
	server.SetPath(tempDir)
	session := server.Session()

	return &TestDatabase{&server, session}
}

func ResetMongo(database *TestDatabase) {
	database.Drop()
	database.Session.Close()
	database.Server.Stop()
	os.RemoveAll(os.TempDir() + dbTempDir)
}
