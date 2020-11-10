package main

import (
	"path/filepath"
	"strconv"
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

func (re *ReportEmailer) configs() (*auth.AuthConfig, *auth.EmailConfig, int) {
	authConfig := auth.NewAuthConfig(re.sql)
	emailConfig := auth.NewEmailConfig(re.sql)
	settings := re.sql.GetSettings()

	return authConfig, emailConfig, settings.DatasourceID
}

func (re *ReportEmailer) cleanup(schedules []dbstore.Schedule) {
	// for _, schedule := range schedules {
	// 	path := filepath.Join("data", schedule.ID+".xlsx")
	// 	os.Remove(path)

	// 	schedule.NextReportTime = int(time.Now().Unix()) + schedule.Interval
	// 	re.sql.UpdateSchedule(schedule.ID, schedule)
	// }

	re.inProgress = false
}

func (re *ReportEmailer) createReports() {
	re.inProgress = true
	log.DefaultLogger.Info("Creating Reports")

	authConfig, emailConfig, datasourceID := re.configs()

	em := emailer.New(emailConfig)

	schedules, _ := re.sql.OverdueSchedules()

	log.DefaultLogger.Info("Found Schedules to create reports for:")

	for _, schedule := range schedules {
		log.DefaultLogger.Info("- " + schedule.Name)
	}

	emails := make(map[string][]string)
	panels := make(map[string][]api.TablePanel)

	for _, schedule := range schedules {

		log.DefaultLogger.Info("Collecting Emails for: " + schedule.Name)

		reportGroup := re.sql.ReportGroupFromSchedule(schedule)
		userIDs, _ := re.sql.GroupMemberUserIDs(*reportGroup)
		emails[schedule.ID] = api.GetEmails(*authConfig, userIDs, datasourceID)

		log.DefaultLogger.Info("Sending Emails to:")

		for k, v := range emails {
			log.DefaultLogger.Info("Emails for :" + k)
			for _, email := range v {
				log.DefaultLogger.Info("- " + email)
			}
		}

		reportContent, _ := re.sql.GetReportContent(schedule.ID)
		panels[schedule.ID] = []api.TablePanel{}

		for _, content := range reportContent {
			interval := int64(schedule.Interval)
			to := strconv.FormatInt(time.Now().Unix(), 10)
			from := strconv.FormatInt(time.Now().Unix()-interval, 10)

			dashboard, _ := api.NewDashboard(authConfig, content.DashboardID, from, to, datasourceID)
			panel, _ := dashboard.Panel(content.PanelID)
			panel.PrepSql(dashboard.Variables, content.StoreID)
			panels[schedule.ID] = append(panels[schedule.ID], *panel)
		}
	}

	templatePath := filepath.Join("data", "template.xlsx")
	reporter := reporter.NewReporter(templatePath)

	for scheduleID, reportSheetPanels := range panels {

		log.DefaultLogger.Info("Creating report for : " + scheduleID)
		report := reporter.CreateNewReport(scheduleID)
		report.SetSheets(reportSheetPanels)
		report.Write(*authConfig)
	}

	for scheduleID, recipientEmails := range emails {
		log.DefaultLogger.Info("Sending Emails!!")
		attachmentPath := filepath.Join("data", scheduleID+".xlsx")
		em.BulkCreateAndSend(attachmentPath, recipientEmails)
	}

	re.cleanup(schedules)

}
