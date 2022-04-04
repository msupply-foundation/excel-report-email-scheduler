package emailer

import (
	"fmt"

	"excel-report-email-scheduler/pkg/auth"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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

func (e *Emailer) CreateAndSend(attachmentPath, email, subject, body string) error {
	log.DefaultLogger.Info(fmt.Sprintf("Sending email to %s...", email))
	m := gomail.NewMessage()

	m.SetHeader("From", e.email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	m.Attach(attachmentPath)
	d := gomail.NewDialer(e.host, e.port, e.email, e.password)

	if err := d.DialAndSend(m); err != nil {
		log.DefaultLogger.Error("CreateAndSend: DialAndSend: " + err.Error())
		return err
	}

	log.DefaultLogger.Info(fmt.Sprintf("Sent email to %s!", email))
	return nil
}

func (e *Emailer) BulkCreateAndSend(attachmentPath string, emails []string, subject string, body string) {
	for _, email := range emails {
		if err := e.CreateAndSend(attachmentPath, email, subject, body); err != nil {
			log.DefaultLogger.Error("BulkCreateAndSend: Could not send to: " + email)
		}
	}
}
