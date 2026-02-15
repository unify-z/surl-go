package helpers

import (
	gomail "gopkg.in/gomail.v2"
)

type SMTPHelper struct {
	d        gomail.Dialer
	user     string
	host     string
	port     int
	password string
	skipTLS  bool
}

func NewSMTPHelper(host string, port int, user, password string, skipTLS bool) *SMTPHelper {
	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = nil
	return &SMTPHelper{
		d:        *d,
		user:     user,
		host:     host,
		port:     port,
		password: password,
		skipTLS:  skipTLS,
	}
}

func (s *SMTPHelper) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	return s.d.DialAndSend(m)
}
