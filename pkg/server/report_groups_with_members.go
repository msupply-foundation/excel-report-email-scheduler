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
	"github.com/pkg/errors"
)

func (server *HttpServer) fetchSingleReportGroupWithMembers(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	frame := trace()

	settings, err := setting.NewSettings(request.Context())
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	authConfig, err := auth.NewAuthConfig(settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	var group *datasource.ReportGroup
	group, err = server.db.GetSingleReportGroup(id)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	groupMemberUserIDs, err := server.db.GroupMemberUserIDs(group)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	memberDetails, err := api.GetMemberDeatailsFromUserIDs(authConfig, groupMemberUserIDs, settings.DatasourceID)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	reportGroupWithMembership := datasource.ReportGroupWithMembership{ID: group.ID, Name: group.Name, Description: group.Description, Members: memberDetails}

	err = json.NewEncoder(rw).Encode(reportGroupWithMembership)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) fetchReportGroupsWithMembers(rw http.ResponseWriter, request *http.Request) {
	frame := trace()

	settings, err := setting.NewSettings(request.Context())
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	authConfig, err := auth.NewAuthConfig(settings)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	var groups []datasource.ReportGroup
	groups, err = server.db.GetReportGroups()
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	reportGroupsWithMembership := []datasource.ReportGroupWithMembership{}

	if len(groups) > 0 {
		for _, group := range groups {
			groupMemberUserIDs, err := server.db.GroupMemberUserIDs(&group)
			if err != nil {
				server.Error(rw, errors.Wrap(err, frame.Function))
				return
			}

			err = server.validator.GroupMemberUserIDsMustHaveElements(groupMemberUserIDs)
			if err != nil {
				server.Error(rw, errors.Wrap(err, frame.Function))
				break
			}

			memberDetails, err := api.GetMemberDeatailsFromUserIDs(authConfig, groupMemberUserIDs, settings.DatasourceID)
			if err != nil {
				server.Error(rw, errors.Wrap(err, frame.Function))
				return
			}

			reportGroupWithMembership := datasource.ReportGroupWithMembership{ID: group.ID, Name: group.Name, Description: group.Description, Members: memberDetails}

			reportGroupsWithMembership = append(reportGroupsWithMembership, reportGroupWithMembership)
		}
	}

	err = json.NewEncoder(rw).Encode(reportGroupsWithMembership)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (server *HttpServer) CreateReportGroupWithMembers(rw http.ResponseWriter, request *http.Request) {
	frame := trace()
	var group datasource.ReportGroupWithMembersRequest

	requestBody, err := request.GetBody()
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	bodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = json.Unmarshal(bodyAsBytes, &group)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = server.validator.ReportGroupDuplicates(group)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	err = server.validator.ReportGroupMustHaveMembers(group)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	_, err = server.db.CreateReportGroupWithMembers(group)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
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
	frame := trace()
	vars := mux.Vars(request)
	id := vars["id"]

	err := server.db.DeleteReportGroupsWithMembers(id)
	if err != nil {
		server.Error(rw, errors.Wrap(err, frame.Function))
		return
	}

	server.Success(rw, "Report group deleted")
}
