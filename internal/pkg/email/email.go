package email

import (
	"crypto/tls"
	"fmt"
	"html"
	"net"
	"net/smtp"
	"time"

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

	addr := net.JoinHostPort(m.config.SMTPHost, fmt.Sprintf("%d", m.config.SMTPPort))
	auth := smtp.PlainAuth("", m.config.SMTPUser, m.config.SMTPPass, m.config.SMTPHost)

	// Build headers
	headers := make(map[string]string)
	headers["From"] = m.config.SMTPFrom
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Format message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to SMTP server
	var err error
	var conn net.Conn
	tlsConfig := &tls.Config{
		ServerName: m.config.SMTPHost,
	}

	if m.config.SMTPPort == 465 {
		conn, err = tls.Dial("tcp", addr, tlsConfig)
	} else {
		conn, err = net.DialTimeout("tcp", addr, 5*time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// If not port 465, attempt STARTTLS if supported
	if m.config.SMTPPort != 465 {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err = client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("failed to upgrade connection with STARTTLS: %w", err)
			}
		}
	}

	// Authenticate
	if m.config.SMTPUser != "" && m.config.SMTPPass != "" {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// Send envelope
	if err = client.Mail(m.config.SMTPFrom); err != nil {
		return fmt.Errorf("SMTP Mail command failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP Rcpt command failed: %w", err)
	}

	// Send message body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP Data command failed: %w", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (m *Mailer) SendVerificationEmail(to, username, token string) error {
	subject := "Verify your Devix account"
	link := fmt.Sprintf("%s/verify-email?token=%s", m.config.FrontendURL, token)
	escapedUsername := html.EscapeString(username)
	body := fmt.Sprintf(`
		<h1>Welcome to Devix, %s!</h1>
		<p>Please click the link below to verify your email address:</p>
		<a href="%s">Verify Email</a>
		<p>If you didn't create an account, you can ignore this email.</p>
	`, escapedUsername, link)
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
