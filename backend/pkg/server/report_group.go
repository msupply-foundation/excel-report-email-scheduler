package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

func (server *HttpServer) fetchReportGroup(rw http.ResponseWriter, request *http.Request) {
	var groups []dbstore.ReportGroup

	groups, err := server.db.GetReportGroups()
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroup: db.GetReportGroups")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	err = json.NewEncoder(rw).Encode(groups)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroup: json.NewEncoder().Encode()", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createReportGroup(rw http.ResponseWriter, request *http.Request) {
	result, err := server.db.CreateReportGroup()
	if err != nil {
		log.DefaultLogger.Error("createReportGroup: CreateReportGroup(): " + err.Error())
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
