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
	StoreID     string `json:"storeID"`
}

func ReportContentFields() string {
	return "\n{\n\tID string\n\tScheduleID string\n\tPanelID string\n\tDashboardID string\n\tLookback int\n\tStoreID string\n}"
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

func (datasource *SQLiteDatasource) DeleteReportContent(id string) error {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM ReportContent WHERE ID = ?")

	if err != nil {
		log.DefaultLogger.Error("DeleteReportContent - db.Prepare :" + err.Error())
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		log.DefaultLogger.Error("DeleteReportContent - stmt.Exec :" + err.Error())
		return err
	}

	defer stmt.Close()

	return nil
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

func (datasource *SQLiteDatasource) UpdateReportContent(id string, reportContent ReportContent) (*ReportContent, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	if err != nil {
		log.DefaultLogger.Error("UpdateReportContent: sql.Open: ", err.Error())
		return nil, err
	}

	stmt, err := db.Prepare("UPDATE ReportContent SET scheduleID = ?, panelID = ?, storeID = ?, lookback = ? where id = ?")

	if err != nil {
		log.DefaultLogger.Error("UpdateReportContent: db.Prepare: ", err.Error())
		return nil, err
	}

	_, err = stmt.Exec(reportContent.ScheduleID, reportContent.PanelID, reportContent.StoreID, reportContent.Lookback, id)

	if err != nil {
		log.DefaultLogger.Error("UpdateReportContent: db.Exec: ", err.Error())
		return nil, err
	}

	defer stmt.Close()

	return &reportContent, nil
}
