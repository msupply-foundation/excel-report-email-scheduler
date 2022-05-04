package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/setting"
	"io/ioutil"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (server *HttpServer) updateSettings(rw http.ResponseWriter, request *http.Request) {
	var settings setting.Settings
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateSettings: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyAsBytes, &settings)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: json.Unmarshal: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = server.db.CreateOrUpdateSettings(settings)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: db.CreateOrUpdateSettings: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(rw).Encode(settings)
	if err != nil {
		log.DefaultLogger.Error("updateSettings: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
