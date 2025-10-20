package notifier

import (
	"crypto-alerts/internal/config"
	"fmt"
	"net/smtp"
)

type Notifier interface {
	SendEmailAlert(to string, subject string, message string) error
}

type emailNotifier struct {
	smtpConfig *config.SMTPConfig
}

func NewEmailNotifier(cfg *config.SMTPConfig) Notifier {
	return &emailNotifier{
		smtpConfig: cfg,
	}
}

func (e *emailNotifier) SendEmailAlert(to string, subject string, message string) error {
	if e.smtpConfig == nil {
		return fmt.Errorf("SMTP config not initialized")
	}
	smtpAddr := fmt.Sprintf("%s:%d", e.smtpConfig.Host, e.smtpConfig.Port)

	headers := make(map[string]string)
	headers["From"] = e.smtpConfig.Username
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	body := ""
	for key, value := range headers {
		body += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	body += "\r\n" + message

	auth := smtp.PlainAuth("", e.smtpConfig.Username, e.smtpConfig.Password, e.smtpConfig.Host)

	err := smtp.SendMail(smtpAddr, auth, e.smtpConfig.Username, []string{to}, []byte(body))
	if err != nil {
		return fmt.Errorf("erro ao enviar e-mail de alerta: %w", err)
	}

	return nil
}
