package server

import (
	"encoding/json"
	"runtime"

	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/ereserror"
	"excel-report-email-scheduler/pkg/validation"
	"net/http"

	"github.com/bugsnag/bugsnag-go"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/pkg/errors"
)

type HttpServer struct {
	db        *datasource.MsupplyEresDatasource
	validator *validation.Validation
}

func NewServer(sqliteDatasource *datasource.MsupplyEresDatasource) *HttpServer {
	validator, _ := validation.New(sqliteDatasource)

	return &HttpServer{db: sqliteDatasource, validator: validator}
}

func (server *HttpServer) ResourceHandler(mSupplyEresDatasource *datasource.MsupplyEresDatasource) backend.CallResourceHandler {
	mux := mux.NewRouter()

	mux.HandleFunc("/report-group", bugsnag.HandlerFunc(server.fetchReportGroupsWithMembers)).Methods("GET")
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

func (server *HttpServer) Error(rw http.ResponseWriter, err error) {
	log.DefaultLogger.Error(err.Error())

	var ew ereserror.EresError
	if errors.As(err, &ew) {
		ew = ew.Dig()
		http.Error(rw, ew.Message, ew.Code)
	} else {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func trace() *runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return &frame
}
