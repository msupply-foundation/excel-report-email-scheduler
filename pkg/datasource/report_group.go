package datasource

import (
	"database/sql"
	"excel-report-email-scheduler/pkg/ereserror"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

type ReportGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ReportGroupFields() string {
	return "\n{\n\tid string" +
		"\n\tname string" +
		"\n\tdescription string\n}"
}

func NewReportGroup(ID string, name string, description string) *ReportGroup {
	return &ReportGroup{ID: ID, Name: name, Description: description}
}

func (datasource *MsupplyEresDatasource) GetReportGroups() ([]ReportGroup, error) {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.Db.Close()

	var reportGroups []ReportGroup

	rows, err := sqlClient.Db.Query("SELECT * FROM ReportGroup")
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not get report group list")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description string
		err = rows.Scan(&ID, &Name, &Description)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not scan report group rows")
			return nil, err
		}

		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	return reportGroups, nil
}

func (datasource *MsupplyEresDatasource) GetSingleReportGroup(ID string) (*ReportGroup, error) {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.Db.Close()

	var reportGroups []ReportGroup

	rows, err := sqlClient.Db.Query("SELECT * FROM ReportGroup where id=?", ID)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function),
			fmt.Sprintf("Could not find ReportGroup with id: %s", ID))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description string

		err = rows.Scan(&ID, &Name, &Description)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), fmt.Sprintf("Could not scan rows of ReportGroup with id: %s", ID))
			return nil, err
		}

		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	if len(reportGroups) > 0 {
		return &reportGroups[0], nil
	} else {
		return nil, ereserror.New(500, errors.Wrap(err, frame.Function), fmt.Sprintf("Could not find of ReportGroup with id: %s", ID))
	}
}

func (datasource *MsupplyEresDatasource) ReportGroupFromSchedule(schedule Schedule) (*ReportGroup, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("ReportGroupFromSchedule: sql.Open", err.Error())
		return nil, err
	}

	row := db.QueryRow("SELECT * FROM ReportGroup WHERE ID = ?", schedule.ReportGroupID)

	var ID, name, description string
	err = row.Scan(&ID, &name, &description)
	if err != nil {
		log.DefaultLogger.Error("ReportGroupFromSchedule: rows.Scan(): ", err.Error())
		return nil, err
	}

	return NewReportGroup(ID, name, description), nil
}
