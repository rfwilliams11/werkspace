package database

import (
	"testing"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/testapi"
)

// func TestReadyAgainstRealDatabase(t *testing.T) {
// 	d, err := GetDatabase()
// 	defer d.Close()
// 	expected := true
// 	actual, err := d.Ready()
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	} else if actual != expected {
// 		t.Errorf("Database not ready, expected: %t got: %t", expected, actual)
// 	}
// }

func TestReady(t *testing.T) {
	db := testapi.InitMongo()
	d := Database{db.GetDatabase()}
	defer testapi.ResetMongo(db)

	expected := true
	actual, err := d.Ready()
	d.Close()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	} else if actual != expected {
		t.Errorf("Database not ready, expected: %t got: %t", expected, actual)
	}
}
