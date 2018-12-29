package env

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCallSuccess(t *testing.T) {
	query := SplitPeaRequest{"euuid", []string{"key1", "key2"}}
	expected := "secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, base64.StdEncoding.EncodeToString([]byte(expected)))
	}))
	defer server.Close()

	client := SplitPeaClient{}
	res, _ := client.Call(server.URL, "username", "password", query)

	if res != expected {
		t.Error("Did not call url")
	}
}

func TestCallFailure(t *testing.T) {
	query := SplitPeaRequest{"euuid", []string{"key1", "key2"}}
	expected := http.StatusUnauthorized
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failure", expected)
	}))
	defer server.Close()

	client := SplitPeaClient{}
	_, err := client.Call(server.URL, "username", "password", query)

	if !strings.HasPrefix(err.Error(), strconv.Itoa(expected)) {
		t.Error("No error calling splitpea")
	}
}

