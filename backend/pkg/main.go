package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/robfig/cron"
)

// Main entry point of the backend plugin.
// Using the provided SDK, it is required to simply make a call to
// datasource.Serve() with the options required (as below)
// {
//  QueryDataHandler: A datasource with a QueryData method which handles /tsdb/query calls.
//  CallResourceHandler: An HTTP handler. Create with sdk-go/backend/resource/httpadapter
//  CheckHealthHandler: A datasource with a CheckHealth method which handles /api/{pluginID}/health endpoint
//  TransformDataHandler: Transformer [experimental]
//  GRPCSettings: Settings..?
// }
func main() {
	log.DefaultLogger.Info("Starting up")
	serveOptions, sql, err := getServeOptions()
	if err != nil {
		log.DefaultLogger.Error("FATAL. Error when retreiving serveOptions", err.Error())
		panic(err)
	}

	re := NewReportEmailer(sql)

	// Try to send reports on loading
	re.createReports()

	// Set up scheduler which will try to send reports every 10 minutes
	c := cron.New()
	c.AddFunc("@every 10m", re.createReports)
	c.Start()

	// Start listening to requests sent from Grafana. This call is blocking and
	// and waits until Grafana shutsdown or the plugin exits.
	err = datasource.Serve(serveOptions)

	if err != nil {
		log.DefaultLogger.Error("FATAL. Plugin datasource.Serve() returned an error: " + err.Error())
		os.Exit(1)
	}
}
