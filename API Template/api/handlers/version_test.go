package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
)

func TestVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Version))
	defer ts.Close()

	env.Version = "test version"
	env.Branch = "test branch"
	env.SHA1 = "test sha1"
	env.BuildTime = "test build time"

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

	ver := &version{}
	err = json.Unmarshal(body, &ver)
	if err != nil {
		t.Errorf("Error unmarshalling body: %v", err)
	}

	expected := &version{
		"test version",
		"test branch",
		"test sha1",
		"test build time",
	}
	if !reflect.DeepEqual(ver, expected) {
		t.Errorf("Wrong version - expected: %v, got : %v", expected, ver)
	}
}
