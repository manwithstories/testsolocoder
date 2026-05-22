package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"travel-planner/internal/config"
	"travel-planner/internal/logger"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	if !s.config.Enabled {
		logger.Infof("Email sending is disabled. Would send to %s: %s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	from := s.config.From
	if from == "" {
		from = s.config.Username
	}

	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	var err error
	if s.config.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: s.config.SkipVerify,
			ServerName:         s.config.Host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.config.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Quit()

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}

		if err := client.Mail(from); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		for _, recipient := range strings.Split(to, ",") {
			if err := client.Rcpt(strings.TrimSpace(recipient)); err != nil {
				return fmt.Errorf("failed to set recipient: %w", err)
			}
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}

		if _, err := w.Write([]byte(msg)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}
	} else {
		err = smtp.SendMail(addr, auth, from, strings.Split(to, ","), []byte(msg))
	}

	if err != nil {
		logger.Errorf("Failed to send email to %s: %v", to, err)
		return err
	}

	logger.Infof("Email sent successfully to %s", to)
	return nil
}

func (s *EmailService) SendReminderEmail(to, reminderTitle, reminderDesc, planTitle string) error {
	subject := fmt.Sprintf("🔔 旅行提醒: %s", reminderTitle)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 10px 10px; }
        .reminder-title { font-size: 20px; font-weight: bold; margin-bottom: 10px; }
        .plan-info { background: #e9ecef; padding: 10px; border-radius: 5px; margin: 15px 0; }
        .footer { text-align: center; color: #999; font-size: 12px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🔔 旅行提醒</h2>
        </div>
        <div class="content">
            <div class="reminder-title">%s</div>
            <p>%s</p>
            <div class="plan-info">
                <strong>关联计划:</strong> %s
            </div>
            <p>请登录旅行规划系统查看详情。</p>
        </div>
        <div class="footer">
            <p>此邮件由旅行规划系统自动发送，请勿直接回复。</p>
        </div>
    </div>
</body>
</html>
`, reminderTitle, reminderDesc, planTitle)

	return s.SendEmail(to, subject, body)
}
