package auth

import (
	"excel-report-email-scheduler/pkg/dbstore"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type EmailConfig struct {
	Email    string
	Password string
	Host     string
	Port     int
}

func NewEmailConfig(datasource *dbstore.SQLiteDatasource) (*EmailConfig, error) {

	settings, err := datasource.GetSettings()
	if err != nil {
		log.DefaultLogger.Error("NewEmailConfig: datasource.GetSettings(): ", err.Error())
		return nil, err
	}

	return &EmailConfig{Email: settings.Email, Password: settings.EmailPassword, Host: settings.EmailHost, Port: settings.EmailPort}, nil
}
