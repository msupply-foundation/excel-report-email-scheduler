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

type GroupSchedule struct {
	ID            string `json:"id"`
	ReportGroupID string `json:"reportGroupID"`
	ScheduleID    string `json:"scheduleID"`
}

type ReportGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReportGroupMembership struct {
	ID                string `json:"id"`
	ReportRecipientID string `json:"reportRecipientID"`
	ReportGroupID     string `json:"reportGroupID"`
}

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
	Name           string `json:"name"`
	Description    string `json:"description"`
	Lookback       int    `json:"lookback"`
}

type ReportRecipient struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
}

type ReportContent struct {
	ID         string `json:"id"`
	ScheduleID string `json:"scheduleID"`
	PanelID    int    `json:"panelID"`
	Lookback   int    `json:"lookback"`
	StoreID    string `json:"storeID"`
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

	mux.HandleFunc("/report-group", getFetchReportGroupHandler(sqliteDatasource)).Methods("GET")
	mux.HandleFunc("/report-group/{id}", getUpdateReportGroupHandler(sqliteDatasource)).Methods("PUT")
	mux.HandleFunc("/report-group", getCreateReportGroupHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/report-group/{id}", getDeleteReportGroupHandler(sqliteDatasource)).Methods("DELETE")

	mux.HandleFunc("/report-group-membership", getFetchReportGroupMembershipHandler(sqliteDatasource)).Queries("group-id", "{group-id}").Methods("GET")
	mux.HandleFunc("/report-group-membership", getCreateReportGroupMembershipHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/report-group-membership/{id}", getDeleteReportGroupMembershipHandler(sqliteDatasource)).Methods("DELETE")

	mux.HandleFunc("/report-content", getFetchReportContentHandler(sqliteDatasource)).Queries("schedule-id", "{schedule-id}").Methods("GET")
	mux.HandleFunc("/report-content", getCreateReportContentHandler(sqliteDatasource)).Methods("POST")
	mux.HandleFunc("/report-content/{id}", getUpdateReportContentHandler(sqliteDatasource)).Methods("PUT")
	mux.HandleFunc("/report-content/{id}", getDeleteReportContentHandler(sqliteDatasource)).Methods("DELETE")

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
		schedule, _ := sqliteDatasource.createSchedule()
		json.NewEncoder(rw).Encode(schedule)
		rw.WriteHeader(http.StatusOK)
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

func getFetchReportGroupHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var groups []ReportGroup

		groups = sqliteDatasource.getReportGroups()

		json.NewEncoder(rw).Encode(groups)
		rw.WriteHeader(http.StatusOK)
	}
}

func getCreateReportGroupHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		result, _ := sqliteDatasource.createReportGroup()

		json.NewEncoder(rw).Encode(result)
		rw.WriteHeader(http.StatusOK)
	}
}

func getUpdateReportGroupHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		var group ReportGroup
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &group)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.updateReportGroup(id, group)

		json.NewEncoder(rw).Encode(group)
		rw.WriteHeader(http.StatusOK)
	}
}

func getDeleteReportGroupHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		// TODO: Handle error
		success, _ := sqliteDatasource.deleteReportGroup(id)

		json.NewEncoder(rw).Encode(success)
		rw.WriteHeader(http.StatusOK)
	}
}

func getFetchReportGroupMembershipHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["group-id"]

		var assignment []ReportGroupMembership

		assignment = sqliteDatasource.getReportGroupMemberships(id)

		json.NewEncoder(rw).Encode(assignment)
		rw.WriteHeader(http.StatusOK)
	}
}

func getCreateReportGroupMembershipHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var assignment []ReportGroupMembership
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &assignment)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		result, _ := sqliteDatasource.createReportGroupMembership(assignment)

		json.NewEncoder(rw).Encode(result)
		rw.WriteHeader(http.StatusOK)
	}
}

func getDeleteReportGroupMembershipHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		// TODO: Handle error
		success, _ := sqliteDatasource.deleteReportGroupMembership(id)

		json.NewEncoder(rw).Encode(success)
		rw.WriteHeader(http.StatusOK)
	}
}

func getFetchReportContentHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		scheduleID := vars["schedule-id"]

		result, _ := sqliteDatasource.getReportContent(scheduleID)

		json.NewEncoder(rw).Encode(result)
		rw.WriteHeader(http.StatusOK)
	}
}

func getDeleteReportContentHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		// TODO: Handle error
		success, _ := sqliteDatasource.deleteReportContent(id)

		json.NewEncoder(rw).Encode(success)
		rw.WriteHeader(http.StatusOK)
	}
}

func getCreateReportContentHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		var reportContent ReportContent
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &reportContent)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		result, _ := sqliteDatasource.createReportContent(reportContent)

		json.NewEncoder(rw).Encode(result)
		rw.WriteHeader(http.StatusOK)
	}
}

func getUpdateReportContentHandler(sqliteDatasource *SQLiteDatasource) func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		var group ReportContent
		requestBody, err := request.GetBody()
		bodyAsBytes, _ := ioutil.ReadAll(requestBody)
		err = json.Unmarshal(bodyAsBytes, &group)

		if err != nil {
			http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
			return
		}

		sqliteDatasource.updateReportContent(id, group)

		json.NewEncoder(rw).Encode(group)
		rw.WriteHeader(http.StatusOK)
	}
}
