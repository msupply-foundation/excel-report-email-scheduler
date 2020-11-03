package auth

import (
	"github.com/grafana/simple-datasource-backend/pkg/dbstore"
)

type AuthConfig struct {
	Username string
	Password string
}

func NewAuthConfig(datasource *dbstore.SQLiteDatasource) *AuthConfig {
	settings := datasource.GetSettings()
	return &AuthConfig{Username: settings.GrafanaUsername, Password: settings.GrafanaPassword}
}

func (config AuthConfig) AuthString() string {
	return config.Username + ":" + config.Password + "@"
}
