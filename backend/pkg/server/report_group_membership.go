package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

func (server *HttpServer) fetchReportGroupMembership(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["group-id"]

	var assignment []dbstore.ReportGroupMembership

	assignment = server.db.GetReportGroupMemberships(id)

	json.NewEncoder(rw).Encode(assignment)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createReportGroupMembership(rw http.ResponseWriter, request *http.Request) {
	var membership []dbstore.ReportGroupMembership
	requestBody, err := request.GetBody()
	bodyAsBytes, _ := ioutil.ReadAll(requestBody)
	err = json.Unmarshal(bodyAsBytes, &membership)

	if err != nil {
		http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
		return
	}

	result, _ := server.db.CreateReportGroupMembership(membership)

	json.NewEncoder(rw).Encode(result)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteReportGroupMembership(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	// TODO: Handle error
	success, _ := server.db.DeleteReportGroupMembership(id)

	json.NewEncoder(rw).Encode(success)
	rw.WriteHeader(http.StatusOK)

}
