package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/datasource"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

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
