package api

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
)

func TestSwaggerRedirects(t *testing.T) {
	os.Mkdir("../public/docs/", os.ModeDir | os.ModePerm)
	os.Create("../public/docs/test-swagger.json")
	defer os.RemoveAll("../public/docs/test-swagger.json")
	os.Setenv("SWAGGER_ENDPOINT", "/docs/test-swagger.json")
	ts := httptest.NewServer(router(Default))
	defer ts.Close()
	url := fmt.Sprintf("%v/api/", ts.URL)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("Error getting URL: %v", err)
	} else if res.StatusCode != http.StatusOK {
		t.Errorf("Did not return correct status code. Expected %v got %v", http.StatusOK, res.StatusCode)
	}

	v, exists := res.Request.URL.Query()["url"]
	if !exists {
		t.Error("Redirect query parameter not supplied")
	} else if len(v) != 1 && v[0] != ts.URL + "/docs/test-swagger.json" {
		t.Error("Invalid redirect query parameter")
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}
}
