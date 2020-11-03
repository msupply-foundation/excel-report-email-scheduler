package dbstore

import (
	"database/sql"
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Settings struct {
	GrafanaUsername string `json:"grafanaUsername"`
	GrafanaPassword string `json:"grafanaPassword"`
	Email           string `json:"email"`
	EmailPassword   string `json:"emailPassword"`
}

func (datasource *SQLiteDatasource) settingsExists() bool {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var exists bool
	rows, _ := db.Query("SELECT EXISTS(SELECT 1 FROM Config)")

	defer rows.Close()
	rows.Next()
	rows.Scan(&exists)

	log.DefaultLogger.Warn(string(strconv.FormatBool(exists)))
	return exists
}

func (datasource *SQLiteDatasource) CreateOrUpdateSettings(settings Settings) (bool, error) {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	if datasource.settingsExists() {
		stmt, _ := db.Prepare("UPDATE Config set id = ?, grafanaUsername = ?, grafanaPassword = ?, email = ?, emailPassword = ?")
		stmt.Exec("ID", settings.GrafanaUsername, settings.GrafanaPassword, settings.Email, settings.EmailPassword)
		stmt.Close()
	} else {
		stmt, _ := db.Prepare("INSERT INTO Config (id, grafanaUsername, grafanaPassword, email, emailPassword) VALUES (?,?,?,?,?)")
		stmt.Exec("ID", settings.GrafanaUsername, settings.GrafanaPassword, settings.Email, settings.EmailPassword)
		stmt.Close()
	}
	return true, nil
}

func (datasource *SQLiteDatasource) GetSettings() *Settings {
	db, _ := sql.Open("sqlite3", datasource.Path)
	defer db.Close()

	var grafanaUsername, grafanaPassword, email, emailPassword string

	if datasource.settingsExists() {
		var id, grafanaUsername, grafanaPassword, email, emailPassword string
		rows, _ := db.Query("SELECT * FROM Config")
		defer rows.Close()
		rows.Next()
		rows.Scan(&id, &grafanaUsername, &grafanaPassword, &email, &emailPassword)
		log.DefaultLogger.Warn(id, grafanaUsername, grafanaPassword, email, emailPassword)
		return &Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword}
	}
	log.DefaultLogger.Warn("found")

	return &Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword}
}
