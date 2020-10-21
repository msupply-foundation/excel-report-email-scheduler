package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

// Returns the options to be passed to the `Serve` API which starts the
// gRPC server and attaches these handlers as listeners
func getServeOptions() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.

	sqlDatasource := getDataSource()
	httpHandler := getHttpHandler(sqlDatasource)

	return datasource.ServeOpts{
		QueryDataHandler:    sqlDatasource,
		CheckHealthHandler:  sqlDatasource,
		CallResourceHandler: httpHandler,
	}
}
