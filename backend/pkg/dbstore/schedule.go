package dbstore

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	Time           string `json:"time"`
	Day            int    `json:"day"`
}

func ScheduleFields() string {
	return "\n{\n\tID string" +
		"\n\tinterval int" +
		"\n\tnextReportTime int\n" +
		"\n\tname string\n" +
		"\n\tdescription string\n" +
		"\n\tlookback int\n" +
		"\n\treportGroupID string\n" +
		"\n\time string\n" +
		"\n\nday int\n}"
}

func NewSchedule(ID string, interval int, nextReportTime int, name string, description string, lookback int, reportGroupID string, time string, day int) Schedule {
	schedule := Schedule{ID: ID, Interval: interval, NextReportTime: nextReportTime, Name: name, Description: description, Lookback: lookback, ReportGroupID: reportGroupID, Time: time, Day: day}
	return schedule
}

func (schedule *Schedule) UpdateNextReportTime() {
	now := time.Now()
	daysOffset := 1
	if schedule.Day > 0 {
		daysOffset = schedule.Day
	}

	reportTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	// there is probably a better way to parse the time string, it didn't work for me though
	timeOfDay, err := time.Parse(time.RFC3339, "1970-01-01T"+schedule.Time+":00+00:00")
	if err == nil {
		log.DefaultLogger.Info(fmt.Sprintf("Adding time of '%s' to '%s'", schedule.Time, reportTime))
		reportTime = time.Date(reportTime.Year(), reportTime.Month(), reportTime.Day(), timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, now.Location())
	}

	// for the intervals using x Day of y, remove the current Day value
	daysOffset = (-1 * int(reportTime.Day())) + daysOffset

	switch schedule.Interval {
	case 5: //yearly
		reportTime = reportTime.AddDate(1, 0, daysOffset)
	case 4: // quarterly
		reportTime = reportTime.AddDate(0, 3, daysOffset)
	case 3: // monthly
		reportTime = reportTime.AddDate(0, 1, daysOffset)
	case 2: // fortnightly
		reportTime = reportTime.AddDate(0, 0, 14)
	case 1: // weekly
		reportTime = reportTime.AddDate(0, 0, 7)
	default: // 0 == daily
		if reportTime.Unix() < now.Unix() {
			// run tomorrow
			reportTime = reportTime.AddDate(0, 0, 1)
		}
	}
	schedule.NextReportTime = int(reportTime.Unix())
	log.DefaultLogger.Info(fmt.Sprintf("Setting time of schedule '%s' to '%s'", schedule.Name, reportTime))
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
		var ID, Name, Description, ReportGroupID, Time string
		var Day, Interval, NextReportTime, Lookback int
		err = rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID, &Time, &Day)
		if err != nil {
			log.DefaultLogger.Error("OverdueSchedules: sql.Open", err.Error())
			return nil, err
		}
		schedules = append(schedules, NewSchedule(ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID, Time, Day))
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
	schedule := Schedule{ID: newUuid, NextReportTime: 0, Interval: 0, Name: "", Description: "", Lookback: 0, ReportGroupID: "", Time: "", Day: 1}
	schedule.UpdateNextReportTime()
	stmt, err := db.Prepare("INSERT INTO Schedule (ID,  nextReportTime, interval, name, description, lookback, reportGroupID, time, day) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.DefaultLogger.Error("CreateSchedule: db.Prepare()", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(newUuid, 0, 60*60*24, "New report schedule", "", 0, "", "", 1)
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

	stmt, err := db.Prepare("UPDATE Schedule SET nextReportTime = ?, interval = ?, name = ?, description = ?, lookback = ?, reportGroupID = ?, time = ?, day = ? where id = ?")
	if err != nil {
		log.DefaultLogger.Error("UpdateSchedule: db.Prepare()", err.Error())
		return nil, err
	}

	schedule.UpdateNextReportTime()
	_, err = stmt.Exec(schedule.NextReportTime, schedule.Interval, schedule.Name, schedule.Description, schedule.Lookback, schedule.ReportGroupID, schedule.Time, schedule.Day, id)
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

func (datasource *SQLiteDatasource) GetSchedule(id string) (*Schedule, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("GetSchedule: sql.Open(): ", err.Error())
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
		var ID, Name, Description, ReportGroupID, Time string
		var Day, Interval, NextReportTime, Lookback int

		err = rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID, &Time, &Day)
		if err != nil {
			log.DefaultLogger.Error("GetSchedules: rows.Scan(): ", err.Error())
			return nil, err
		}

		schedule := Schedule{ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID, Time, Day}
		schedules = append(schedules, schedule)
	}

	if len(schedules) > 0 {
		return &schedules[0], nil
	} else {
		return nil, errors.New("No schedule found")
	}

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
		var ID, Name, Description, ReportGroupID, Time string
		var Day, Interval, NextReportTime, Lookback int

		err = rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID, &Time, &Day)
		if err != nil {
			log.DefaultLogger.Error("GetSchedules: rows.Scan(): ", err.Error())
			return nil, err
		}

		schedule := Schedule{ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID, Time, Day}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
