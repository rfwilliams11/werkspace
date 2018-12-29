package handlers

import (
	"encoding/json"
	"net/http"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/database"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
	log "github.com/sirupsen/logrus"
)

// Health response containing the status of the mongo database
//
// swagger:model healthResponse
type health struct {
	DatabaseReady bool   `json:"databaseReady"`
	Reason        string `json:"reason,omitempty"`
	Endpoint      string `json:"endpoint"`
}

//swagger:operation GET /health getHealth
// Returns the health status of the mongo database and the api
// ---
// consumes:
//   - "*/*"
// produces:
//   - "application/json"
// responses:
//   200:
//     description: "Valid response"
//     schema:
//       $ref: "#/definitions/healthResponse"
func Health(w http.ResponseWriter, _ *http.Request) {
	url := env.Config().MongoUrl()
	var h = health{Endpoint: url}
	db, err := database.GetDatabase()
	if db != nil {
		defer db.Close()
		_, err = db.Ready()
	}

	if err != nil {
		h.Reason = err.Error()
	} else {
		h.DatabaseReady = true
		log.Info("Database Ready")
	}

	json.NewEncoder(w).Encode(h)
}
