package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"freelancer-management/internal/config"
	"freelancer-management/internal/logger"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{config: cfg}
}

func (s *EmailService) SendWelcomeEmail(to, firstName string) error {
	data := struct {
		FirstName string
	}{
		FirstName: firstName,
	}

	return s.sendEmail(to, "Welcome to Freelancer Management", "welcome", data)
}

func (s *EmailService) SendInvoiceEmail(to, invoiceNumber, clientName string, amount float64) error {
	data := struct {
		InvoiceNumber string
		ClientName    string
		Amount        float64
	}{
		InvoiceNumber: invoiceNumber,
		ClientName:    clientName,
		Amount:        amount,
	}

	return s.sendEmail(to, fmt.Sprintf("Invoice %s is ready", invoiceNumber), "invoice", data)
}

func (s *EmailService) sendEmail(to, subject, templateName string, data interface{}) error {
	body, err := s.renderTemplate(templateName, data)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		s.config.FromName, s.config.FromAddress, to, subject, body)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.FromAddress,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		logger.LogError("Failed to send email: %v", err)
		return err
	}

	logger.LogInfo("Email sent successfully to %s", to)
	return nil
}

func (s *EmailService) renderTemplate(name string, data interface{}) (string, error) {
	templates := map[string]string{
		"welcome": `
			<h1>Welcome, {{.FirstName}}!</h1>
			<p>Thank you for joining Freelancer Management. Start managing your projects and invoices today.</p>
		`,
		"invoice": `
			<h1>Invoice {{.InvoiceNumber}}</h1>
			<p>Dear {{.ClientName}},</p>
			<p>Your invoice is ready. Total amount: ${{printf "%.2f" .Amount}}</p>
		`,
	}

	tmpl, err := template.New(name).Parse(templates[name])
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
