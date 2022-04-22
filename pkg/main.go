package main

import (
	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/server"

	"github.com/bugsnag/bugsnag-go"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

const ERES_PLUGIN_ID = "msupplyfoundation-excelreportemailscheduler-app"

func main() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "90618551d4e8d52a45260c61033093df",
		ReleaseStage: "production",
	})

	backend.SetupPluginEnvironment(ERES_PLUGIN_ID)

	pluginLogger := log.New()

	ds, server, err := Init(pluginLogger)
	if err != nil {
		pluginLogger.Error("Error starting mSupply Excel report e-mail scheduler datasource", "error", err.Error())
	}

	pluginLogger.Debug("Starting mSupply Excel report e-mail scheduler datasource")
	err = backend.Serve(backend.ServeOpts{
		CallResourceHandler: server.ResourceHandler(ds),
		QueryDataHandler:    ds,
		CheckHealthHandler:  ds,
	})
	if err != nil {
		pluginLogger.Error("Error starting mSupply Excel report e-mail scheduler datasource", "error", err.Error())
	}
}

func Init(logger log.Logger) (*datasource.MsupplyEresDatasource, *server.HttpServer, error) {
	mSupplyEresDatasource := datasource.NewMsupplyEresDatasource()

	server := server.NewServer(mSupplyEresDatasource)

	return mSupplyEresDatasource, server, nil

}
