package main

import (
	"context"
	"database/sql"
	_ "database/sql/driver"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	_ "github.com/mattn/go-sqlite3"
)

// TODOs:
// SQL In separate files.
// Repository structs for each table? i.e. `ScheduleRepository` which has the methods `ScheduleRepository.update()` etc to partition the datasource struct.
// Common serialization/deserialization methods
// Better/consistent erroring
// Better/consistent return values

// Basic SQLite datasource
type SQLiteDatasource struct {
	instanceManager instancemgmt.InstanceManager
	path            string
}

// Basic plugin instance settings
type InstanceSettings struct {
	httpClient *http.Client
}

type queryModel struct {
	Format string `json:"format"`
}

func getDataSource() (*SQLiteDatasource, error) {
	instanceManager := datasource.NewInstanceManager(getDataSourceInstanceSettings)

	sqlDatasource := &SQLiteDatasource{
		instanceManager: instanceManager,
		path:            "./data/msupply.db",
	}

	sqlDatasource.Init()

	return sqlDatasource, nil
}

func getDataSourceInstanceSettings(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &InstanceSettings{
		httpClient: &http.Client{},
	}, nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (datasource *SQLiteDatasource) QueryData(ctx context.Context, request *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if os.Getenv("APP_ENV") != "production" {
		log.DefaultLogger.Info("QueryData", "request", request)
	}

	// TODO: Use context to timeout/cancel?
	response := backend.NewQueryDataResponse()

	for _, query := range request.Queries {
		res := datasource.query(ctx, query)
		response.Responses[query.RefID] = res
	}

	return response, nil
}

func (datasource *SQLiteDatasource) query(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	var queryModel queryModel

	response := backend.DataResponse{}
	response.Error = json.Unmarshal(query.JSON, &queryModel)

	if response.Error != nil {
		return response
	}

	if queryModel.Format == "" {
		log.DefaultLogger.Warn("format is empty. defaulting to time series")
	}

	frame := data.NewFrame("response")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("values", nil, []int64{10, 20}),
	)
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (datasource *SQLiteDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Yeah, nah, All good"

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (datasource *SQLiteDatasource) Init() {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.DefaultLogger.Warn("Could not ping database")
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS Schedule (id TEXT PRIMARY KEY, interval INTEGER, nextReportTime INTEGER, name TEXT, description TEXT, lookback INTEGER, reportGroupID TEXT, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id))")
	stmt.Exec()
	defer stmt.Close()

	if err != nil {
		log.DefaultLogger.Warn("Could not create table!")
		log.DefaultLogger.Warn(err.Error())
	}
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS Config (id TEXT PRIMARY KEY, grafanaUsername TEXT, grafanaPassword TEXT, emailPassword TEXT, email TEXT)")
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportRecipient (id TEXT PRIMARY KEY, userID TEXT)")
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportGroup (id TEXT PRIMARY KEY, name TEXT, description TEXT)")
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportGroupMembership (id TEXT PRIMARY KEY, userID TEXT, reportGroupID TEXT, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id))")
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS GroupSchedule (id TEXT PRIMARY KEY, scheduleID TEXT, reportGroupID TEXT, FOREIGN KEY(reportGroupID) REFERENCES ReportGroup(id), FOREIGN KEY(scheduleID) REFERENCES Schedule(id))")
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS ReportContent (id TEXT PRIMARY KEY, scheduleID TEXT, panelID INTEGER, lookback INTEGER, storeID TEXT, FOREIGN KEY(scheduleID) REFERENCES Schedule(id))")
	stmt.Exec()

}

func (datasource *SQLiteDatasource) settingsExists() bool {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var exists bool
	rows, _ := db.Query("SELECT EXISTS(SELECT 1 FROM Config)")

	defer rows.Close()
	rows.Next()
	rows.Scan(&exists)

	log.DefaultLogger.Warn(string(strconv.FormatBool(exists)))
	return exists
}

func (datasource *SQLiteDatasource) createOrUpdateSettings(config Config) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	if datasource.settingsExists() {
		stmt, _ := db.Prepare("UPDATE Config set id = ?, grafanaUsername = ?, grafanaPassword = ?, email = ?, emailPassword = ?")
		stmt.Exec("ID", config.GrafanaUsername, config.GrafanaPassword, config.Email, config.EmailPassword)
		stmt.Close()
	} else {
		stmt, _ := db.Prepare("INSERT INTO Config (id, grafanaUsername, grafanaPassword, email, emailPassword) VALUES (?,?,?,?,?)")
		stmt.Exec("ID", config.GrafanaUsername, config.GrafanaPassword, config.Email, config.EmailPassword)
		stmt.Close()
	}
	return true, nil
}

