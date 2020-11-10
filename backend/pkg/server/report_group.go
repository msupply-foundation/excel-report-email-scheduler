package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

func (server *HttpServer) fetchReportGroup(rw http.ResponseWriter, request *http.Request) {
	var groups []dbstore.ReportGroup

	groups = server.db.GetReportGroups()

	json.NewEncoder(rw).Encode(groups)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createReportGroup(rw http.ResponseWriter, request *http.Request) {
	result, _ := server.db.CreateReportGroup()

	json.NewEncoder(rw).Encode(result)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) updateReportGroup(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var group dbstore.ReportGroup
	requestBody, err := request.GetBody()
	bodyAsBytes, _ := ioutil.ReadAll(requestBody)
	err = json.Unmarshal(bodyAsBytes, &group)

	if err != nil {
		http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
		return
	}

	server.db.UpdateReportGroup(id, group)

	json.NewEncoder(rw).Encode(group)
	rw.WriteHeader(http.StatusOK)

}

func (server *HttpServer) deleteReportGroup(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	// TODO: Handle error
	success, _ := server.db.DeleteReportGroup(id)

	json.NewEncoder(rw).Encode(success)
	rw.WriteHeader(http.StatusOK)

}
