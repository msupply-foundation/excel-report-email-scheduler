package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

func (server *HttpServer) fetchReportContent(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	scheduleID := vars["schedule-id"]

	result, err := server.db.GetReportContent(scheduleID)
	if err != nil {
		log.DefaultLogger.Error("fetchReportContent: db.GetReportContent(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.DefaultLogger.Error("fetchReportContent: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createReportContent(rw http.ResponseWriter, request *http.Request) {
	var reportContent dbstore.ReportContent

	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("createReportContent: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("createReportContent: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyAsBytes, &reportContent)
	if err != nil {
		log.DefaultLogger.Error("createReportContent: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.ReportContentFields()).Error(), http.StatusBadRequest)
		return
	}

	result, err := server.db.CreateReportContent(reportContent)
	if err != nil {
		log.DefaultLogger.Error("createReportContent: db.CreateReportContent: " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.DefaultLogger.Error("createReportContent: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteReportContent(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteReportContent(id)
	if err != nil {
		log.DefaultLogger.Error("deleteReportContent: db.DeleteReportContent(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

func (server *HttpServer) updateReportContent(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var group dbstore.ReportContent
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateReportContent: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("updateReportContent: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyAsBytes, &group)
	if err != nil {
		log.DefaultLogger.Error("updateReportContent: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.ReportContentFields()).Error(), http.StatusBadRequest)
		return
	}

	reportContent, err := server.db.UpdateReportContent(id, group)
	if err != nil {
		log.DefaultLogger.Error("updateReportContent: db.UpdateReportContent: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(rw).Encode(reportContent)
	if err != nil {
		log.DefaultLogger.Error("updateReportContent: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