func (datasource *SQLiteDatasource) getSettings() *Config {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var grafanaUsername, grafanaPassword, email, emailPassword string

	if datasource.settingsExists() {
		var id, grafanaUsername, grafanaPassword, email, emailPassword string
		rows, _ := db.Query("SELECT * FROM Config")
		defer rows.Close()
		rows.Next()
		rows.Scan(&id, &grafanaUsername, &grafanaPassword, &email, &emailPassword)
		log.DefaultLogger.Warn(id, grafanaUsername, grafanaPassword, email, emailPassword)
		return &Config{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword}
	}
	log.DefaultLogger.Warn("found")

	return &Config{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword}
}

func (datasource *SQLiteDatasource) createSchedule() (Schedule, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	newUuid := uuid.New().String()
	schedule := Schedule{ID: newUuid, NextReportTime: 0, Interval: 0, Name: "", Description: "", Lookback: 0, ReportGroupID: ""}
	stmt, _ := db.Prepare("INSERT INTO Schedule (ID,  nextReportTime, interval, name, description, lookback, reportGroupID) VALUES (?,?,?,?,?,?,?)")
	stmt.Exec(newUuid, 0, 60*1000*60*24, "New report schedule", "", 0, "")
	defer stmt.Close()

	return schedule, nil
}

func (datasource *SQLiteDatasource) updateSchedule(id string, schedule Schedule) *Schedule {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE Schedule SET nextReportTime = ?, interval = ?, name = ?, description = ?, lookback = ?, reportGroupID = ? where id = ?")
	stmt.Exec(schedule.NextReportTime, schedule.Interval, schedule.Name, schedule.Description, schedule.Lookback, schedule.ReportGroupID, id)
	defer stmt.Close()

	return nil
}

func (datasource *SQLiteDatasource) deleteSchedule(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM Schedule WHERE id = ?")
	stmt.Exec(id)
	stmt, _ = db.Prepare("DELETE FROM ReportContent WHERE scheduleID = ?")
	stmt.Exec(id)

	defer stmt.Close()

	return true, nil
}

func (datasource *SQLiteDatasource) getSchedules() []Schedule {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var schedules []Schedule

	rows, _ := db.Query("SELECT * FROM Schedule")
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description, ReportGroupID string
		var Interval, NextReportTime, Lookback int

		rows.Scan(&ID, &Interval, &NextReportTime, &Name, &Description, &Lookback, &ReportGroupID)
		schedule := Schedule{ID, Interval, NextReportTime, Name, Description, Lookback, ReportGroupID}

		schedules = append(schedules, schedule)
	}

	return schedules
}

func (datasource *SQLiteDatasource) createReportRecipient(reportRecipient ReportRecipient) *ReportRecipient {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("INSERT INTO ReportRecipient (ID, userID) VALUES (?,?)")
	stmt.Exec("a", reportRecipient.UserID)

	defer stmt.Close()

	return nil
}

func (datasource *SQLiteDatasource) updateReportRecipient(id string, reportRecipient ReportRecipient) *ReportRecipient {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE ReportRecipient SET UserID = ? WHERE id = ?")
	stmt.Exec(reportRecipient.UserID, id)
	defer stmt.Close()

	return nil
}

func (datasource *SQLiteDatasource) deleteReportRecipient(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportRecipient WHERE ID = ?")
	stmt.Exec(id)
	defer stmt.Close()

	return true, nil
}

func (datasource *SQLiteDatasource) getReportRecipients() []ReportRecipient {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var recipients []ReportRecipient

	rows, _ := db.Query("SELECT * FROM ReportRecipient")
	defer rows.Close()

	for rows.Next() {
		var ID, UserID string
		rows.Scan(&ID, &UserID)
		recipient := ReportRecipient{ID, UserID}
		recipients = append(recipients, recipient)
	}

	return recipients
}

func (datasource *SQLiteDatasource) getReportRecipient(id string) ReportRecipient {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	row := db.QueryRow("SELECT * FROM ReportRecipient WHERE ID = ?", id)

	// TODO: Handle the case where it doesn't exist
	var ID, UserID string
	row.Scan(&ID, &UserID)
	recipient := ReportRecipient{ID, UserID}

	return recipient
}

func (datasource *SQLiteDatasource) getReportGroups() []ReportGroup {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var reportGroups []ReportGroup

	rows, _ := db.Query("SELECT * FROM ReportGroup")
	defer rows.Close()

	for rows.Next() {
		var ID, Name, Description string
		rows.Scan(&ID, &Name, &Description)
		reportGroup := ReportGroup{ID, Name, Description}
		reportGroups = append(reportGroups, reportGroup)
	}

	return reportGroups
}

