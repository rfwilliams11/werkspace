package handlers

import (
	"encoding/json"
	"net/http"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
)

// Version response containing build information of the current api
//
// swagger:model versionResponse
type version struct {
	Version   string `json:"version"`
	Branch    string `json:"branch"`
	Sha1      string `json:"sha1"`
	BuildTime string `json:"buildTime"`
}

//swagger:operation GET /version getVersion
// Returns the version information for the api
// ---
// consumes:
//   - "*/*"
// produces:
//   - "application/json"
// responses:
//   200:
//     description: "Valid response"
//     schema:
//       $ref: "#/definitions/versionResponse"
func Version(w http.ResponseWriter, _ *http.Request) {
	v := version{
		env.Version,
		env.Branch,
		env.SHA1,
		env.BuildTime,
	}
	json.NewEncoder(w).Encode(v)
}
