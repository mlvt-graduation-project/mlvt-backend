package email_service

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"mlvt/internal/infra/env"
	"net/smtp"
)

type EmailService interface {
	SendHTMLEmail(subject string, htmlBody string, receiverEmail string) error
	CreateAccountSignUpEmail(username string, token string) (string, error)
}

type emailService struct {
	SMTPEmail    string
	SMTPPassword string
	SMTPHost     string
	SMTPPort     string
}

func NewEmailService() EmailService {
	return &emailService{
		SMTPEmail:    env.EnvConfig.SMTPEmail,
		SMTPPassword: env.EnvConfig.SMTPPassword,
		SMTPHost:     env.EnvConfig.SMTPHost,
		SMTPPort:     env.EnvConfig.SMTPPort,
	}
}

//go:embed EmailTemplate/AccountSignUp.html
var signupTemplate string

type EmailData struct {
	Username string
	Token    string
}

func (e *emailService) SendHTMLEmail(subject string, htmlBody string, receiverEmail string) error {
	auth := smtp.PlainAuth("", e.SMTPEmail, e.SMTPPassword, e.SMTPHost)

	headers := map[string]string{
		"From":         e.SMTPEmail,
		"To":           receiverEmail,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	smtpAddr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
	err := smtp.SendMail(
		smtpAddr,
		auth,
		e.SMTPEmail,
		[]string{receiverEmail},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *emailService) CreateAccountSignUpEmail(username string, token string) (string, error) {
	tmpl, err := template.New("signup").Parse(signupTemplate)
	if err != nil {
		return "", fmt.Errorf("cannot read template: %v", err)
	}

	data := EmailData{
		Username: username,
		Token:    token,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("cannot execute template: %v", err)
	}

	return body.String(), nil
}
