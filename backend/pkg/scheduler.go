package main

import "github.com/grafana/grafana-plugin-sdk-go/backend/log"

func getScheduler(sqlite *SQLiteDatasource) func() {
	return func() {
		log.DefaultLogger.Info("Scheduler!")
	}
}
