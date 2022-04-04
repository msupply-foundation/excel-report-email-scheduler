package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/reporter"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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
		log.DefaultLogger.Error("exportPanel: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, ExportPanelArgsFields()).Error(), http.StatusBadRequest)
		panic(err)
	}

	authConfig, err := auth.NewAuthConfig(server.db)
	if err != nil {
		log.DefaultLogger.Error("exportPanel: auth.NewAuthConfig: ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	settings, err := server.db.GetSettings()
	if err != nil {
		log.DefaultLogger.Error("exportPanel: db.GetSettings: ", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	templatePath := reporter.GetFilePath("template")
	reporter := reporter.NewReporter(templatePath)

	url, err := reporter.ExportPanel(authConfig, settings.DatasourceID, args.DashboardID, args.PanelID, args.Query, args.Title)
	if err != nil {
		log.DefaultLogger.Error("exportPanel: reporter.SaveReport: ", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	fmt.Fprint(rw, url)
	rw.WriteHeader(http.StatusOK)
}
