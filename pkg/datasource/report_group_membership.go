package datasource

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportGroupMembership struct {
	ID            string `json:"id"`
	UserID        string `json:"userID"`
	ReportGroupID string `json:"reportGroupID"`
}

func ReportGroupMembershipFields() string {
	return "\n{\n\tID string\n\tUserID string\nReportGroupID string\n}"
}

func NewReportGroupMembership(ID string, userID string, reportGroupID string) *ReportGroupMembership {
	return &ReportGroupMembership{ID: ID, UserID: userID, ReportGroupID: reportGroupID}
}

func (datasource *MsupplyEresDatasource) GroupMemberUserIDs(reportGroup *ReportGroup) ([]string, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("GroupMemberUserIDs: sql.Open", err.Error())
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ReportGroupMembership WHERE reportGroupID = ?", reportGroup.ID)
	if err != nil {
		log.DefaultLogger.Error("GroupMemberUserIDs: db.Query()", err.Error())
		return nil, err
	}

	var memberships []ReportGroupMembership
	for rows.Next() {
		var ID, UserID, ReportGroupID string
		err = rows.Scan(&ID, &UserID, &ReportGroupID)
		if err != nil {
			log.DefaultLogger.Error("GroupMemberUserIDs: rows.Scan(): ", err.Error())
			return nil, err
		}
		membership := ReportGroupMembership{ID, UserID, ReportGroupID}
		memberships = append(memberships, membership)
	}

	var userIDs []string
	for _, member := range memberships {
		userIDs = append(userIDs, member.UserID)
	}

	return userIDs, nil
}

func (datasource *MsupplyEresDatasource) CreateReportGroupMembership(members []ReportGroupMembership) (*[]ReportGroupMembership, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("CreateReportGroupMembership: sql.Open(): ", err.Error())
		return nil, err
	}
	defer db.Close()

	var addedMemberships []ReportGroupMembership
	for _, member := range members {
		newUuid := uuid.New().String()

		stmt, err := db.Prepare("INSERT INTO ReportGroupMembership (ID, userID, reportGroupID) VALUES (?,?,?)")
		if err != nil {
			log.DefaultLogger.Error("CreateReportGroupMembership: db.Prepare(): ", err.Error())
			return nil, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(newUuid, member.UserID, member.ReportGroupID)
		if err != nil {
			log.DefaultLogger.Error("CreateReportGroupMembership: stmt.Exec() ", err.Error())
			return nil, err
		}

		member.ID = newUuid
		addedMemberships = append(addedMemberships, member)
	}

	return &addedMemberships, nil
}

func (datasource *MsupplyEresDatasource) DeleteReportGroupMembersByGroupID(id string) error {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroupMembership: sql.Open(): ", err.Error())
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroupMembership: db.Prepare(): ", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.DefaultLogger.Error("DeleteReportGroupMembership: stmt.Exec(): ", err.Error())
		return err
	}

	return nil
}
