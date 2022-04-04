package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"excel-report-email-scheduler/pkg/dbstore"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (server *HttpServer) fetchSettings(rw http.ResponseWriter, request *http.Request) {

	settings, err := server.db.GetSettings()
	if err != nil {
		log.DefaultLogger.Error("fetchSettings: db.GetSettings(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(settings)
	if err != nil {
		log.DefaultLogger.Error("fetchSettings: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) updateSettings(rw http.ResponseWriter, request *http.Request) {
	var settings dbstore.Settings
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateSettings: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	defer request.Body.Close()
	if err != nil {
		log.DefaultLogger.Error("updateSettings: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.Unmarshal(bodyAsBytes, &settings)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.SettingsFields()).Error(), http.StatusBadRequest)
		panic(err)
	}

	err = server.db.CreateOrUpdateSettings(settings)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: db.CreateOrUpdateSettings: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(settings)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}
