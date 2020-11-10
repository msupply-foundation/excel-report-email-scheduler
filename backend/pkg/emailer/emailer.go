package emailer

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
	"gopkg.in/gomail.v2"
)

type Emailer struct {
	email    string
	password string
	host     string
	port     int
}

func New(config *auth.EmailConfig) *Emailer {
	return &Emailer{email: config.Email, password: config.Password, host: config.Host, port: config.Port}
}

func (e *Emailer) CreateAndSend(attachmentPath string, email string) {
	m := gomail.NewMessage()

	m.SetHeader("From", e.email)
	m.SetHeader("To", email)

	// m.SetHeader("Subject", "Hello!")
	// m.SetBody("text/html", "Hello")
	m.Attach(attachmentPath)
	d := gomail.NewDialer(e.host, e.port, e.email, e.password)

	if err := d.DialAndSend(m); err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

func (e *Emailer) BulkCreateAndSend(attachmentPath string, emails []string) {
	for _, email := range emails {
		e.CreateAndSend(attachmentPath, email)
	}
}
