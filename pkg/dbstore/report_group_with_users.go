package dbstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportGroupWithMembers struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"selectedUsers"`
}

func NewReportGroupWithUsers(ID string, name string, description string, members []string) *ReportGroupWithMembers {
	return &ReportGroupWithMembers{ID: ID, Name: name, Description: description, Members: members}
}

func (datasource *SQLiteDatasource) CreateReportGroupWithMembers(reportGroupWithMembers ReportGroupWithMembers) (*ReportGroupWithMembers, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroupWithMembers: sql.Open(): ", err.Error())
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO ReportGroup (id, name, description) VALUES (?,?,?) RETURNING *")
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroupWithMembers: db.Prepare(): ", err.Error())
		return nil, err
	}
	defer stmt.Close()

	reportGroupID := uuid.New().String()

	_, err = stmt.Exec(reportGroupID, reportGroupWithMembers.Name, reportGroupWithMembers.Description)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroup: stmt.Exec(): ", err.Error())
		return nil, err
	}

	var reportGroupMemberships []ReportGroupMembership
	for _, member := range reportGroupWithMembers.Members {
		newUuid := uuid.New().String()
		reportGroupMember := ReportGroupMembership{ID: newUuid, ReportGroupID: reportGroupID, UserID: member}
		reportGroupMemberships = append(reportGroupMemberships, reportGroupMember)
	}

	_, err = datasource.CreateReportGroupMembership(reportGroupMemberships)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroupWithMembers: datasource.CreateReportGroupMembership: ", err.Error())
		return nil, err
	}

	return &reportGroupWithMembers, nil
}
