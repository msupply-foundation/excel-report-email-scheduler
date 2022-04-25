package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/datasource"
	"net/http"

	"github.com/bugsnag/bugsnag-go"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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

	mux.HandleFunc("/report-group", server.fetchReportGroupsWithMembers).Methods("GET")
	mux.HandleFunc("/report-group/{id}", bugsnag.HandlerFunc(server.fetchSingleReportGroupWithMembers)).Methods("GET")
	mux.HandleFunc("/report-group", bugsnag.HandlerFunc(server.CreateReportGroupWithMembers)).Methods("POST")
	mux.HandleFunc("/report-group/{id}", bugsnag.HandlerFunc(server.deleteReportGroupsWithMembers)).Methods("DELETE")

	return httpadapter.New(mux)
}

func (server *HttpServer) Success(rw http.ResponseWriter, message string) {
	rw.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: db.UpdateReportGroup: " + err.Error())
	}
	rw.Write(jsonResp)
}

func (server *HttpServer) Error(rw http.ResponseWriter, status int, message string, err error) {
	data := make(map[string]interface{})

	data["status"] = status

	switch status {
	case 404:
		data["message"] = "Not Found"
	case 500:
		data["message"] = "Internal Server Error"
	}

	if message != "" {
		data["message"] = message
	}

	if err != nil {
		data["error"] = err.Error()
	}

	rw.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: db.UpdateReportGroup: " + err.Error())
	}
	rw.Write(jsonResp)
}
