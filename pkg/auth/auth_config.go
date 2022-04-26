package auth

import (
	"excel-report-email-scheduler/pkg/ereserror"
	"excel-report-email-scheduler/pkg/setting"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type AuthConfig struct {
	Username string
	Password string
	URL      string
}

func NewAuthConfig(settings *setting.Settings) (*AuthConfig, error) {
	return &AuthConfig{Username: settings.GrafanaUsername, Password: settings.GrafanaPassword, URL: settings.GrafanaURL}, nil
}

func (config AuthConfig) AuthString() string {
	return config.Username + ":" + config.Password + "@"
}

func (config AuthConfig) AuthURL() (*string, error) {
	authUrl, err := config.InjectAuthString()
	if err != nil {
		return nil, err
	}
	authUrl1 := strings.TrimSuffix(*authUrl, "/")
	return &authUrl1, nil
}

func (config AuthConfig) InjectAuthString() (*string, error) {
	http := regexp.MustCompile("http://")
	https := regexp.MustCompile("https://")

	index := http.FindStringIndex(config.URL)
	if index == nil {
		index = https.FindStringIndex(config.URL)
		if index == nil {
			err := errors.New("Could not authorise")
			err = ereserror.New(500, err, err.Error())
			return nil, err
		}
	}

	authURL := config.URL[:index[1]] + config.AuthString() + config.URL[index[1]:]

	return &authURL, nil
}
