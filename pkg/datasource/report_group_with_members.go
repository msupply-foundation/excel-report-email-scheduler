package datasource

import (
	"database/sql"
	"excel-report-email-scheduler/pkg/api"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportGroupWithMembership struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Members     []api.MemberDetail `json:"members"`
}

type ReportGroupWithMembersRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

func (datasource *MsupplyEresDatasource) CreateReportGroupWithMembers(reportGroupWithMembers ReportGroupWithMembersRequest) (*ReportGroupWithMembersRequest, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
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

	reportGroupWithMembers.ID = uuid.New().String()

	_, err = stmt.Exec(reportGroupWithMembers.ID, reportGroupWithMembers.Name, reportGroupWithMembers.Description)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroup: stmt.Exec(): ", err.Error())
		return nil, err
	}

	var reportGroupMemberships []ReportGroupMembership
	for _, member := range reportGroupWithMembers.Members {
		newUuid := uuid.New().String()
		reportGroupMember := ReportGroupMembership{ID: newUuid, ReportGroupID: reportGroupWithMembers.ID, UserID: member}
		reportGroupMemberships = append(reportGroupMemberships, reportGroupMember)
	}

	_, err = datasource.CreateReportGroupMembership(reportGroupMemberships)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroupWithMembers: datasource.CreateReportGroupMembership: ", err.Error())
		return nil, err
	}

	return &reportGroupWithMembers, nil
}

func (datasource *MsupplyEresDatasource) DeleteReportGroupsWithMembers(id string) error {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: sql.Open(): ", err.Error())
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM ReportGroup WHERE id = ?")
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: db.Prepare(): ", err.Error())
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	stmt, err = db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroup: db.Prepare(): ", err.Error())
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return nil
}
