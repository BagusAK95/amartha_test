package mail

import (
	"bytes"

	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/utils/html"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSender(cfg config.MailConfig) *Sender {
	return &Sender{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
	}
}

func (s *Sender) SendEmailWithTemplate(to, subject, file string, data any) error {
	tmpl, err := html.NewTemplate()
	if err != nil {
		return err
	}

	// Execute the template with the provided data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, file, data); err != nil {
		return err
	}

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", s.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Create a new dialer
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
