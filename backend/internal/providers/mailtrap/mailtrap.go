package mailtrap

import (
	"fmt"
	"net/smtp"
)

type MailtrapProvider struct {
	host      string
	port      string
	username  string
	password  string
	fromEmail string
}

func New(host, port, username, password, fromEmail string) *MailtrapProvider {
	return &MailtrapProvider{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
	}
}

func (p *MailtrapProvider) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", p.username, p.password, p.host)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		p.fromEmail, to, subject, body))

	addr := fmt.Sprintf("%s:%s", p.host, p.port)

	if err := smtp.SendMail(addr, auth, p.fromEmail, []string{to}, msg); err != nil {
		return fmt.Errorf("send email to %s: %w", to, err)
	}

	return nil
}
