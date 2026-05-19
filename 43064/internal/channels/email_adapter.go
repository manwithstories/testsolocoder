package channels

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"go.uber.org/zap"
)

type EmailAdapter struct{}

func NewEmailAdapter() *EmailAdapter {
	return &EmailAdapter{}
}

func (e *EmailAdapter) Send(message *models.Message, config string) error {
	var emailConfig models.EmailConfig
	if err := json.Unmarshal([]byte(config), &emailConfig); err != nil {
		return fmt.Errorf("invalid email config: %w", err)
	}

	from := mail.Address{Name: emailConfig.FromName, Address: emailConfig.FromAddress}
	to := mail.Address{Address: message.Recipient}

	subject := message.Subject
	if subject == "" {
		subject = "Notification"
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var body string
	for k, v := range headers {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body += "\r\n" + message.Content

	addr := fmt.Sprintf("%s:%d", emailConfig.SMTPHost, emailConfig.SMTPPort)

	var auth smtp.Auth
	if emailConfig.Username != "" && emailConfig.Password != "" {
		auth = smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTPHost)
	}

	if emailConfig.UseSSL {
		return e.sendSSL(addr, auth, emailConfig.FromAddress, []string{message.Recipient}, []byte(body))
	}

	if emailConfig.UseTLS {
		return e.sendTLS(addr, auth, emailConfig.SMTPHost, emailConfig.FromAddress, []string{message.Recipient}, []byte(body))
	}

	return smtp.SendMail(addr, auth, emailConfig.FromAddress, []string{message.Recipient}, []byte(body))
}

func (e *EmailAdapter) sendSSL(addr string, auth smtp.Auth, from string, to []string, body []byte) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, addr[:strings.Index(addr, ":")])
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	return w.Close()
}

func (e *EmailAdapter) sendTLS(addr string, auth smtp.Auth, host string, from string, to []string, body []byte) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.StartTLS(tlsConfig); err != nil {
		return err
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	return w.Close()
}

func (e *EmailAdapter) TestConnection(config string) error {
	var emailConfig models.EmailConfig
	if err := json.Unmarshal([]byte(config), &emailConfig); err != nil {
		return fmt.Errorf("invalid email config: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", emailConfig.SMTPHost, emailConfig.SMTPPort)

	var auth smtp.Auth
	if emailConfig.Username != "" && emailConfig.Password != "" {
		auth = smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTPHost)
	}

	if emailConfig.UseSSL {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("ssl connection failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, emailConfig.SMTPHost)
		if err != nil {
			return fmt.Errorf("smtp client creation failed: %w", err)
		}
		defer client.Close()

		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}

		if err := client.Noop(); err != nil {
			return fmt.Errorf("smtp noop failed: %w", err)
		}
	} else if emailConfig.UseTLS {
		conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			return fmt.Errorf("tcp connection failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, emailConfig.SMTPHost)
		if err != nil {
			return fmt.Errorf("smtp client creation failed: %w", err)
		}
		defer client.Close()

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         emailConfig.SMTPHost,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}

		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}

		if err := client.Noop(); err != nil {
			return fmt.Errorf("smtp noop failed: %w", err)
		}
	} else {
		conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			return fmt.Errorf("tcp connection failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, emailConfig.SMTPHost)
		if err != nil {
			return fmt.Errorf("smtp client creation failed: %w", err)
		}
		defer client.Close()

		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}

		if err := client.Noop(); err != nil {
			return fmt.Errorf("smtp noop failed: %w", err)
		}
	}

	logger.Info("email connection test successful", zap.String("host", emailConfig.SMTPHost), zap.Int("port", emailConfig.SMTPPort))
	return nil
}
