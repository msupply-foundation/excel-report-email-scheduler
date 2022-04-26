package validation

import (
	"database/sql"

	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/ereserror"

	"github.com/pkg/errors"
)

func (validator *Validation) ReportGroupDuplicates(reportGroupWithMembers datasource.ReportGroupWithMembersRequest) error {
	frame := trace()
	var id string
	row := validator.sqlClient.Db.QueryRow("SELECT id FROM ReportGroup WHERE name = $1 LIMIT 1", reportGroupWithMembers.Name)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return nil
	default:
		if reportGroupWithMembers.ID == "" || (reportGroupWithMembers.ID != id) {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Cannot have more than one report groups with same name")
			return err
		}
	}
	return nil
}

func (validator *Validation) ReportGroupMustHaveMembers(reportGroupWithMembers datasource.ReportGroupWithMembersRequest) error {
	frame := trace()
	memberLength := len(reportGroupWithMembers.Members)
	if memberLength <= 0 {
		err := errors.New("report group must have at least one member")
		err = ereserror.New(500, errors.Wrap(err, frame.Function), err.Error())
		return err
	}

	return nil
}

func (validator *Validation) GroupMemberUserIDsMustHaveElements(groupMemberUserIDs []string) error {
	memberLength := len(groupMemberUserIDs)
	frame := trace()
	if memberLength <= 0 {
		err := errors.New("report group must have members")
		err = ereserror.New(500, errors.Wrap(err, frame.Function), err.Error())
		return err
	}

	return nil
}
