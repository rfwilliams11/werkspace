package main

import (
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/api/handlers"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/rtrapi"
	log "github.com/sirupsen/logrus"
)

func main() {
	splitPea := env.SplitPeaClient{http.Client{Timeout: 30 * time.Second}}
	env.Initialize(splitPea)

	rtrapi.Start(router)
}

func router(r *mux.Router) *mux.Router {

	// These generic handlers can be used to directly access the database
	log.Info("Mongo Endpoints Configured")
	r.HandleFunc("mgo/{collection}/count", handlers.CollectionCountHandler).Methods("GET")
	r.HandleFunc("mgo/{collection}", handlers.CollectionHandler).Methods("GET")

	return r
}
