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

func (e *Emailer) CreateAndSend(attachmentPath string, email string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.email)
	m.SetHeader("To", email)

	m.Attach(attachmentPath)
	d := gomail.NewDialer(e.host, e.port, e.email, e.password)

	if err := d.DialAndSend(m); err != nil {
		log.DefaultLogger.Error("CreateAndSend: DialAndSend: " + err.Error())
		return err
	}

	return nil
}

func (e *Emailer) BulkCreateAndSend(attachmentPath string, emails []string) {
	for _, email := range emails {
		if err := e.CreateAndSend(attachmentPath, email); err != nil {
			log.DefaultLogger.Error("BulkCreateAndSend: Could not send to: " + email)
		}

	}
}
