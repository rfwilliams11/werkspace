package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	dao "bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/dao/mongo"
	"github.com/globalsign/mgo/bson"
	"strings"
)

func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	var documents []bson.M

	collection := mux.Vars(r)["collection"]
	project := r.URL.Query().Get("select")

	if len(project) > 0 {
		cleanedQuery := r.URL.Query()
		cleanedQuery.Del("select")

		distinct := cleanedQuery.Get("distinct")

		if len(distinct) > 0 {
			cleanedQuery.Del("distinct")
			documents = dao.ProjectedDistinctDocuments(collection, cleanedQuery, strings.Split(project, ","),distinct)
		} else {
			documents = dao.ProjectedDocuments(collection, cleanedQuery, strings.Split(project, ","))
		}
	} else {
		documents = dao.Documents(collection, r.URL.Query())
	}

	json.NewEncoder(w).Encode(documents)
}

type countResponse struct {
	Collection  string `json:"collection"`
	Count   	int `json:"count"`
}

func CollectionCountHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := mux.Vars(r)["collection"]


	count := countResponse{
		collectionName,
		dao.CountDocuments(collectionName, r.URL.Query()),
	}

	json.NewEncoder(w).Encode(count)
}