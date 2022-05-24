package server

import (
	"excel-report-email-scheduler/pkg/auth"
	reportEmailer "excel-report-email-scheduler/pkg/report-emailer"
	"excel-report-email-scheduler/pkg/setting"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

func (server *HttpServer) testEmail(rw http.ResponseWriter, request *http.Request) {
	frame := trace()
	vars := mux.Vars(request)
	id := vars["schedule-id"]

	settings, err := setting.NewSettings(request.Context())
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	authConfig, err := auth.NewAuthConfig(settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	emailConfig, err := auth.NewEmailConfig(settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	schedule, err := server.db.GetSchedule(id)
	if err != nil {
		log.DefaultLogger.Error("testData: server.db.GetSchedule: ", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	} else {
		log.DefaultLogger.Debug("testEmail: schedule:", schedule)
	}

	em := reportEmailer.NewEmailSender(emailConfig)

	re := reportEmailer.NewReportEmailer(server.db)

	re.CreateReport(*schedule, authConfig, settings.DatasourceID, *em)

	server.Success(rw, "test report successfully sent")
}
