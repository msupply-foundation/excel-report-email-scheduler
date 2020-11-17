package dbstore

import (
	"database/sql"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Settings struct {
	GrafanaUsername string `json:"grafanaUsername"`
	GrafanaPassword string `json:"grafanaPassword"`
	GrafanaURL      string `json:"grafanaURL"`
	Email           string `json:"email"`
	EmailPassword   string `json:"emailPassword"`
	EmailPort       int    `json:"emailPort"`
	EmailHost       string `json:"emailHost"`
	DatasourceID    int    `json:"datasourceID"`
}

func SettingsFields() string {
	return "\n{\n\tgrafanaUsername string" +
		"\n\tgrafanaPassword string" +
		"\n\tgrafanaURL string" +
		"\n\temail string\n}" +
		"\n\temailPassword string\n}" +
		"\n\temailPort int\n}" +
		"\n\temailHost string\n}" +
		"\n\tDatasourceID int\n}"
}

func (datasource *SQLiteDatasource) settingsExists() (bool, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("settingsExist: sql.Open(): ", err.Error())
		return false, err
	}

	rows, err := db.Query("SELECT EXISTS(SELECT 1 FROM Config)")
	defer rows.Close()
	if err != nil {
		log.DefaultLogger.Error("GetSchedules: db.Query(): ", err.Error())
		return false, err
	}

	var exists bool
	rows.Next()
	err = rows.Scan(&exists)
	if err != nil {
		log.DefaultLogger.Error("GetSchedules: rows.Scan(): ", err.Error())
		return false, err
	}

	return exists, nil
}

func (datasource *SQLiteDatasource) CreateOrUpdateSettings(settings Settings) error {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("CreateOrUpdateSettings: sql.Open(): ", err.Error())
		return err
	}

	exists, err := datasource.settingsExists()
	if err != nil {
		log.DefaultLogger.Error("CreateOrUpdateSettings: datasource.SettingsExist(): ", err.Error())
		return err
	}

	if exists {
		stmt, err := db.Prepare("UPDATE Config set id = ?, grafanaUsername = ?, grafanaPassword = ?, email = ?, emailPassword = ?, datasourceID = ?, emailHost = ?, emailPort = ?, grafanaURL = ?")
		defer stmt.Close()
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: db.Prepare()1: ", err.Error())
			return err
		}

		_, err = stmt.Exec("ID", settings.GrafanaUsername, settings.GrafanaPassword, settings.Email, settings.EmailPassword, settings.DatasourceID, settings.EmailHost, settings.EmailPort, settings.GrafanaURL)
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: stmt.Exec()2: ", err.Error())
			return err
		}

	} else {
		stmt, err := db.Prepare("INSERT INTO Config (id, grafanaUsername, grafanaPassword, email, emailPassword, datasourceID, emailHost, emailPort, grafanaURL) VALUES (?,?,?,?,?,?,?,?,?)")
		defer stmt.Close()
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: db.Prepare()2: ", err.Error())
			return err
		}

		_, err = stmt.Exec("ID", settings.GrafanaUsername, settings.GrafanaPassword, settings.Email, settings.EmailPassword, settings.DatasourceID, settings.EmailHost, settings.EmailPort, settings.GrafanaURL)
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: stmt.Exec(): ", err.Error())
			return err
		}

	}
	return nil
}

func (datasource *SQLiteDatasource) GetSettings() (*Settings, error) {
	db, err := sql.Open("sqlite3", datasource.Path)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("GetSettings: sql.Open(): ", err.Error())
		return nil, err
	}

	var id, grafanaUsername, grafanaPassword, email, emailPassword, emailHost, grafanaURL string
	var emailPort, datasourceID int

	exists, err := datasource.settingsExists()
	if err != nil {
		log.DefaultLogger.Error("GetSettings: settingsExists(): ", err.Error())
		return nil, err
	}

	if exists {
		rows, err := db.Query("SELECT id, grafanaUsername, grafanaPassword, email, emailPassword, datasourceID, emailHost, emailPort, grafanaURL FROM Config")
		defer rows.Close()
		if err != nil {
			log.DefaultLogger.Error("GetSettings: db.Query(): ", err.Error())
			return nil, err
		}

		rows.Next()
		err = rows.Scan(&id, &grafanaUsername, &grafanaPassword, &email, &emailPassword, &datasourceID, &emailHost, &emailPort, &grafanaURL)
		if err != nil {
			log.DefaultLogger.Error("GetSettings: rows.Scan(): ", err.Error())
			return nil, err
		}

		return &Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword, DatasourceID: datasourceID, EmailPort: emailPort, EmailHost: emailHost, GrafanaURL: grafanaURL}, nil
	}

	return &Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword, DatasourceID: datasourceID, EmailPort: emailPort, EmailHost: emailHost, GrafanaURL: grafanaURL}, nil
}
