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

func (re *ReportEmailer) configs() (*auth.AuthConfig, *auth.EmailConfig, int, error) {
	authConfig, err := auth.NewAuthConfig(re.sql)
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.configs: NewAuthConfig: " + err.Error())
		return nil, nil, 0, err
	}

	emailConfig, err := auth.NewEmailConfig(re.sql)
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.configs: NewEmailConfig: " + err.Error())
		return nil, nil, 0, err
	}

	settings, err := re.sql.GetSettings()
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.configs: GetSettings: " + err.Error())
		return nil, nil, 0, err
	}

	return authConfig, emailConfig, settings.DatasourceID, nil
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

func (re *ReportEmailer) createReports()  {
	re.inProgress = true

	authConfig, emailConfig, datasourceID, err := re.configs()
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.createReports: re.configs: " + err.Error())
		return
	}

	em := emailer.New(emailConfig)

	schedules, err := re.sql.OverdueSchedules()
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.createReports: OverdueSchedules: " + err.Error())
		return
	}

	emails := make(map[string][]string)
	panels := make(map[string][]api.TablePanel)

	for _, schedule := range schedules {
		reportGroup, err := re.sql.ReportGroupFromSchedule(schedule)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReports: ReportGroupFromSchedule: " + err.Error())
			return
		}

		userIDs, err := re.sql.GroupMemberUserIDs(*reportGroup)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReports: GroupMemberUserIDs: " + err.Error())
			return
		}

		emailsFromUsers, err := api.GetEmails(*authConfig, userIDs, datasourceID)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReports: emailsFromUsers: " + err.Error())
			return
		}
		emails[schedule.ID] = emailsFromUsers

		reportContent, err := re.sql.GetReportContent(schedule.ID)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReports: GetReportContent: " + err.Error())
			return
		}

		panels[schedule.ID] = []api.TablePanel{}

		for _, content := range reportContent {
			interval := int64(schedule.Interval)
			to := strconv.FormatInt(time.Now().Unix(), 10)
			from := strconv.FormatInt(time.Now().Unix()-interval, 10)

			dashboard, err := api.NewDashboard(authConfig, content.DashboardID, from, to, datasourceID)
			if err != nil {
				log.DefaultLogger.Error("ReportEmailer.createReports: NewDashboard: " + err.Error())
				return
			}

			panel := dashboard.Panel(content.PanelID)
			panel.PrepSql(dashboard.Variables, content.StoreID)
			panels[schedule.ID] = append(panels[schedule.ID], *panel)
		}
	}

	templatePath := filepath.Join("data", "template.xlsx")
	reporter := reporter.NewReporter(templatePath)

	for scheduleID, reportSheetPanels := range panels {
		report := reporter.CreateNewReport(scheduleID)
		report.SetSheets(reportSheetPanels)
		err := report.Write(*authConfig)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReports: report.Write: " + err.Error())
			return
		}
	}

	for scheduleID, recipientEmails := range emails {
		attachmentPath := filepath.Join("data", scheduleID+".xlsx")
		em.BulkCreateAndSend(attachmentPath, recipientEmails)
	}

	re.cleanup(schedules)

	
}
