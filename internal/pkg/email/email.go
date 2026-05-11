package email

import (
	"fmt"
	"net/smtp"
	"devix-backend/internal/config"
)

type Mailer struct {
	config config.EmailConfig
}

func NewMailer(cfg config.EmailConfig) *Mailer {
	return &Mailer{config: cfg}
}

func (m *Mailer) Send(to, subject, body string) error {
	if m.config.SMTPHost == "" {
		return fmt.Errorf("SMTP host not configured")
	}

	auth := smtp.PlainAuth("", m.config.SMTPUser, m.config.SMTPPass, m.config.SMTPHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, subject, body))
	
	addr := fmt.Sprintf("%s:%d", m.config.SMTPHost, m.config.SMTPPort)
	return smtp.SendMail(addr, auth, m.config.SMTPFrom, []string{to}, msg)
}

func (m *Mailer) SendVerificationEmail(to, username, token string) error {
	subject := "Verify your Devix account"
	link := fmt.Sprintf("%s/verify-email?token=%s", m.config.FrontendURL, token)
	body := fmt.Sprintf(`
		<h1>Welcome to Devix, %s!</h1>
		<p>Please click the link below to verify your email address:</p>
		<a href="%s">Verify Email</a>
		<p>If you didn't create an account, you can ignore this email.</p>
	`, username, link)
	return m.Send(to, subject, body)
}

func (m *Mailer) SendPasswordResetEmail(to, token string) error {
	subject := "Reset your Devix password"
	link := fmt.Sprintf("%s/reset-password?token=%s", m.config.FrontendURL, token)
	body := fmt.Sprintf(`
		<h1>Password Reset Request</h1>
		<p>You requested a password reset. Click the link below to set a new password:</p>
		<a href="%s">Reset Password</a>
		<p>This link will expire in 1 hour. If you didn't request this, please secure your account.</p>
	`, link)
	return m.Send(to, subject, body)
}
