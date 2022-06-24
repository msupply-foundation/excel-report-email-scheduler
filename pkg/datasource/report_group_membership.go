package datasource

import (
	"excel-report-email-scheduler/pkg/ereserror"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.Db.Close()

	rows, err := sqlClient.Db.Query("SELECT * FROM ReportGroupMembership WHERE reportGroupID = ?", reportGroup.ID)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function),
			fmt.Sprintf("Could not find Report Group members for report group with id: %s", reportGroup.ID))
		return nil, err
	}

	var memberships []ReportGroupMembership
	for rows.Next() {
		var ID, UserID, ReportGroupID string
		err = rows.Scan(&ID, &UserID, &ReportGroupID)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function),
				fmt.Sprintf("Could not find Report Group members for report group with id: %s", reportGroup.ID))
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
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.Db.Close()

	var addedMemberships []ReportGroupMembership
	for _, member := range members {
		newUuid := uuid.New().String()

		stmt, err := sqlClient.Db.Prepare("INSERT INTO ReportGroupMembership (ID, userID, reportGroupID) VALUES (?,?,?)")
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not create report group membership")
			return nil, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(newUuid, member.UserID, member.ReportGroupID)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not create report group membership")
			return nil, err
		}

		member.ID = newUuid
		addedMemberships = append(addedMemberships, member)
	}

	return &addedMemberships, nil
}

func (datasource *MsupplyEresDatasource) DeleteReportGroupMembersByGroupID(id string) error {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return err
	}
	defer sqlClient.Db.Close()

	stmt, err := sqlClient.Db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not delete report group membership")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not delete report group membership")
		return err
	}

	return nil
}
