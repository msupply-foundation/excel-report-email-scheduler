package server

import (
	"encoding/json"
	"net/http"

	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/emailer"
	"excel-report-email-scheduler/pkg/reportEmailer"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (server *HttpServer) testEmail(rw http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["schedule-id"]

	emailConfig, err := auth.NewEmailConfig(server.db)
	if err != nil {
		log.DefaultLogger.Error("testEmail: auth.NewEmailConfig: ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	} else {
		log.DefaultLogger.Debug("testEmail: emailConfig:", emailConfig)
	}

	authConfig, err := auth.NewAuthConfig(server.db)
	if err != nil {
		log.DefaultLogger.Error("testEmail: auth.NewAuthConfig: ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	} else {
		log.DefaultLogger.Debug("testEmail: authConfig:", authConfig)
	}

	settings, err := server.db.GetSettings()
	if err != nil {
		log.DefaultLogger.Error("testEmail: db.GetSettings: ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	} else {
		log.DefaultLogger.Debug("testEmail: db.getSettings:", settings)
	}

	schedule, err := server.db.GetSchedule(id)
	if err != nil {
		log.DefaultLogger.Error("testData: server.db.GetSchedule: ", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	} else {
		log.DefaultLogger.Debug("testEmail: schedule:", schedule)
	}

	em := emailer.New(emailConfig)
	re := reportEmailer.NewReportEmailer(server.db)
	re.CreateReport(*schedule, authConfig, settings.DatasourceID, *em)

	err = json.NewEncoder(rw).Encode("success")
	if err != nil {
		log.DefaultLogger.Error("testEmail: Email successfully sent")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}
