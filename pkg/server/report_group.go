package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"excel-report-email-scheduler/pkg/dbstore"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportGroupWithMembers struct {
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	Description string                          `json:"description"`
	Members     []dbstore.ReportGroupMembership `json:"members"`
}

func (server *HttpServer) fetchReportGroup(rw http.ResponseWriter, request *http.Request) {
	var groups []dbstore.ReportGroup

	groups, err := server.db.GetReportGroups()
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroup: db.GetReportGroups")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	var reportGroupsWithMembers []ReportGroupWithMembers

	for _, group := range groups {
		groupMembers, err := server.db.GetReportGroupMemberships(group.ID)
		if err != nil {
			log.DefaultLogger.Error("fetchReportGroup: db.GetReportGroups", err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		reportGroupWithMembers := ReportGroupWithMembers{ID: group.ID, Name: group.Name, Description: group.Description, Members: groupMembers}

		reportGroupsWithMembers = append(reportGroupsWithMembers, reportGroupWithMembers)
	}

	err = json.NewEncoder(rw).Encode(reportGroupsWithMembers)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroup: json.NewEncoder().Encode()", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) fetchSingleReportGroup(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	reportGroup, err := server.db.GetReportGroup(id)
	if err != nil {
		log.DefaultLogger.Error("fetchSingleReportGroup: server.db.GetReportGroup(id): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(reportGroup)
	if err != nil {
		log.DefaultLogger.Error("fetchSingleReportGroup: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) CreateReportGroupWithMembers(rw http.ResponseWriter, request *http.Request) {
	var group dbstore.ReportGroupWithMembers
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.Unmarshal(bodyAsBytes, &group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.ReportGroupFields()).Error(), http.StatusBadRequest)
		panic(err)
	}

	result, err := server.db.CreateReportGroupWithMembers(group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: db.UpdateReportGroup: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.DefaultLogger.Error("createReportGroup: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) updateReportGroup(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var group dbstore.ReportGroup
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.Unmarshal(bodyAsBytes, &group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.ReportGroupFields()).Error(), http.StatusBadRequest)
		panic(err)
	}

	_, err = server.db.UpdateReportGroup(id, group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: db.UpdateReportGroup: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)

}

func (server *HttpServer) deleteReportGroup(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteReportGroup(id)
	if err != nil {
		log.DefaultLogger.Error("deleteReportGroup: db.DeleteReportGroup(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)

}
