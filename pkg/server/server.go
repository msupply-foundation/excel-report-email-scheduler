package server

import (
	"net/http"

	"excel-report-email-scheduler/pkg/dbstore"

	"github.com/bugsnag/bugsnag-go"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

type HttpServer struct {
	db *dbstore.SQLiteDatasource
}

func NewServer(sqliteDatasource *dbstore.SQLiteDatasource) *HttpServer {
	return &HttpServer{db: sqliteDatasource}
}

func (server *HttpServer) ResourceHandler(sqliteDatasource *dbstore.SQLiteDatasource) backend.CallResourceHandler {

	mux := mux.NewRouter()

	mux.HandleFunc("/settings", bugsnag.HandlerFunc(server.updateSettings)).Methods("POST")
	mux.HandleFunc("/settings", bugsnag.HandlerFunc(server.fetchSettings)).Methods("GET")

	mux.HandleFunc("/schedule", bugsnag.HandlerFunc(server.createSchedule)).Methods("POST")
	mux.HandleFunc("/schedule/{id}", bugsnag.HandlerFunc(server.fetchSingleSchedules)).Methods("GET")
	mux.HandleFunc("/schedule/{id}", bugsnag.HandlerFunc(server.updateSchedule)).Methods("PUT")
	mux.HandleFunc("/schedule", bugsnag.HandlerFunc(server.fetchSchedules)).Methods("GET")
	mux.HandleFunc("/schedule/{id}", bugsnag.HandlerFunc(server.deleteSchedule)).Methods("DELETE")

	mux.HandleFunc("/report-group", bugsnag.HandlerFunc(server.fetchReportGroup)).Methods("GET")
	mux.HandleFunc("/report-group/{id}", bugsnag.HandlerFunc(server.updateReportGroup)).Methods("PUT")
	mux.HandleFunc("/report-group", bugsnag.HandlerFunc(server.createReportGroup)).Methods("POST")
	mux.HandleFunc("/report-group/{id}", bugsnag.HandlerFunc(server.deleteReportGroup)).Methods("DELETE")

	mux.HandleFunc("/report-group-membership", bugsnag.HandlerFunc(server.fetchReportGroupMembership)).Queries("group-id", "{group-id}").Methods("GET")
	mux.HandleFunc("/report-group-membership", bugsnag.HandlerFunc(server.createReportGroupMembership)).Methods("POST")
	mux.HandleFunc("/report-group-membership/{id}", bugsnag.HandlerFunc(server.deleteReportGroupMembership)).Methods("DELETE")

	mux.HandleFunc("/report-content", bugsnag.HandlerFunc(server.fetchReportContent)).Queries("schedule-id", "{schedule-id}").Methods("GET")
	mux.HandleFunc("/report-content", bugsnag.HandlerFunc(server.createReportContent)).Methods("POST")
	mux.HandleFunc("/report-content/{id}", bugsnag.HandlerFunc(server.updateReportContent)).Methods("PUT")
	mux.HandleFunc("/report-content/{id}", bugsnag.HandlerFunc(server.deleteReportContent)).Methods("DELETE")

	mux.HandleFunc("/test-email", bugsnag.HandlerFunc(server.testEmail)).Queries("schedule-id", "{schedule-id}").Methods("GET")
	mux.HandleFunc("/export-panel", bugsnag.HandlerFunc(server.exportPanel)).Methods("POST")
	mux.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("../data")))).Methods("GET")

	return httpadapter.New(mux)
}
