package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/auth"
	reportEmailer "excel-report-email-scheduler/pkg/report-emailer"
	"excel-report-email-scheduler/pkg/setting"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

type ExportPanelArgs struct {
	DashboardID string `json:"dashboardID"`
	PanelID     int    `json:"panelID"`
	Query       string `json:"query"`
	Title       string `json:"title"`
}

func ExportPanelArgsFields() string {
	return "\n{\n\tPanelID int\n\t" +
		"DashboardID string\n\t" +
		"Query string\n\t" +
		"Title string\n\t" +
		"\n}"
}

func (server *HttpServer) exportPanel(rw http.ResponseWriter, request *http.Request) {
	frame := trace()
	var args ExportPanelArgs

	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("exportPanel: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("exportPanel: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = json.Unmarshal(bodyAsBytes, &args)
	if err != nil {
		server.Error(rw, errors.Wrap(NewRequestBodyError(err, ExportPanelArgsFields()), frame.Function))
		return
	}

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

	reportEmailer.NewReportEmailer(server.db)
	templatePath := reportEmailer.GetFilePath("template")
	reporter := reportEmailer.NewReporter(templatePath)

	url, err := reporter.ExportPanel(authConfig, settings.DatasourceID, args.DashboardID, args.PanelID, args.Query, args.Title)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	fmt.Fprint(rw, url)
	rw.WriteHeader(http.StatusOK)
}
