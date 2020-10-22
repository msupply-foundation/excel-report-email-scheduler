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

type ReportRecipient struct {
	ID         string `json:"id"`
	UserID     string `json:"userID"`
	ScheduleID string `json:"scheduleID"`
	Email      string `json:"email"`
}

func getHttpHandler(sqliteDatasource *SQLiteDatasource) backend.CallResourceHandler {
	mux := mux.NewRouter()

	mux.HandleFunc("/settings", getUpdateSettingsHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/settings", getAllUsersHandler(sqliteDatasource)).Methods("GET")

	mux.HandleFunc("/schedule", getCreateScheduleHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/schedule/{id}", getUpdateScheduleHandler(sqliteDatasource)).Methods("PUT")
	mux.HandleFunc("/schedule", getFetchSchedulesHandler(sqliteDatasource)).Methods("GET")

	mux.HandleFunc("/report-recipient/{id}", getFetchReportRecipientHandler(sqliteDatasource)).Methods("GET")
	mux.HandleFunc("/report-recipient", getFetchReportRecipientsHandler(sqliteDatasource)).Methods("GET")
	mux.HandleFunc("/report-recipient/{id}", getUpdateReportRecipientHandler(sqliteDatasource)).Methods("PUT")
	mux.HandleFunc("/report-recipient", getCreateReportRecipientHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/report-recipient/{id}", getDeleteReportRecipientHandler(sqliteDatasource)).Methods("DELETE")

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

func getFetchReportRecipientHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		reportRecipient := sqliteDatasource.getReportRecipient(id)

		json.NewEncoder(rw).Encode(reportRecipient)
		rw.WriteHeader(http.StatusOK)
	}
}

func getFetchReportRecipientsHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var recipients []ReportRecipient

		recipients = sqliteDatasource.getReportRecipients()

		json.NewEncoder(rw).Encode(recipients)
		rw.WriteHeader(http.StatusOK)
	}
}

func getCreateReportRecipientHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var recipient ReportRecipient
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &recipient)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.createReportRecipient(recipient)

		json.NewEncoder(rw).Encode(recipient)
		rw.WriteHeader(http.StatusOK)
	}
}

func getUpdateReportRecipientHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		var recipient ReportRecipient
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &recipient)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.updateReportRecipient(id, recipient)

		json.NewEncoder(rw).Encode(recipient)
		rw.WriteHeader(http.StatusOK)
	}
}

func getDeleteReportRecipientHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		// TODO: Handle error
		success, _ := sqliteDatasource.deleteReportRecipient(id)

		json.NewEncoder(rw).Encode(success)
		rw.WriteHeader(http.StatusOK)
	}
}
