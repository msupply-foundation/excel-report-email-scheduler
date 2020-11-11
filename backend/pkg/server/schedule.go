package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

func (server *HttpServer) fetchSchedules(rw http.ResponseWriter, request *http.Request) {
	var schedules []dbstore.Schedule

	schedules, err := server.db.GetSchedules()
	if err != nil {
		log.DefaultLogger.Error("fetchSchedules: db.GetSchedules(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(schedules)
	if err != nil {
		log.DefaultLogger.Error("fetchSchedules: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createSchedule(rw http.ResponseWriter, request *http.Request) {
	schedule, err := server.db.CreateSchedule()
	if err != nil {
		log.DefaultLogger.Error("createSchedule: db.CreateSchedule(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(rw).Encode(schedule)
	if err != nil {
		log.DefaultLogger.Error("createSchedule: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteSchedule(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteSchedule(id)
	if err != nil {
		log.DefaultLogger.Error("deleteSchedule: db.DeleteSchedule(): " + id + " : " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) updateSchedule(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var schedule dbstore.Schedule
	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateSchedule: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("updateSchedule: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyAsBytes, &schedule)
	if err != nil {
		log.DefaultLogger.Error("updateSchedule: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, dbstore.ScheduleFields()).Error(), http.StatusBadRequest)
		return
	}

	_, err = server.db.UpdateSchedule(id, schedule)

	if err != nil {
		log.DefaultLogger.Error("updateSchedule: db.UpdateSchedule: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(rw).Encode(schedule)
	if err != nil {
		log.DefaultLogger.Error("updateSchedule: json.NewEncoder().Encode(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
