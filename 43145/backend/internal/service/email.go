package service

import (
	"log"
	"survey-platform/internal/utils"

	"gopkg.in/gomail.v2"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	cfg := utils.GetConfig()

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EmailUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.EmailHost, cfg.EmailPort, cfg.EmailUser, cfg.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func (s *EmailService) SendEmailWithRetry(to, subject, body string, maxRetries int) error {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying email to %s (attempt %d/%d)", to, i+1, maxRetries+1)
		}
		err := s.SendEmail(to, subject, body)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return lastErr
}
