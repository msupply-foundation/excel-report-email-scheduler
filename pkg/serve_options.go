package main

import (
	"excel-report-email-scheduler/pkg/dbstore"
	"excel-report-email-scheduler/pkg/server"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// Returns the options to be passed to the `Serve` API which starts the
// gRPC server and attaches these handlers as listeners
func getServeOptions() (datasource.ServeOpts, *dbstore.SQLiteDatasource, error) {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	log.DefaultLogger.Info("getServeOptions")

	sqlDatasource := dbstore.GetDataSource()

	server := server.NewServer(sqlDatasource)

	return datasource.ServeOpts{
		QueryDataHandler:    sqlDatasource,
		CheckHealthHandler:  sqlDatasource,
		CallResourceHandler: server.ResourceHandler(sqlDatasource),
	}, sqlDatasource, nil
}
