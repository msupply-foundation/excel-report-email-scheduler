package server

import (
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"

	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

type HttpServer struct {
	db *dbstore.SQLiteDatasource
}

func NewServer(sqliteDatasource *dbstore.SQLiteDatasource) *HttpServer {
	return &HttpServer{db: sqliteDatasource}
}

func (server *HttpServer) ResourceHandler(sqliteDatasource *dbstore.SQLiteDatasource) backend.CallResourceHandler {
	mux := mux.NewRouter()

	mux.HandleFunc("/settings", server.updateSettings).Methods("POST")
	mux.HandleFunc("/settings", server.fetchSettings).Methods("GET")

	mux.HandleFunc("/schedule", server.createSchedule).Methods("POST")
	mux.HandleFunc("/schedule/{id}", server.updateSchedule).Methods("PUT")
	mux.HandleFunc("/schedule", server.fetchSchedules).Methods("GET")
	mux.HandleFunc("/schedule/{id}", server.deleteSchedule).Methods("DELETE")

	mux.HandleFunc("/report-group", server.fetchReportGroup).Methods("GET")
	mux.HandleFunc("/report-group/{id}", server.updateReportGroup).Methods("PUT")
	mux.HandleFunc("/report-group", server.createReportGroup).Methods("POST")
	mux.HandleFunc("/report-group/{id}", server.deleteReportGroup).Methods("DELETE")

	mux.HandleFunc("/report-group-membership", server.fetchReportGroupMembership).Queries("group-id", "{group-id}").Methods("GET")
	mux.HandleFunc("/report-group-membership", server.createReportGroupMembership).Methods("POST")
	mux.HandleFunc("/report-group-membership/{id}", server.deleteReportGroupMembership).Methods("DELETE")

	mux.HandleFunc("/report-content", server.fetchReportContent).Queries("schedule-id", "{schedule-id}").Methods("GET")
	mux.HandleFunc("/report-content", server.createReportContent).Methods("POST")
	mux.HandleFunc("/report-content/{id}", server.updateReportContent).Methods("PUT")
	mux.HandleFunc("/report-content/{id}", server.deleteReportContent).Methods("DELETE")

	return httpadapter.New(mux)
}
