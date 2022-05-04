package datasource

import (
	"database/sql"
	"excel-report-email-scheduler/pkg/ereserror"
	"excel-report-email-scheduler/pkg/setting"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

func (datasource *MsupplyEresDatasource) CreateOrUpdateSettings(settings setting.Settings) error {
	db, err := sql.Open("sqlite", datasource.DataPath)
	if err != nil {
		log.DefaultLogger.Error("CreateOrUpdateSettings: sql.Open(): ", err.Error())
		return err
	}
	defer db.Close()

	exists, err := datasource.settingsExists()
	if err != nil {
		log.DefaultLogger.Error("CreateOrUpdateSettings: datasource.SettingsExist(): ", err.Error())
		return err
	}

	if exists {
		stmt, err := db.Prepare("UPDATE Config set id = ?, grafanaUsername = ?, grafanaPassword = ?, email = ?, emailPassword = ?, datasourceID = ?, emailHost = ?, emailPort = ?, grafanaURL = ?")
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: db.Prepare()1: ", err.Error())
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec("ID", settings.GrafanaUsername, settings.GrafanaPassword, settings.Email, settings.EmailPassword, settings.DatasourceID, settings.EmailHost, settings.EmailPort, settings.GrafanaURL)
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: stmt.Exec()2: ", err.Error())
			return err
		}

	} else {
		stmt, err := db.Prepare("INSERT INTO Config (id, grafanaUsername, grafanaPassword, email, emailPassword, datasourceID, emailHost, emailPort, grafanaURL) VALUES (?,?,?,?,?,?,?,?,?)")
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: db.Prepare()2: ", err.Error())
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec("ID", settings.GrafanaUsername, settings.GrafanaPassword, settings.Email, settings.EmailPassword, settings.DatasourceID, settings.EmailHost, settings.EmailPort, settings.GrafanaURL)
		if err != nil {
			log.DefaultLogger.Error("CreateOrUpdateSettings: stmt.Exec(): ", err.Error())
			return err
		}

	}
	return nil
}

func (datasource *MsupplyEresDatasource) NewSettings() (*setting.Settings, error) {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return nil, err
	}
	defer sqlClient.Db.Close()

	var id, grafanaUsername, grafanaPassword, email, emailPassword, emailHost, grafanaURL string
	var emailPort, datasourceID int

	exists, err := datasource.settingsExists()
	if err != nil {
		log.DefaultLogger.Error("GetSettings: settingsExists(): ", err.Error())
		return nil, err
	}

	if exists {
		rows, err := sqlClient.Db.Query("SELECT id, grafanaUsername, grafanaPassword, email, emailPassword, datasourceID, emailHost, emailPort, grafanaURL FROM Config")
		if err != nil {
			log.DefaultLogger.Error("GetSettings: db.Query(): ", err.Error())
			return nil, err
		}
		defer rows.Close()

		rows.Next()
		err = rows.Scan(&id, &grafanaUsername, &grafanaPassword, &email, &emailPassword, &datasourceID, &emailHost, &emailPort, &grafanaURL)
		if err != nil {
			log.DefaultLogger.Error("GetSettings: rows.Scan(): ", err.Error())
			return nil, err
		}

		return &setting.Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword, DatasourceID: datasourceID, EmailPort: emailPort, EmailHost: emailHost, GrafanaURL: grafanaURL}, nil
	}

	return &setting.Settings{GrafanaUsername: grafanaUsername, GrafanaPassword: grafanaPassword, Email: email, EmailPassword: emailPassword, DatasourceID: datasourceID, EmailPort: emailPort, EmailHost: emailHost, GrafanaURL: grafanaURL}, nil
}

func (datasource *MsupplyEresDatasource) settingsExists() (bool, error) {
	frame := trace()
	sqlClient, err := datasource.NewSqlClient()
	if err != nil {
		err = ereserror.New(500, errors.Wrap(err, frame.Function), "Could not open database")
		return false, err
	}
	defer sqlClient.Db.Close()

	rows, err := sqlClient.Db.Query("SELECT EXISTS(SELECT 1 FROM Config)")
	if err != nil {
		log.DefaultLogger.Error("GetSchedules: db.Query(): ", err.Error())
		return false, err
	}
	defer rows.Close()

	var exists bool
	rows.Next()
	err = rows.Scan(&exists)
	if err != nil {
		log.DefaultLogger.Error("GetSchedules: rows.Scan(): ", err.Error())
		return false, err
	}

	return exists, nil
}
