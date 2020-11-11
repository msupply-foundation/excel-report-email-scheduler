package auth

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
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
		log.DefaultLogger.Error("NewAuthConfig: datasource.GetSettings(): ", err.Error())
		return nil, err
	}

	return &EmailConfig{Email: settings.Email, Password: settings.EmailPassword, Host: settings.EmailHost, Port: settings.EmailPort}, nil
}
