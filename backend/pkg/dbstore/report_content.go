package dbstore

import (
	"database/sql"

	"github.com/google/uuid"
)

type ReportContent struct {
	ID          string `json:"id"`
	ScheduleID  string `json:"scheduleID"`
	PanelID     int    `json:"panelID"`
	DashboardID string `json:"dashboardID"`
	Lookback    int    `json:"lookback"`
	StoreID     string `json:"storeID"`
}

func (datasource *SQLiteDatasource) GetReportContent(scheduleID string) ([]ReportContent, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var reportContent []ReportContent

	rows, _ := db.Query("SELECT * FROM ReportContent WHERE scheduleID = ?", scheduleID)
	defer rows.Close()

	for rows.Next() {
		var ID, ScheduleID, StoreID, DashboardID string
		var Lookback, PanelID int
		rows.Scan(&ID, &ScheduleID, &PanelID, &DashboardID, &Lookback, &StoreID)
		content := ReportContent{ID, ScheduleID, PanelID, DashboardID, Lookback, StoreID}
		reportContent = append(reportContent, content)
	}

	return reportContent, nil
}

func (datasource *SQLiteDatasource) DeleteReportContent(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportContent WHERE ID = ?")
	stmt.Exec(id)
	defer stmt.Close()

	return true, nil
}

func (datasource *SQLiteDatasource) CreateReportContent(newReportContentValues ReportContent) (ReportContent, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	reportContent := ReportContent{ID: uuid.New().String(), ScheduleID: newReportContentValues.ScheduleID, PanelID: newReportContentValues.PanelID, DashboardID: newReportContentValues.DashboardID, Lookback: 0, StoreID: ""}

	stmt, _ := db.Prepare("INSERT INTO ReportContent (id, scheduleID, panelID, dashboardID, lookback, storeID) VALUES (?,?,?,?,?,?)")

	stmt.Exec(reportContent.ID, reportContent.ScheduleID, reportContent.PanelID, reportContent.DashboardID, reportContent.Lookback, reportContent.StoreID)

	defer stmt.Close()

	return reportContent, nil
}

func (datasource *SQLiteDatasource) UpdateReportContent(id string, reportGroup ReportContent) ReportContent {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE ReportContent SET scheduleID = ?, panelID = ?, storeID = ?, lookback = ? where id = ?")
	stmt.Exec(reportGroup.ScheduleID, reportGroup.PanelID, reportGroup.StoreID, reportGroup.Lookback, id)
	defer stmt.Close()

	return reportGroup
}
