package validation

import (
	"database/sql"
	"errors"

	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/ereserror"
)

type Validation struct {
	datasource *datasource.MsupplyEresDatasource
	sqlClient  *datasource.SqlClient
}

func New(datasource *datasource.MsupplyEresDatasource) (*Validation, error) {
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, err, "Could not open database")
		return nil, err
	}

	return &Validation{datasource: datasource, sqlClient: sqlClient}, nil
}

func (validator *Validation) ReportGroupDuplicates(reportGroupWithMembers datasource.ReportGroupWithMembersRequest) error {
	var id string
	row := validator.sqlClient.Db.QueryRow("SELECT id FROM ReportGroup WHERE name = $1 LIMIT 1", reportGroupWithMembers.Name)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return nil
	default:
		if reportGroupWithMembers.ID == "" || (reportGroupWithMembers.ID != id) {
			err = ereserror.New(500, err, "Cannot have more than one report groups with same name")
			return err
		}
	}
	return nil
}

func (validator *Validation) ReportGroupMustHaveMembers(reportGroupWithMembers datasource.ReportGroupWithMembersRequest) error {
	memberLength := len(reportGroupWithMembers.Members)
	if memberLength <= 0 {
		err := errors.New("report group must have at least one member")
		err = ereserror.New(500, err, err.Error())
		return err
	}

	return nil
}

func (validator *Validation) GroupMemberUserIDsMustHaveElements(groupMemberUserIDs []string) error {
	memberLength := len(groupMemberUserIDs)

	if memberLength <= 0 {
		err := errors.New("report group must have members")
		err = ereserror.New(500, err, err.Error())
		return err
	}

	return nil
}
