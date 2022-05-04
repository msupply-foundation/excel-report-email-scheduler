package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/datasource"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

func (server *HttpServer) fetchSchedules(rw http.ResponseWriter, request *http.Request) {
	frame := trace()
	var schedules []datasource.Schedule

	schedules, err := server.db.GetSchedules()
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = json.NewEncoder(rw).Encode(schedules)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) fetchSingleSchedule(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	frame := trace()

	schedule, err := server.db.GetSchedule(id)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = json.NewEncoder(rw).Encode(schedule)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) createSchedule(rw http.ResponseWriter, request *http.Request) {
	frame := trace()
	var schedule datasource.Schedule

	requestBody, err := request.GetBody()
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = json.Unmarshal(bodyAsBytes, &schedule)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = server.validator.ScheduleDuplicates(schedule)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = server.validator.ScheduleMustHavePanes(schedule)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = server.validator.ScheduleMustHaveReportGroup(schedule)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	_, err = server.db.CreateScheduleWithDetails(schedule)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	successMessageChunk := ""
	if schedule.ID != "" {
		successMessageChunk = "updated"
	} else {
		successMessageChunk = "created"
	}

	server.Success(rw, "Schedule successfully "+successMessageChunk)
}

func (server *HttpServer) deleteSchedule(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteSchedule(id)
	if err != nil {
		log.DefaultLogger.Error("deleteSchedule: db.DeleteSchedule(): " + id + " : " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}
