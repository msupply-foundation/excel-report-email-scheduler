package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	grafanaUsername `json:"grafanaUsername"`
	grafanaPassword `json:"grafanaPassword"`
	email           `json:"email"`
	emailPassword   `json:"grafanaEmailPassword"`
}

func getHttpHandler(sqliteDatasource *SQLiteDatasource) backend.CallResourceHandler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", getHelloHandler(sqliteDatasource))
	return httpadapter.New(mux)
}

func getHelloHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		log.DefaultLogger.Info("handleHello")

		if request.Body != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(request.Body)
			newStr := buf.String()
			log.DefaultLogger.Info("Request body:", newStr)
		}

		var config Config
		err := json.NewDecoder(request.Body).Decode(&config)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword } ", http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusOK)

	}
}
