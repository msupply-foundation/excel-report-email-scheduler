package emailer

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
	"gopkg.in/gomail.v2"
)

type Emailer struct {
	email    string
	password string
}

func New(config *auth.EmailConfig) *Emailer {
	return &Emailer{email: config.Email, password: config.Password}
}

func (e *Emailer) CreateAndSend(attachmentPath string, email string) {
	m := gomail.NewMessage()

	m.SetHeader("From", e.email)
	m.SetHeader("To", email)

	// m.SetHeader("Subject", "Hello!")
	// m.SetBody("text/html", "Hello")
	m.Attach(attachmentPath)

	// // I don't really know what I'm doing with this auth.
	// // PlainAuth works and reading the docs it seems to fail
	// // if not using TLS. So I guess it's probably OK.
	// // TODO: Host and port need to be added to datasource config?
	// // This password is an app-specific password. The real password
	// // to the account is kathmandu312. Seems to require me to generate
	// // and use an app-specific password. :shrug: // "ybtkmpesjptowmru"
	// d := gomail.NewDialer("smtp.gmail.com", 587, e.email, "ybtkmpesjptowmru")

	if err := d.DialAndSend(m); err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

func (e *Emailer) BulkCreateAndSend(attachmentPath string, emails []string) {
	for _, email := range emails {
		e.CreateAndSend(attachmentPath, email)
	}
}
