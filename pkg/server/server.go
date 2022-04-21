package server

import (
	"excel-report-email-scheduler/pkg/datasource"

	"github.com/bugsnag/bugsnag-go"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

type HttpServer struct {
	db *datasource.MsupplyEresDatasource
}

func NewServer(sqliteDatasource *datasource.MsupplyEresDatasource) *HttpServer {
	return &HttpServer{db: sqliteDatasource}
}

func (server *HttpServer) ResourceHandler(mSupplyEresDatasource *datasource.MsupplyEresDatasource) backend.CallResourceHandler {
	mux := mux.NewRouter()

	mux.HandleFunc("/report-group", bugsnag.HandlerFunc(server.fetchReportGroupsWithMembers)).Methods("GET")
	mux.HandleFunc("/report-group/{id}", bugsnag.HandlerFunc(server.fetchSingleReportGroupWithMembers)).Methods("GET")
	mux.HandleFunc("/report-group", bugsnag.HandlerFunc(server.CreateReportGroupWithMembers)).Methods("POST")
	mux.HandleFunc("/report-group/{id}", bugsnag.HandlerFunc(server.deleteReportGroupsWithMembers)).Methods("DELETE")

	return httpadapter.New(mux)
}
