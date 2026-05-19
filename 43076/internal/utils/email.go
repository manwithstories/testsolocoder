package utils

import (
	"fmt"
	"log"
	"ticket-system/internal/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(to []string, subject, body string) error {
	emailCfg := config.AppConfig.Email
	if !emailCfg.Enabled {
		log.Println("Email is disabled in config, skipping email notification")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailCfg.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		emailCfg.Host,
		emailCfg.Port,
		emailCfg.Username,
		emailCfg.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to: %v, subject: %s", to, subject)
	return nil
}

func SendSLAEscalationNotification(ticketNo, title, priority, customerName string, assigneeEmail string, dutyEmails []string) error {
	subject := fmt.Sprintf("[SLA升级] 工单 %s - %s", ticketNo, title)

	body := fmt.Sprintf(`
SLA超时升级通知
========================

工单号: %s
标题: %s
优先级: %s
客户: %s
处理人: %s

该工单已超过SLA响应或解决时限，请尽快处理！

请登录工单系统查看详情。

--
自动发送，请勿回复
`, ticketNo, title, priority, customerName, assigneeEmail)

	recipients := []string{assigneeEmail}
	for _, email := range dutyEmails {
		if email != assigneeEmail {
			recipients = append(recipients, email)
		}
	}

	return SendEmail(recipients, subject, body)
}

func SendSLAAboutToExpireNotification(ticketNo, title, priority, customerName string, assigneeEmail string, minutesLeft int) error {
	subject := fmt.Sprintf("[SLA预警] 工单 %s 即将超时 (%d分钟)", ticketNo, minutesLeft)

	body := fmt.Sprintf(`
SLA即将超时预警
========================

工单号: %s
标题: %s
优先级: %s
客户: %s
处理人: %s

该工单将在 %d 分钟内超时，请尽快处理！

请登录工单系统查看详情。

--
自动发送，请勿回复
`, ticketNo, title, priority, customerName, assigneeEmail, minutesLeft)

	return SendEmail([]string{assigneeEmail}, subject, body)
}

func SendTicketAssignmentNotification(ticketNo, title, priority, customerName, assigneeEmail, assignerName string) error {
	subject := fmt.Sprintf("[工单分配] 新工单 %s - %s", ticketNo, title)

	body := fmt.Sprintf(`
新工单分配通知
========================

工单号: %s
标题: %s
优先级: %s
客户: %s
分配人: %s

您有一个新的工单需要处理，请及时登录系统查看。

--
自动发送，请勿回复
`, ticketNo, title, priority, customerName, assignerName)

	return SendEmail([]string{assigneeEmail}, subject, body)
}
