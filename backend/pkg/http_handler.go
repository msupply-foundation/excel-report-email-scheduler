package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
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

type Schedule struct {
	ID             string `json:"id"`
	Interval       int    `json:"interval"`
	NextReportTime int    `json:"nextReportTime"`
	DashboardID    string `json:"dashboardID"`
}

func getHttpHandler(sqliteDatasource *SQLiteDatasource) backend.CallResourceHandler {
	mux := mux.NewRouter()

	mux.HandleFunc("/settings", getUpdateSettingsHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/settings", getAllUsersHandler(sqliteDatasource)).Methods("GET")

	mux.HandleFunc("/schedule", getCreateScheduleHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/schedule/{id}", getUpdateScheduleHandler(sqliteDatasource)).Methods("PUT")
	mux.HandleFunc("/schedule", getFetchSchedulesHandler(sqliteDatasource)).Methods("GET")

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

func getAllUsersHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		config := sqliteDatasource.getSettings()

		json.NewEncoder(rw).Encode(config)
		rw.WriteHeader(http.StatusOK)
	}
}

func getCreateScheduleHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var schedule Schedule
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &schedule)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.createSchedule(schedule)
	}
}

func getUpdateScheduleHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		var schedule Schedule
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &schedule)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.updateSchedule(id, schedule)
	}
}

func getFetchSchedulesHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var schedules []Schedule

		schedules = sqliteDatasource.getSchedules()

		json.NewEncoder(rw).Encode(schedules)
		rw.WriteHeader(http.StatusOK)
	}
}
