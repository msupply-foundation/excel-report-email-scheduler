package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/api"
	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/setting"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (server *HttpServer) fetchReportGroupsWithMembers(rw http.ResponseWriter, request *http.Request) {
	settings, err := setting.NewSettings(request.Context())
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: db.GetReportGroups")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	authConfig, err := auth.NewAuthConfig(settings)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: auth.NewAuthConfig", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	var groups []datasource.ReportGroup
	groups, err = server.db.GetReportGroups()
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: server.db.GetReportGroups", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	var reportGroupsWithMembership []datasource.ReportGroupWithMembership

	for _, group := range groups {
		groupMemberUserIDs, err := server.db.GroupMemberUserIDs(group)
		if err != nil {
			log.DefaultLogger.Error("fetchReportGroupsWithMembers: server.db.GroupMemberUserIDs", err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		memberDetails, err := api.GetMemberDeatailsFromUserIDs(authConfig, groupMemberUserIDs, settings.DatasourceID)
		if err != nil {
			log.DefaultLogger.Error("fetchReportGroupsWithMembers: api.GetMemberDeatailsFromUserIDs", err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		reportGroupWithMembership := datasource.ReportGroupWithMembership{ID: group.ID, Name: group.Name, Description: group.Description, Members: memberDetails}

		reportGroupsWithMembership = append(reportGroupsWithMembership, reportGroupWithMembership)
	}

	err = json.NewEncoder(rw).Encode(reportGroupsWithMembership)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: json.NewEncoder().Encode()", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) deleteReportGroupsWithMembers(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteReportGroupsWithMembers(id)
	if err != nil {
		log.DefaultLogger.Error("deleteReportGroup: db.DeleteReportGroup(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
}
