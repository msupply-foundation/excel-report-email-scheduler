package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"excel-report-email-scheduler/pkg/dbstore"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (server *HttpServer) fetchReportGroupMembership(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["group-id"]

	var assignment []dbstore.ReportGroupMembership

	assignment, err := server.db.GetReportGroupMemberships(id)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupMembership: db.GetReportGroupMemberships():" + id + ": " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(assignment)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupMembership: json.NewEncoder().Encode()")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createReportGroupMembership(rw http.ResponseWriter, request *http.Request) {
	var membership []dbstore.ReportGroupMembership
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("createReportGroupMembership: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("createReportGroupMembership: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.Unmarshal(bodyAsBytes, &membership)

	if err != nil {
		log.DefaultLogger.Error("createReportGroupMembership: json.Unmarshal: ", err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.ReportGroupMembershipFields()).Error(), http.StatusBadRequest)
		panic(err)
	}

	result, err := server.db.CreateReportGroupMembership(membership)
	if err != nil {
		log.DefaultLogger.Error("createReportGroupMembership: db.CreateReportGroupMembership: ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.DefaultLogger.Error("createReportGroupMembership: json.NewEncoder().Encode()", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteReportGroupMembership(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteReportGroupMembership(id)
	if err != nil {
		log.DefaultLogger.Error("deleteReportGroupMembership: db.DeleteReportGroupMembership(): ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)

}
