package dbstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Schedule struct {
	ID             string `json:"id"`
	Interval       int    `json:"interval"`
	NextReportTime int    `json:"nextReportTime"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Lookback       int    `json:"lookback"`
	ReportGroupID  string `json:"reportGroupID"`
}

func NewSchedule(ID string, interval int, nextReportTime int, name string, description string, lookback int, reportGroupID string) Schedule {
	schedule := Schedule{ID: ID, Interval: interval, NextReportTime: nextReportTime, Name: name, Description: description, Lookback: lookback, ReportGroupID: reportGroupID}
	return schedule
}

func (datasource *SQLiteDatasource) OverdueSchedules() ([]Schedule, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Schedule WHERE strftime(\"%s\", \"now\") > nextReportTime")

	var schedules []Schedule
	for rows.Next() {
		var ID, Name, Description, ReportGroupID string
		var Interval, NextReportTime, Lookback int
		rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID)
		schedules = append(schedules, NewSchedule(ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID))
	}

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	return schedules, nil
}

func (datasource *SQLiteDatasource) CreateSchedule() (Schedule, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	newUuid := uuid.New().String()
	schedule := Schedule{ID: newUuid, NextReportTime: 0, Interval: 0, Name: "", Description: "", Lookback: 0, ReportGroupID: ""}
	stmt, _ := db.Prepare("INSERT INTO Schedule (ID,  nextReportTime, interval, name, description, lookback, reportGroupID) VALUES (?,?,?,?,?,?,?)")
	stmt.Exec(newUuid, 0, 60*1000*60*24, "New report schedule", "", 0, "")
	defer stmt.Close()

	return schedule, nil
}

func (datasource *SQLiteDatasource) UpdateSchedule(id string, schedule Schedule) *Schedule {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE Schedule SET nextReportTime = ?, interval = ?, name = ?, description = ?, lookback = ?, reportGroupID = ? where id = ?")
	stmt.Exec(schedule.NextReportTime, schedule.Interval, schedule.Name, schedule.Description, schedule.Lookback, schedule.ReportGroupID, id)
	defer stmt.Close()

	return nil
}

func (datasource *SQLiteDatasource) DeleteSchedule(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM Schedule WHERE id = ?")
	stmt.Exec(id)
	stmt, _ = db.Prepare("DELETE FROM ReportContent WHERE scheduleID = ?")
	stmt.Exec(id)

	defer stmt.Close()

	return true, nil
}

func (datasource *SQLiteDatasource) GetSchedules() []Schedule {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var schedules []Schedule

	rows, _ := db.Query("SELECT * FROM Schedule")
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description, ReportGroupID string
		var Interval, NextReportTime, Lookback int

		rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID)
		schedule := Schedule{ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID}

		schedules = append(schedules, schedule)
	}

	return schedules
}
