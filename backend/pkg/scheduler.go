package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/api"
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	email    string
	password string
}

// Possibly take the email config, an attachment and some array of recipients?
func sendEmail() {

	// // TODO: Password and email from params. Leave hard coded for now
	m := gomail.NewMessage()
	m.SetHeader("From", "testemailsussol@gmail.com")
	m.SetHeader("To", "griffinjoshua5@gmail.com")

	// // TODO: Subject, message to be added to datasource config? Or schedule config?
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello")
	m.Attach("./data/Book1.xlsx")

	// // I don't really know what I'm doing with this auth.
	// // PlainAuth works and reading the docs it seems to fail
	// // if not using TLS. So I guess it's probably OK.
	// // TODO: Host and port need to be added to datasource config?
	// // This password is an app-specific password. The real password
	// // to the account is kathmandu312. Seems to require me to generate
	// // and use an app-specific password. :shrug:
	d := gomail.NewDialer("smtp.gmail.com", 587, "testemailsussol@gmail.com", "ybtkmpesjptowmru")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		// log.DefaultLogger.Error(err.Error())
	}
}

// TODO: Handle cases where an email config doesn't exist or invalid etc.
// Could send an email to somewhere else, sussol support or something
// if this is the case?
func getEmailConfig(sqlite *SQLiteDatasource) EmailConfig {
	return sqlite.getEmailConfig()
}

func getSchedulesToReportOn(sqlite *SQLiteDatasource) {
	// const schedules = sqlite.getSchedules()
	// const timeNow = Date.now()
	// const schedulesToReportOn = schedules.filter(schedule => schedule.nextReportTime < timeNow)

	// return schedulesToReportOn
}

func getUsersToReportTo(sqlite *SQLiteDatasource) {

	// const schedules = getSchedulesToReportOn()

	// Get the report groups assigned to the schedules that are being reported on.
	// const reportGroups = sqlite.getReportGroups(schedules)

	// Find the report group members (Join table between report group and user)
	// const reportGroupMembers = sqlite.getReportGroupMembers(reportGroups)

	// const users = need to do an http call here to the postgres db with the userIDs from the report group member records.

	// const lookup = create some sort of lookup table with a shape:
	// { {scheduleID}: [user emails] }

	// return users
}

func getPanelsForSchedules(sqlite *SQLiteDatasource) {
	// Get the schedules
	// const schedules = getSchedulesToReportOn()

	// Create a lookup in the shape:
	// { {scheduleID}: [panelIDs] }
	// const lookup = schedules.reduce((acc, value) => ({...acc, [value.id]: sqlite.getPanels(scheduleID)}) , {})
}

func createReports() {
	// Get the panels for each schedule
	// const panelLookup = getPanelsForSchedules()

	// For each schedule, for each panel, query for the panel data
	// and create an excel file where each tab is a panel table.
	// save in a new lookup the shape:
	// { {scheduleID}: excelFilePath }

	dashboard, err := api.GetDashboard("")

	if err != nil {
		log.DefaultLogger.Error(err.Error())
	}

	rawSQL := dashboard.GetRawSQL(2)
	queryRequest, err := api.NewQueryRequest(rawSQL).ToRequestBody()

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		// return nil, err
	}

	response, err := http.Post("http://admin:admin@localhost:3000/api/tsdb/query", "application/json", queryRequest)

	if err != nil {
		log.DefaultLogger.Error(err.Error())
		// return nil, err
	}

	qr, err := api.QueryFromResponse(response)

	// Open the template
	f, e := excelize.OpenFile("./data/template.xlsx")

	// Regexp which matches a cell reference [TODO: Handle $ refs]
	reg := regexp.MustCompile("([A-Za-z]+)|([0-9]+)")

	// Search for the {title} cell to add into the template
	result := f.SearchSheet("Sheet1", "{{title}}")
	f.SetCellValue("Sheet1", result[0], "This is a dynamically added title")

	// Same for the date
	result = f.SearchSheet("Sheet1", "{{date}}")
	f.SetCellValue("Sheet1", result[0], "This is a dynamically added date")

	// Find where to start the rows from
	rowStart := f.SearchSheet("Sheet1", "{{rows}}")
	split := reg.FindAllString(rowStart[0], 2)
	rowIndex, _ := strconv.Atoi(split[1])

	// Same with the headers
	headerStart := f.SearchSheet("Sheet1", "{{headers}}")
	split = reg.FindAllString(headerStart[0], 2)
	headerIndex, _ := strconv.Atoi(split[1])

	if e != nil {
		log.DefaultLogger.Error("erorr opening", e.Error())
	}

	// Replace the header row with the column text values.
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}
	ci := 0
	for _, column := range qr.Results.A.Tables[0].Columns {
		cellRef := columns[ci] + strconv.Itoa(headerIndex)
		f.SetCellValue("Sheet1", cellRef, column.Text)
		ci += 1
	}

	// Duplicate enough rows
	ri := rowIndex + 1
	i := 0
	for i < len(qr.Results.A.Tables[0].Rows)-1 {
		f.DuplicateRow("Sheet1", rowIndex)
		i += 1
	}

	// Insert row content into each duplicated row
	ri = rowIndex
	for _, row := range qr.Results.A.Tables[0].Rows {
		ci = 0
		for range qr.Results.A.Tables[0].Columns {
			cellRef := columns[ci] + strconv.Itoa(ri)
			f.SetCellValue("Sheet1", cellRef, row[ci])
			ci += 1
		}
		ri += 1
	}

	if err := f.SaveAs("./data/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func sendEmails(sqlite *SQLiteDatasource) {
	// const attachments = createReports()
	// const recipients = getUsersToReportTo()

	// We could do some intersection and send a single
	// email with multiple attachments to some users
	// or combine the attachments or something like
	// that.. but seems too hard basket

	// for each key, value in attachments
	// user the key to lookup in the recipients
	// lookup to get the recipient emails.
	// and call sendEmail(config, recipients, attachment)
}

func cleanup() {
	// Not sure exactly but need to:
	// 1. Delete all temp attachment files?
	// 2. Set the new times on all of the schedules
}

func getScheduler(sqlite *SQLiteDatasource) func() {
	return func() {
		log.DefaultLogger.Info("Scheduler!")
		createReports()
		sendEmail()
		log.DefaultLogger.Info("Scheduler2!")

	}
}
