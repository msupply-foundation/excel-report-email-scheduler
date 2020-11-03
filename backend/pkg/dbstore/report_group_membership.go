package dbstore

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

func NewReportGroupMembership(ID string, userID string, reportGroupID string) *ReportGroupMembership {
	return &ReportGroupMembership{ID: ID, UserID: userID, ReportGroupID: reportGroupID}
}

func (datasource *SQLiteDatasource) GroupMemberUserIDs(reportGroup ReportGroup) ([]string, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ReportGroupMembership WHERE reportGroupID = ?", reportGroup.ID)

	if err != nil {
		log.DefaultLogger.Error("GroupMemberUserIDs", err.Error())
		return nil, err
	}

	var memberships []ReportGroupMembership
	for rows.Next() {
		var ID, UserID, ReportGroupID string
		rows.Scan(&ID, &UserID, &ReportGroupID)
		membership := ReportGroupMembership{ID, UserID, ReportGroupID}
		memberships = append(memberships, membership)
	}

	var userIDs []string
	for _, member := range memberships {
		userIDs = append(userIDs, member.UserID)
	}

	return userIDs, nil
}

func (datasource *SQLiteDatasource) GetReportGroupMemberships(groupID string) []ReportGroupMembership {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var memberships []ReportGroupMembership

	rows, _ := db.Query("SELECT * FROM ReportGroupMembership WHERE reportGroupID = ?", groupID)
	defer rows.Close()

	for rows.Next() {
		var ID, UserID, ReportGroupID string
		rows.Scan(&ID, &UserID, &ReportGroupID)
		membership := ReportGroupMembership{ID, UserID, ReportGroupID}
		memberships = append(memberships, membership)
	}

	return memberships
}

func (datasource *SQLiteDatasource) CreateReportGroupMembership(members []ReportGroupMembership) ([]ReportGroupMembership, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var addedMemberships []ReportGroupMembership
	for _, member := range members {
		newUuid := uuid.New().String()
		stmt, _ := db.Prepare("INSERT INTO ReportGroupMembership (ID, userID, reportGroupID) VALUES (?,?,?)")
		stmt.Exec(newUuid, member.UserID, member.ReportGroupID)
		member.ID = newUuid
		addedMemberships = append(addedMemberships, member)
		defer stmt.Close()
	}

	// TODO: Return report assignment
	return addedMemberships, nil
}

func (datasource *SQLiteDatasource) DeleteReportGroupMembership(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportGroupMembership WHERE id = ?")
	stmt.Exec(id)
	defer stmt.Close()

	// TODO: Proper return values, returning error or false? or just an error, probably
	return true, nil
}
