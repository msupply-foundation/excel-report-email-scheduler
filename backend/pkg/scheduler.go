package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/api"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
	dbstore "github.com/grafana/simple-datasource-backend/pkg/db"
	"github.com/grafana/simple-datasource-backend/pkg/report"
	"gopkg.in/gomail.v2"
)

func sendEmail(attachmentPath string, email string, emailConfig auth.EmailConfig) {
	log.DefaultLogger.Info("Sending Email!!:" + " " + email)
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.Email)
	m.SetHeader("To", email)

	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello")
	m.Attach(attachmentPath)

	// // I don't really know what I'm doing with this auth.
	// // PlainAuth works and reading the docs it seems to fail
	// // if not using TLS. So I guess it's probably OK.
	// // TODO: Host and port need to be added to datasource config?
	// // This password is an app-specific password. The real password
	// // to the account is kathmandu312. Seems to require me to generate
	// // and use an app-specific password. :shrug: // "ybtkmpesjptowmru"
	d := gomail.NewDialer("smtp.gmail.com", 587, emailConfig.Email, "ybtkmpesjptowmru")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

func sendEmails(attachmentPath string, emails []string, emailConfig auth.EmailConfig) {
	for _, email := range emails {
		sendEmail(attachmentPath, email, emailConfig)
	}
}

func createReports(sql *dbstore.SQLiteDatasource) {
	authConfig := auth.NewAuthConfig(sql)
	emailConfig := auth.NewEmailConfig(sql)

	schedules, _ := sql.OverdueSchedules()

	emails := make(map[string][]string)
	panels := make(map[string][]api.TablePanel)

	for _, schedule := range schedules {
		reportGroup := sql.ReportGroupFromSchedule(schedule)
		userIDs, _ := sql.GroupMemberUserIDs(*reportGroup)
		emails[schedule.ID] = api.GetEmails(*authConfig, userIDs)

		reportContent, _ := sql.GetReportContent(schedule.ID)
		panels[schedule.ID] = []api.TablePanel{}

		for _, content := range reportContent {
			dashboard, _ := api.NewDashboard(authConfig, content.DashboardID)
			panel, _ := dashboard.Panel(content.PanelID)
			panels[schedule.ID] = append(panels[schedule.ID], *panel)
		}
	}

	reporter := report.NewReporter("./data/template.xlsx")

	for scheduleID, reportSheetPanels := range panels {
		report := reporter.CreateNewReport(scheduleID)
		report.SetSheets(reportSheetPanels)
		report.Write(*authConfig)
	}

	for scheduleID, recipientEmails := range emails {
		log.DefaultLogger.Info("Sending Emails!!")
		sendEmails("./data/"+scheduleID+".xlsx", recipientEmails, emailConfig)
	}

}

func cleanup() {
	// 	// Not sure exactly but need to:
	// 	// 1. Delete all temp attachment files?
	// 	// 2. Set the new times on all of the schedules
}

// func getScheduler(sqlite *SQLiteDatasource) func() {
// 	return func() {
// 		log.DefaultLogger.Info("Scheduler!")
// 		createReports()
// 		sendEmail()
// 		log.DefaultLogger.Info("Scheduler2!")

// 	}
