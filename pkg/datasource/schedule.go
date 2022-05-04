package datasource

import (
	"database/sql"
	"errors"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func NewSchedule(ID string, interval int, nextReportTime int, name string, description string, lookback int, reportGroupID string, time string, day int) Schedule {
	schedule := Schedule{
		ID:             ID,
		Interval:       interval,
		NextReportTime: nextReportTime,
		Name:           name,
		Description:    description,
		Lookback:       lookback,
		ReportGroupID:  reportGroupID,
		Time:           time,
		Day:            day,
		PanelDetails:   []ReportContent{},
	}
	return schedule
}

func (datasource *MsupplyEresDatasource) OverdueSchedules() ([]Schedule, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
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

func (datasource *MsupplyEresDatasource) GetSchedule(id string) (*Schedule, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("GetSchedule: sql.Open(): ", err.Error())
		return nil, err
	}

	var schedules []Schedule

	rows, err := db.Query("SELECT * FROM Schedule where id=?", id)
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

		schedule := Schedule{
			ID:             id,
			Interval:       Interval,
			NextReportTime: NextReportTime,
			Name:           Name,
			Description:    Description,
			Lookback:       Lookback,
			ReportGroupID:  ReportGroupID,
			Time:           Time,
			Day:            Day,
			PanelDetails:   []ReportContent{},
		}
		schedules = append(schedules, schedule)
	}

	if len(schedules) > 0 {
		return &schedules[0], nil
	} else {
		return nil, errors.New("no schedule found")
	}

}

func (datasource *MsupplyEresDatasource) UpdateSchedule(id string, schedule Schedule) (*Schedule, error) {
	db, err := sql.Open("sqlite", datasource.DataPath)
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
