package validation

import (
	"database/sql"

	"excel-report-email-scheduler/pkg/datasource"
	"excel-report-email-scheduler/pkg/ereserror"

	"github.com/pkg/errors"
)

func (validator *Validation) ScheduleDuplicates(schedule datasource.Schedule) error {
	frame := trace()
	var id string
	row := validator.sqlClient.Db.QueryRow("SELECT id FROM Schedule WHERE name = $1 LIMIT 1", schedule.Name)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return nil
	default:
		if schedule.ID == "" || (schedule.ID != id) {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Cannot have more than one schedule with same name")
			return err
		}
	}

	return nil
}

func (validator *Validation) ScheduleMustHaveReportGroup(schedule datasource.Schedule) error {
	frame := trace()
	if schedule.ReportGroupID == "" {
		err := errors.New("schedule must have a report group selected")
		err = ereserror.New(500, errors.Wrap(err, frame.Function), err.Error())
		return err
	}

	var id string
	row := validator.sqlClient.Db.QueryRow("SELECT id FROM ReportGroup WHERE id = $1 LIMIT 1", schedule.ReportGroupID)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Must select valid report group")
		return err
	}

	return nil
}

func (validator *Validation) ScheduleMustHavePanes(schedule datasource.Schedule) error {
	frame := trace()
	paneLength := len(schedule.PanelDetails)
	if paneLength <= 0 {
		err := errors.New("schedule must have at least one panel selected")
		err = ereserror.New(500, errors.Wrap(err, frame.Function), err.Error())
		return err
	}

	return nil
}
