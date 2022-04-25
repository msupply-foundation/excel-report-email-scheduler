package datasource

import (
	"excel-report-email-scheduler/pkg/api"
	"excel-report-email-scheduler/pkg/ereserror"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.db.Close()

	if reportGroupWithMembers.ID == "" {
		stmt, err := sqlClient.db.Prepare("INSERT INTO ReportGroup (id, name, description) VALUES (?,?,?) RETURNING *")
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not create report group record")
			return nil, err
		}
		defer stmt.Close()

		reportGroupWithMembers.ID = uuid.New().String()

		_, err = stmt.Exec(reportGroupWithMembers.ID, reportGroupWithMembers.Name, reportGroupWithMembers.Description)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not create report group record")
			return nil, err
		}
	} else {
		stmt, err := sqlClient.db.Prepare("UPDATE ReportGroup SET name = ?, description = ? where id = ? RETURNING *")
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not update report group record")
			return nil, err
		}

		_, err = stmt.Exec(reportGroupWithMembers.Name, reportGroupWithMembers.Description, reportGroupWithMembers.ID)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not update report group record")
			return nil, err
		}

		err = datasource.DeleteReportGroupMembersByGroupID(reportGroupWithMembers.ID)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not update report group record")
			return nil, err
		}
	}

	var reportGroupMemberships []ReportGroupMembership
	for _, member := range reportGroupWithMembers.Members {
		newUuid := uuid.New().String()
		reportGroupMember := ReportGroupMembership{ID: newUuid, ReportGroupID: reportGroupWithMembers.ID, UserID: member}
		reportGroupMemberships = append(reportGroupMemberships, reportGroupMember)
	}

	_, err = datasource.CreateReportGroupMembership(reportGroupMemberships)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not update report group record")
		return nil, err
	}

	return &reportGroupWithMembers, nil
}

func (datasource *MsupplyEresDatasource) DeleteReportGroupsWithMembers(id string) error {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return err
	}
	defer sqlClient.db.Close()

	stmt, err := sqlClient.db.Prepare("DELETE FROM ReportGroup WHERE id = ?")
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not delete report group record")
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	stmt, err = sqlClient.db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not delete report group record")
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return nil
}
