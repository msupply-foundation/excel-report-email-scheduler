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

	serveOptions, sqliteDatasource := getServeOptions()

	c := cron.New()
	c.AddFunc("@every 1m", getScheduler(sqliteDatasource))
	c.Start()

	// Start listening to requests sent from Grafana. This call is blocking and
	// and waits until Grafana shutsdown or the plugin exits.
	err := datasource.Serve(serveOptions)

	// TODO: Defer a Close() ?

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
