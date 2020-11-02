package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	dbstore "github.com/grafana/simple-datasource-backend/pkg/db"
)

func (server *HttpServer) createSchedule(rw http.ResponseWriter, request *http.Request) {
	schedule, _ := server.db.CreateSchedule()
	json.NewEncoder(rw).Encode(schedule)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) updateSchedule(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var schedule dbstore.Schedule
	requestBody, err := request.GetBody()
	bodyAsBytes, _ := ioutil.ReadAll(requestBody)
	err = json.Unmarshal(bodyAsBytes, &schedule)

	if err != nil {
		http.Error(rw, "Invalid Request, received: "+" "+err.Error()+"\n"+"Expecting a JSON body in the shape { grafanaUsername, grafanaPassword, email, emailPassword }", http.StatusBadRequest)
		return
	}

	server.db.UpdateSchedule(id, schedule)
}

func (server *HttpServer) fetchSchedules(rw http.ResponseWriter, request *http.Request) {
	var schedules []dbstore.Schedule

	schedules = server.db.GetSchedules()

	json.NewEncoder(rw).Encode(schedules)
	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteSchedule(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	// TODO: Handle error
	success, _ := server.db.DeleteSchedule(id)

	json.NewEncoder(rw).Encode(success)
	rw.WriteHeader(http.StatusOK)
}
