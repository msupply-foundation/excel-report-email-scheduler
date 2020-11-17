package reportEmailer

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bugsnag/bugsnag-go"

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
	log.DefaultLogger.Info("Starting Clean up...")
	for _, schedule := range schedules {

		fileName := schedule.Name + ".xlsx"
		log.DefaultLogger.Info(fmt.Sprintf("Deleting %s...", fileName))
		path := filepath.Join("..", "data", fileName)
		err := os.Remove(path)
		if err != nil {
			// Failure case shouldn't be much of a problem since we're using the schedule ID for the report name, at the moment
			// as it will just write to the same file and not create infinitely many if deleting always fails.
			log.DefaultLogger.Error(fmt.Sprintf("Could not delete %s... : %s", fileName, err.Error()))
		}

		schedule.NextReportTime = int(time.Now().Unix()) + schedule.Interval
		re.sql.UpdateSchedule(schedule.ID, schedule)
	}

	re.inProgress = false
}

func (re *ReportEmailer) CreateReport(schedule dbstore.Schedule, authConfig *auth.AuthConfig, datasourceID int, em emailer.Emailer) error {
	emails := make(map[string][]string)
	panels := make(map[string][]api.TablePanel)

	reportGroup, err := re.sql.ReportGroupFromSchedule(schedule)
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.createReport: ReportGroupFromSchedule: " + err.Error())
		return err
	}

	userIDs, err := re.sql.GroupMemberUserIDs(*reportGroup)
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.createReport: GroupMemberUserIDs: " + err.Error())
		return err
	}

	emailsFromUsers, err := api.GetEmails(*authConfig, userIDs, datasourceID)
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.createReport: emailsFromUsers: " + err.Error())
		return err
	}

	emails[schedule.ID] = emailsFromUsers

	reportContent, err := re.sql.GetReportContent(schedule.ID)
	if err != nil {
		log.DefaultLogger.Error("ReportEmailer.createReport: GetReportContent: " + err.Error())
		return err
	}

	panels[schedule.ID] = []api.TablePanel{}

	for _, content := range reportContent {
		interval := int64(schedule.Interval)
		to := strconv.FormatInt(time.Now().Unix(), 10)
		from := strconv.FormatInt(time.Now().Unix()-interval, 10)

		dashboard, err := api.NewDashboard(authConfig, content.DashboardID, from, to, datasourceID)

		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReport: NewDashboard: " + err.Error())
			return err
		}

		panel := dashboard.Panel(content.PanelID)
		if panel != nil {
			panel.PrepSql(dashboard.Variables, content.StoreID, content.Variables)
			panels[schedule.ID] = append(panels[schedule.ID], *panel)
		}
	}

	templatePath := filepath.Join("..", "data", "template.xlsx")
	reporter := reporter.NewReporter(templatePath)

	for scheduleID, reportSheetPanels := range panels {
		report := reporter.CreateNewReport(scheduleID)

		report.SetSheets(reportSheetPanels)
		err := report.Write(*authConfig)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer.createReports: report.Write: " + err.Error())
			return err
		}
	}

	for scheduleID, recipientEmails := range emails {
		schedule, err := re.sql.GetSchedule(scheduleID)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer: GetSchedule: Could not create report to send.", err.Error())
			bugsnag.Notify(err)
		} else {
			attachmentPath := filepath.Join("..", "data", schedule.Name+".xlsx")
			em.BulkCreateAndSend(attachmentPath, recipientEmails, schedule.Name, schedule.Description)
		}
	}

	return nil
	// re.cleanup(schedules)
}

func (re *ReportEmailer) CreateReports() {
	log.DefaultLogger.Info("Creating Reports...")
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

	if len(schedules) > 0 {
		log.DefaultLogger.Info("Found schedules which are overdue...")
		for _, schedule := range schedules {
			log.DefaultLogger.Info(fmt.Sprintf("- %s : %s", schedule.Name, schedule.Description))
		}
	} else {
		log.DefaultLogger.Info("No schedules are overdue...")
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

			if panel != nil {
				panel.PrepSql(dashboard.Variables, content.StoreID, content.Variables)
				panels[schedule.ID] = append(panels[schedule.ID], *panel)
			}

		}
	}

	templatePath := filepath.Join("..", "data", "template.xlsx")

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
		schedule, err := re.sql.GetSchedule(scheduleID)
		if err != nil {
			log.DefaultLogger.Error("ReportEmailer: GetSchedule: Could not create report to send.", err.Error())
			bugsnag.Notify(err)
		} else {
			attachmentPath := filepath.Join("..", "data", schedule.Name+".xlsx")
			em.BulkCreateAndSend(attachmentPath, recipientEmails, schedule.Name, schedule.Description)
		}
	}

	re.cleanup(schedules)

}
