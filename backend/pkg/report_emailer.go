package main

import (
	"os"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/api"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
	"github.com/grafana/simple-datasource-backend/pkg/emailer"
	"github.com/grafana/simple-datasource-backend/pkg/reporter"
)

type ReportEmailer struct {
	sql        *dbstore.SQLiteDatasource
	inProgress bool
}

func NewReportEmailer(sql *dbstore.SQLiteDatasource) *ReportEmailer {
	return &ReportEmailer{sql: sql, inProgress: false}
}

func (re *ReportEmailer) configs() (*auth.AuthConfig, *auth.EmailConfig) {
	authConfig := auth.NewAuthConfig(re.sql)
	emailConfig := auth.NewEmailConfig(re.sql)

	return authConfig, emailConfig
}

func (re *ReportEmailer) cleanup(schedules []dbstore.Schedule) {
	for _, schedule := range schedules {
		os.Remove("./data/" + schedule.ID + ".xlsx")

		// Unix is in seconds. Story in Milliseconds
		schedule.NextReportTime = int(time.Now().Unix()*1000) + schedule.Interval
		re.sql.UpdateSchedule(schedule.ID, schedule)
	}

	re.inProgress = false
}

func (re *ReportEmailer) createReports() {
	re.inProgress = true

	authConfig, emailConfig := re.configs()

	em := emailer.New(emailConfig)

	schedules, _ := re.sql.OverdueSchedules()

	emails := make(map[string][]string)
	panels := make(map[string][]api.TablePanel)

	for _, schedule := range schedules {
		reportGroup := re.sql.ReportGroupFromSchedule(schedule)
		userIDs, _ := re.sql.GroupMemberUserIDs(*reportGroup)
		emails[schedule.ID] = api.GetEmails(*authConfig, userIDs)

		reportContent, _ := re.sql.GetReportContent(schedule.ID)
		panels[schedule.ID] = []api.TablePanel{}

		for _, content := range reportContent {
			dashboard, _ := api.NewDashboard(authConfig, content.DashboardID)
			panel, _ := dashboard.Panel(content.PanelID)
			panels[schedule.ID] = append(panels[schedule.ID], *panel)
		}
	}

	reporter := reporter.NewReporter("./data/template.xlsx")

	for scheduleID, reportSheetPanels := range panels {
		report := reporter.CreateNewReport(scheduleID)
		report.SetSheets(reportSheetPanels)
		report.Write(*authConfig)
	}

	for scheduleID, recipientEmails := range emails {
		log.DefaultLogger.Info("Sending Emails!!")
		em.BulkCreateAndSend("./data/"+scheduleID+".xlsx", recipientEmails)
	}

	re.cleanup(schedules)

}
