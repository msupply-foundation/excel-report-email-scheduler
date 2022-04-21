package datasource

import (
	"database/sql"
	"errors"
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
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		datasource.logger.Error("GetReportGroups: sql.Open", err.Error())
		return nil, err
	}
	defer db.Close()

	var reportGroups []ReportGroup

	rows, err := db.Query("SELECT * FROM ReportGroup")
	if err != nil {
		datasource.logger.Error("GetReportGroups: db.Query(): ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description string
		err = rows.Scan(&ID, &Name, &Description)
		if err != nil {
			datasource.logger.Error("GetReportGroups: rows.Scan(): ", err.Error())
			return nil, err
		}

		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	return reportGroups, nil
}

func (datasource *MsupplyEresDatasource) GetSingleReportGroup(ID string) (*ReportGroup, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		datasource.logger.Error("GetSchedule: sql.Open(): ", err.Error())
		return nil, err
	}
	defer db.Close()

	var reportGroups []ReportGroup

	rows, err := db.Query("SELECT * FROM ReportGroup where id=?", ID)
	if err != nil {
		datasource.logger.Error("GetSchedules: db.Query(): ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description string

		err = rows.Scan(&ID, &Name, &Description)
		if err != nil {
			datasource.logger.Error("GetSchedules: rows.Scan(): ", err.Error())
			return nil, err
		}

		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	if len(reportGroups) > 0 {
		return &reportGroups[0], nil
	} else {
		return nil, errors.New("no report group found")
	}
}
