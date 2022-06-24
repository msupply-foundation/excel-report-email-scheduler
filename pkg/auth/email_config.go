package auth

import (
	"excel-report-email-scheduler/pkg/setting"
)

func NewEmailConfig(settings *setting.Settings) (*EmailConfig, error) {
	return &EmailConfig{Email: settings.Email, Password: settings.EmailPassword, Host: settings.EmailHost, Port: settings.EmailPort}, nil
}
