package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	GrafanaUsername string `json:"grafanaUsername"`
	GrafanaPassword string `json:"grafanaPassword"`
	Email           string `json:"email"`
	EmailPassword   string `json:"emailPassword"`
}

func getHttpHandler(sqliteDatasource *SQLiteDatasource) backend.CallResourceHandler {
	mux := http.NewServeMux()

	mux.HandleFunc("/settings", getUpdateSettingsHandler(sqliteDatasource))
	return httpadapter.New(mux)
}

func getUpdateSettingsHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {

		var config Config
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &config)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.createOrUpdateSettings(config)

		rw.WriteHeader(http.StatusOK)
	}
}
