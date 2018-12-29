package env

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)


type DecryptionClient interface {
	Call(url string, username string, password string, query SplitPeaRequest) (response string, err error)
}

type SplitPeaClient struct {
	http.Client
}

// SplitPeaRequest is a data structure that ensures the proper formatting of the POST body sent to the SplitPea API.
// Currently Key2 is a command line parameter passed as KEY2.
type SplitPeaRequest struct {
	Euuid         string   `json:"eUUID"`
	KeyComponents []string `json:"keyComponents"`
}

//CallSplitPea forms an HTTP Post using the supplied arguments and sends the Post to the SplitPea API
//If the call is successful, the response will be what was stored in Splitpea that matches the EUUID and key components.
//Normally this is a password or a database connection string.
func (s SplitPeaClient) Call(url string, username string, password string, query SplitPeaRequest) (response string, err error) {
	if(len(url)==0 || len(username)==0 || len(password)==0){
		log.WithFields(log.Fields{
			"error": fmt.Sprintf("DecryptUrl=%s Username=%s PasswordLength=%d", url, username, len(password)),
		}).Fatal("Incomplete SplitPea Components")
	}

	body, _ := json.Marshal(query)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	resp, err := s.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("SplitPea HTTP Request Error")
		return
	}

	if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"error": fmt.Sprintf("%d when decrypting secret", resp.StatusCode),
		}).Fatal("SplitPea Decrypt Error")
		return
	}

	intermediateBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	result, _ := base64.StdEncoding.DecodeString(string(intermediateBytes))
	response = string(result)
	return
}
