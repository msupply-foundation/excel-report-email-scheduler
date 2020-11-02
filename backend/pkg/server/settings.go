package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	dbstore "github.com/grafana/simple-datasource-backend/pkg/db"
)

func (server *HttpServer) updateSettings(rw http.ResponseWriter, request *http.Request) {
	var settings dbstore.Settings
	requestBody, err := request.GetBody()
	bodyAsBytes, _ := ioutil.ReadAll(requestBody)
	err = json.Unmarshal(bodyAsBytes, &settings)

	if err != nil {
		http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
		return
	}

	server.db.CreateOrUpdateSettings(settings)

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) fetchSettings(rw http.ResponseWriter, request *http.Request) {
	config := server.db.GetSettings()
	json.NewEncoder(rw).Encode(config)
	rw.WriteHeader(http.StatusOK)

}
