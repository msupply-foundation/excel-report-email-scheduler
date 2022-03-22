package dbstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ReportContent struct {
	ID          string `json:"id"`
	ScheduleID  string `json:"scheduleID"`
	PanelID     int    `json:"panelID"`
	DashboardID string `json:"dashboardID"`
	Lookback    int    `json:"lookback"`
	Variables   string `json:"variables"`
}

func ReportContentFields() string {
	return "\n{\n\tID string\n\t" +
		"ScheduleID string\n\t" +
		"PanelID string\n\t" +
		"DashboardID string\n\t" +
		"Lookback int\n\t" +
		"Variables string" +
		"\n}"
}

func (datasource *SQLiteDatasource) GetReportContent(scheduleID string) ([]ReportContent, error) {
	var reportContent []ReportContent

	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("GetReportContent: sql.Open: ", err.Error())
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM ReportContent WHERE scheduleID = ?", scheduleID)
	defer rows.Close()
	if err != nil {
		log.DefaultLogger.Error("GetReportContent: db.Query()", err.Error())
		return nil, err
	}

	for rows.Next() {
		var ID, ScheduleID, DashboardID, Variables string
		var Lookback, PanelID int
		err = rows.Scan(&ID, &ScheduleID, &PanelID, &DashboardID, &Lookback, &Variables)
		if err != nil {
			log.DefaultLogger.Error("GetReportContent: rows.Scan() ", err.Error())
			return nil, err
		}

		content := ReportContent{ID, ScheduleID, PanelID, DashboardID, Lookback, Variables}
		reportContent = append(reportContent, content)
	}

	return reportContent, nil
}

func (datasource *SQLiteDatasource) DeleteReportContent(id string) error {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("DeleteReportContent: sql.Open", err.Error())
		return err
	}

	stmt, err := db.Prepare("DELETE FROM ReportContent WHERE ID = ?")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("db.Prepare: ", err.Error())
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.DefaultLogger.Error("DeleteReportContent - stmt.Exec :" + err.Error())
		return err
	}

	return nil
}

func (datasource *SQLiteDatasource) CreateReportContent(newReportContentValues ReportContent) (*ReportContent, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateReportContent: sql.Open", err.Error())
		return nil, err
	}

	reportContent := ReportContent{ID: uuid.New().String(), ScheduleID: newReportContentValues.ScheduleID, PanelID: newReportContentValues.PanelID, DashboardID: newReportContentValues.DashboardID, Lookback: 0, Variables: ""}

	stmt, err := db.Prepare("INSERT INTO ReportContent (id, scheduleID, panelID, dashboardID, lookback, variables) VALUES (?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateReportContent: ", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(reportContent.ID, reportContent.ScheduleID, reportContent.PanelID, reportContent.DashboardID, reportContent.Lookback, reportContent.Variables)
	if err != nil {
		log.DefaultLogger.Error("CreateReportContent: ", err.Error())
		return nil, err
	}

	return &reportContent, nil
}

func (datasource *SQLiteDatasource) UpdateReportContent(id string, reportContent ReportContent) (*ReportContent, error) {
	db, err := sql.Open("sqlite", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("UpdateReportContent: sql.Open: ", err.Error())
		return nil, err
	}

	stmt, err := db.Prepare("UPDATE ReportContent SET scheduleID = ?, panelID = ?, lookback = ?, variables = ? where id = ?")
	defer stmt.Close()
	if err != nil {
		log.DefaultLogger.Error("UpdateReportContent: db.Prepare: ", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(reportContent.ScheduleID, reportContent.PanelID, reportContent.Lookback, reportContent.Variables, id)
	if err != nil {
		log.DefaultLogger.Error("UpdateReportContent: db.Exec: ", err.Error())
		return nil, err
	}

	return &reportContent, nil
}
