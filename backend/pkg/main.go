package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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

	// Start listening to requests sent from Grafana. This call is blocking and
	// and waits until Grafana shutsdown or the plugin exits.
	serveOptions := getServeOptions()
	err := datasource.Serve(serveOptions)

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
