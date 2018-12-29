package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/database"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/testapi"
)

func TestHealthAgainstRealDatabase(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Health))
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Errorf("Error getting URL: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("Did not return correct status code. Expected %v got %v", http.StatusOK, res.StatusCode)
	}
}

func TestHealthOk(t *testing.T) {
	db := testapi.InitMongo()
	defer testapi.ResetMongo(db)
	database.GetDatabase = func() (*database.Database, error) {
		return &database.Database{Database: db.GetDatabase()}, nil
	}
	ts := httptest.NewServer(http.HandlerFunc(Health))
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Errorf("Error getting URL: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("Did not return correct status code. Expected %v got %v", http.StatusOK, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	h := health{}
	err = json.Unmarshal(body, &h)
	expected := health{true, "", env.Config().MongoUrl()}
	if err != nil {
		t.Errorf("Error unmarshalling body: %v", err)
	} else if h != expected {
		t.Errorf("Received unexpected body. Expected: %v, Actual: %v", expected, h)
	}
}

func TestHealthNotOk(t *testing.T) {
	expectedError := "Expected database error"
	database.GetDatabase = func() (*database.Database, error) {
		return nil, errors.New(expectedError)
	}
	ts := httptest.NewServer(http.HandlerFunc(Health))
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Errorf("Error getting URL: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("Did not return correct status code. Expected %v got %v", http.StatusOK, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	h := health{}
	err = json.Unmarshal(body, &h)
	expected := health{false, expectedError, env.Config().MongoUrl()}
	if err != nil {
		t.Errorf("Error unmarshalling body: %v", err)
	} else if h != expected {
		t.Errorf("Received unexpected body. Expected: %v, Actual: %v", expected, h)
	}
}
