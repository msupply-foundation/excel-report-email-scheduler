package auth

import (
	"excel-report-email-scheduler/pkg/dbstore"
	"regexp"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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

func (config AuthConfig) AuthURL() string {
	authUrl := config.InjectAuthString()
	if strings.HasSuffix(authUrl, "/") {
		authUrl = authUrl[:len(authUrl)-1]
	}
	log.DefaultLogger.Debug("auth url: " + authUrl)

	return authUrl
}

func (config AuthConfig) InjectAuthString() string {
	http := regexp.MustCompile("http://")
	https := regexp.MustCompile("https://")

	index := http.FindStringIndex(config.URL)
	if index == nil {
		index = https.FindStringIndex(config.URL)
		if index == nil {
			log.DefaultLogger.Info("Error injecting Auth: " + config.URL)
			return ""
		}
	}

	return config.URL[:index[1]] + config.AuthString() + config.URL[index[1]:]
}
