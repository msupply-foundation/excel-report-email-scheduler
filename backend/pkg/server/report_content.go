package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

func (server *HttpServer) fetchReportContent(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	scheduleID := vars["schedule-id"]

	result, _ := server.db.GetReportContent(scheduleID)

	json.NewEncoder(rw).Encode(result)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteReportContent(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	success, _ := server.db.DeleteReportContent(id)

	json.NewEncoder(rw).Encode(success)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createReportContent(rw http.ResponseWriter, request *http.Request) {
	var reportContent dbstore.ReportContent
	requestBody, err := request.GetBody()
	bodyAsBytes, _ := ioutil.ReadAll(requestBody)
	err = json.Unmarshal(bodyAsBytes, &reportContent)

	if err != nil {
		http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
		return
	}

	result, _ := server.db.CreateReportContent(reportContent)

	json.NewEncoder(rw).Encode(result)
	rw.WriteHeader(http.StatusOK)

}

func (server *HttpServer) updateReportContent(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var group dbstore.ReportContent
	requestBody, err := request.GetBody()
	bodyAsBytes, _ := ioutil.ReadAll(requestBody)
	err = json.Unmarshal(bodyAsBytes, &group)

	if err != nil {
		http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
		return
	}

	server.db.UpdateReportContent(id, group)

	json.NewEncoder(rw).Encode(group)
	rw.WriteHeader(http.StatusOK)
}
