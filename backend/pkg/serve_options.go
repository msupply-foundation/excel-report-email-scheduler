package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	dbstore "github.com/grafana/simple-datasource-backend/pkg/db"
	"github.com/grafana/simple-datasource-backend/pkg/server"
)

// Returns the options to be passed to the `Serve` API which starts the
// gRPC server and attaches these handlers as listeners
func getServeOptions() (datasource.ServeOpts, *dbstore.SQLiteDatasource) {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.

	// TODO: Handle errors here
	sqlDatasource, _ := dbstore.GetDataSource()
	server := server.NewServer(sqlDatasource)

	return datasource.ServeOpts{
		QueryDataHandler:    sqlDatasource,
		CheckHealthHandler:  sqlDatasource,
		CallResourceHandler: server.ResourceHandler(sqlDatasource),
	}, sqlDatasource
}