func (datasource *SQLiteDatasource) createReportGroup() (ReportGroup, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	reportGroup := ReportGroup{ID: uuid.New().String(), Name: "New report group", Description: ""}
	stmt, _ := db.Prepare("INSERT INTO ReportGroup (id, name, description) VALUES (?,?,?)")
	stmt.Exec(reportGroup.ID, reportGroup.Name, reportGroup.Description)
	defer stmt.Close()

	return reportGroup, nil
}

func (datasource *SQLiteDatasource) updateReportGroup(id string, reportGroup ReportGroup) ReportGroup {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE ReportGroup SET name = ?, description = ? where id = ?")
	stmt.Exec(reportGroup.Name, reportGroup.Description, id)
	defer stmt.Close()

	return reportGroup
}

func (datasource *SQLiteDatasource) deleteReportGroup(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportGroup WHERE id = ?")
	stmt.Exec(id)
	stmt, _ = db.Prepare("DELETE FROM ReportGroupMembership WHERE reportGroupID = ?")
	stmt.Exec(id)

	defer stmt.Close()

	return true, nil
}

func (datasource *SQLiteDatasource) getReportGroupMemberships(groupID string) []ReportGroupMembership {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var memberships []ReportGroupMembership

	rows, _ := db.Query("SELECT * FROM ReportGroupMembership WHERE reportGroupID = ?", groupID)
	defer rows.Close()

	for rows.Next() {
		var ID, UserID, ReportGroupID string
		rows.Scan(&ID, &UserID, &ReportGroupID)
		membership := ReportGroupMembership{ID, UserID, ReportGroupID}
		memberships = append(memberships, membership)
	}

	return memberships
}

func (datasource *SQLiteDatasource) createReportGroupMembership(members []ReportGroupMembership) ([]ReportGroupMembership, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var addedMemberships []ReportGroupMembership
	for _, member := range members {
		newUuid := uuid.New().String()
		stmt, _ := db.Prepare("INSERT INTO ReportGroupMembership (ID, userID, reportGroupID) VALUES (?,?,?)")
		stmt.Exec(newUuid, member.ReportRecipientID, member.ReportGroupID)
		member.ID = newUuid
		addedMemberships = append(addedMemberships, member)
		defer stmt.Close()
	}

	// TODO: Return report assignment
	return addedMemberships, nil
}

func (datasource *SQLiteDatasource) deleteReportGroupMembership(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportGroupMembership WHERE id = ?")
	stmt.Exec(id)
	defer stmt.Close()

	// TODO: Proper return values, returning error or false? or just an error, probably
	return true, nil
}

func (datasource *SQLiteDatasource) getReportContent(scheduleID string) ([]ReportContent, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var reportContent []ReportContent

	rows, _ := db.Query("SELECT * FROM ReportContent WHERE scheduleID = ?", scheduleID)
	defer rows.Close()

	for rows.Next() {
		var ID, ScheduleID, StoreID string
		var Lookback, PanelID int
		rows.Scan(&ID, &ScheduleID, &PanelID, &Lookback, &StoreID)
		content := ReportContent{ID, ScheduleID, PanelID, Lookback, StoreID}
		reportContent = append(reportContent, content)
	}

	return reportContent, nil
}

func (datasource *SQLiteDatasource) deleteReportContent(id string) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM ReportContent WHERE ID = ?")
	stmt.Exec(id)
	defer stmt.Close()

	return true, nil
}

func (datasource *SQLiteDatasource) createReportContent(newReportContentValues ReportContent) (ReportContent, error) {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	reportContent := ReportContent{ID: uuid.New().String(), ScheduleID: newReportContentValues.ScheduleID, PanelID: newReportContentValues.PanelID, Lookback: 0, StoreID: ""}
	stmt, _ := db.Prepare("INSERT INTO ReportContent (id, scheduleID, panelID, lookback, storeID) VALUES (?,?,?,?,?)")
	stmt.Exec(reportContent.ID, reportContent.ScheduleID, reportContent.PanelID, reportContent.Lookback, reportContent.StoreID)
	defer stmt.Close()

	return reportContent, nil
}

func (datasource *SQLiteDatasource) updateReportContent(id string, reportGroup ReportContent) ReportContent {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE ReportContent SET scheduleID = ?, panelID = ?, storeID = ?, lookback = ? where id = ?")
	stmt.Exec(reportGroup.ScheduleID, reportGroup.PanelID, reportGroup.StoreID, reportGroup.Lookback, id)
	defer stmt.Close()

	return reportGroup
}

// TODO: Handle error cases and also might need to add additional
// fields i.e. SMTP etc
func (datasource *SQLiteDatasource) getEmailConfig() EmailConfig {
	db, _ := sql.Open("sqlite3", datasource.path)
	defer db.Close()

	var email, password string

	row := db.QueryRow("SELECT email, emailPassword as password FROM Config")
	row.Scan(&email, &password)

	return EmailConfig{email: email, password: password}

}
