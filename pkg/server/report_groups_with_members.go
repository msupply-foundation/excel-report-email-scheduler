package server

import (
	"encoding/json"
	"excel-report-email-scheduler/pkg/api"
	"excel-report-email-scheduler/pkg/auth"
	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/setting"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (server *HttpServer) fetchSingleReportGroupWithMembers(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	settings, err := setting.NewSettings(request.Context())
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: db.GetReportGroups")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		// panic(err) is causing a RPC error so removing it for now
		return
	}

	authConfig, err := auth.NewAuthConfig(settings)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: auth.NewAuthConfig", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var group *datasource.ReportGroup
	group, err = server.db.GetSingleReportGroup(id)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: server.db.GetReportGroups", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	groupMemberUserIDs, err := server.db.GroupMemberUserIDs(group)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: server.db.GroupMemberUserIDs", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	memberDetails, err := api.GetMemberDeatailsFromUserIDs(authConfig, groupMemberUserIDs, settings.DatasourceID)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: api.GetMemberDeatailsFromUserIDs", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	reportGroupWithMembership := datasource.ReportGroupWithMembership{ID: group.ID, Name: group.Name, Description: group.Description, Members: memberDetails}

	err = json.NewEncoder(rw).Encode(reportGroupWithMembership)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: json.NewEncoder().Encode()", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) fetchReportGroupsWithMembers(rw http.ResponseWriter, request *http.Request) {
	settings, err := setting.NewSettings(request.Context())
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: db.GetReportGroups")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	authConfig, err := auth.NewAuthConfig(settings)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: auth.NewAuthConfig", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var groups []datasource.ReportGroup
	groups, err = server.db.GetReportGroups()
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: server.db.GetReportGroups", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	reportGroupsWithMembership := []datasource.ReportGroupWithMembership{}

	if len(groups) > 0 {
		for _, group := range groups {
			groupMemberUserIDs, err := server.db.GroupMemberUserIDs(&group)
			if err != nil {
				log.DefaultLogger.Error("fetchReportGroupsWithMembers: server.db.GroupMemberUserIDs", err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			memberDetails, err := api.GetMemberDeatailsFromUserIDs(authConfig, groupMemberUserIDs, settings.DatasourceID)
			if err != nil {
				log.DefaultLogger.Error("fetchReportGroupsWithMembers: api.GetMemberDeatailsFromUserIDs", err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			reportGroupWithMembership := datasource.ReportGroupWithMembership{ID: group.ID, Name: group.Name, Description: group.Description, Members: memberDetails}

			reportGroupsWithMembership = append(reportGroupsWithMembership, reportGroupWithMembership)
		}
	}

	err = json.NewEncoder(rw).Encode(reportGroupsWithMembership)
	if err != nil {
		log.DefaultLogger.Error("fetchReportGroupsWithMembers: json.NewEncoder().Encode()", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) CreateReportGroupWithMembers(rw http.ResponseWriter, request *http.Request) {
	var group datasource.ReportGroupWithMembersRequest

	requestBody, err := request.GetBody()
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: request.GetBody(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: ioutil.ReadAll(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyAsBytes, &group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: json.Unmarshal: " + err.Error())
		http.Error(rw, NewRequestBodyError(err, datasource.ReportGroupFields()).Error(), http.StatusBadRequest)
		return
	}

	_, err = server.db.CreateReportGroupWithMembers(group)
	if err != nil {
		log.DefaultLogger.Error("updateReportGroup: db.UpdateReportGroup: " + err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	successMessageChunk := ""
	if group.ID != "" {
		successMessageChunk = "updated"
	} else {
		successMessageChunk = "created"
	}

	server.Success(rw, "Report group successfully "+successMessageChunk)
}

func (server *HttpServer) deleteReportGroupsWithMembers(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteReportGroupsWithMembers(id)
	if err != nil {
		log.DefaultLogger.Error("deleteReportGroupsWithMembers: db.DeleteReportGroupsWithMembers(): " + err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	server.Success(rw, "Report group deleted")
}
