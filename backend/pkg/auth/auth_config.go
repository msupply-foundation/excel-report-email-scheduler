package auth

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

type AuthConfig struct {
	Username string
	Password string
	URL      string
}

func NewAuthConfig(datasource *dbstore.SQLiteDatasource) (*AuthConfig, error) {
	settings, err := datasource.GetSettings()
	if err != nil {
		log.DefaultLogger.Error("NewAuthConfig: datasource.GetSettings(): ", err.Error())
		return nil, err
	}

	return &AuthConfig{Username: settings.GrafanaUsername, Password: settings.GrafanaPassword, URL: settings.GrafanaURL}, nil
}

func (config AuthConfig) AuthString() string {
	return config.Username + ":" + config.Password + "@"
}
