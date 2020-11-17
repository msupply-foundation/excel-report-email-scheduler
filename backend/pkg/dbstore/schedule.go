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

func ScheduleFields() string {
	return "\n{\n\tID string" +
		"\n\tinterval int" +
		"\n\tnextReportTime int\n}" +
		"\n\tname string\n}" +
		"\n\tdescription string\n}" +
		"\n\tlookback int\n}" +
		"\n\treportGroupID string\n}"
}

func NewSchedule(ID string, interval int, nextReportTime int, name string, description string, lookback int, reportGroupID string) Schedule {
	schedule := Schedule{ID: ID, Interval: interval, NextReportTime: nextReportTime, Name: name, Description: description, Lookback: lookback, ReportGroupID: reportGroupID}
	return schedule
}

func (datasource *SQLiteDatasource) OverdueSchedules() ([]Schedule, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("OverdueSchedules: sql.Open", err.Error())
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM Schedule WHERE strftime(\"%s\", \"now\") > nextReportTime")
	if err != nil {
		log.DefaultLogger.Error("OverdueSchedules: db.Query", err.Error())
		return nil, err
	}

	var schedules []Schedule
	for rows.Next() {
		var ID, Name, Description, ReportGroupID string
		var Interval, NextReportTime, Lookback int
		err = rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID)
		if err != nil {
			log.DefaultLogger.Error("OverdueSchedules: sql.Open", err.Error())
			return nil, err
		}
		schedules = append(schedules, NewSchedule(ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID))
	}

	return schedules, nil
}

func (datasource *SQLiteDatasource) CreateSchedule() (*Schedule, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateSchedule: sql.Open", err.Error())
		return nil, err
	}

	newUuid := uuid.New().String()
	schedule := Schedule{ID: newUuid, NextReportTime: 0, Interval: 0, Name: "", Description: "", Lookback: 0, ReportGroupID: ""}
	stmt, err := db.Prepare("INSERT INTO Schedule (ID,  nextReportTime, interval, name, description, lookback, reportGroupID) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.DefaultLogger.Error("CreateSchedule: db.Prepare()", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(newUuid, 0, 60*60*24, "New report schedule", "", 0, "")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateSchedule: stmt.Exec()", err.Error())
		return nil, err
	}

	return &schedule, nil
}

func (datasource *SQLiteDatasource) UpdateSchedule(id string, schedule Schedule) (*Schedule, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("UpdateSchedule: sql.Open()", err.Error())
		return nil, err
	}

	stmt, err := db.Prepare("UPDATE Schedule SET nextReportTime = ?, interval = ?, name = ?, description = ?, lookback = ?, reportGroupID = ? where id = ?")
	if err != nil {
		log.DefaultLogger.Error("UpdateSchedule: db.Prepare()", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(schedule.NextReportTime, schedule.Interval, schedule.Name, schedule.Description, schedule.Lookback, schedule.ReportGroupID, id)
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("UpdateSchedule: stmt.Exec()", err.Error())
		return nil, err
	}

	return &schedule, nil
}

func (datasource *SQLiteDatasource) DeleteSchedule(id string) error {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("DeleteSchedule: sql.Open()", err.Error())
		return err
	}

	stmt, err := db.Prepare("DELETE FROM Schedule WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("DeleteSchedule: db.Prepare()1", err.Error())
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.DefaultLogger.Error("DeleteSchedule: stmt.Exec()1", err.Error())
		return err
	}

	stmt, err = db.Prepare("DELETE FROM ReportContent WHERE scheduleID = ?")
	if err != nil {
		log.DefaultLogger.Error("DeleteSchedule: db.Prepare()2", err.Error())
		return err
	}

	stmt.Exec(id)
	if err != nil {
		log.DefaultLogger.Error("DeleteSchedule: stmt.Exec()2", err.Error())
		return err
	}

	return nil
}

func (datasource *SQLiteDatasource) GetSchedules() ([]Schedule, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("GetSchedules: sql.Open(): ", err.Error())
		return nil, err
	}

	var schedules []Schedule

	rows, err := db.Query("SELECT * FROM Schedule")
	defer rows.Close()
	if err != nil {
		log.DefaultLogger.Error("GetSchedules: db.Query(): ", err.Error())
		return nil, err
	}

	for rows.Next() {
		var ID, Name, Description, ReportGroupID string
		var Interval, NextReportTime, Lookback int

		err = rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID)
		if err != nil {
			log.DefaultLogger.Error("GetSchedules: rows.Scan(): ", err.Error())
			return nil, err
		}

		schedule := Schedule{ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
