package datasource

import (
	"database/sql"
	"excel-report-email-scheduler/pkg/ereserror"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

func (datasource *MsupplyEresDatasource) GetReportContent(scheduleID string) ([]ReportContent, error) {
	var reportContent []ReportContent

	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("GetReportContent: sql.Open: ", err.Error())
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ReportContent WHERE scheduleID = ?", scheduleID)
	if err != nil {
		log.DefaultLogger.Error("GetReportContent: db.Query()", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID, ScheduleID, DashboardID, Variables, Lookback string
		var PanelID int
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

func (datasource *MsupplyEresDatasource) CreateReportContents(reportContents []ReportContent) (*[]ReportContent, error) {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.Db.Close()

	var addedReportContents []ReportContent

	for _, reportContent := range reportContents {
		newUuid := uuid.New().String()

		stmt, err := sqlClient.Db.Prepare("INSERT INTO ReportContent (id, scheduleID, panelID, dashboardID, lookback, variables) VALUES (?,?,?,?,?,?)")
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not create report content")
			return nil, err
		}
		defer stmt.Close()

		reportContent.ID = newUuid

		_, err = stmt.Exec(reportContent.ID, reportContent.ScheduleID, reportContent.PanelID, reportContent.DashboardID, reportContent.Lookback, reportContent.Variables)
		if err != nil {
			err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not create report content")
			return nil, err
		}

		addedReportContents = append(addedReportContents, reportContent)
	}

	return &addedReportContents, nil
}

func (datasource *MsupplyEresDatasource) DeleteReportContentByScheduleID(id string) error {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return err
	}
	defer sqlClient.Db.Close()

	stmt, err := sqlClient.Db.Prepare("DELETE FROM ReportContent WHERE scheduleID = ?")
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not delete report content")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not delete report content")
		return err
	}

	return nil
}
