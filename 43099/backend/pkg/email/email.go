package email

import (
	"fmt"
	"venue-booking/internal/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendVerificationEmail(to, code string) error {
	subject := "邮箱验证"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>邮箱验证</h2>
			<p>您好，</p>
			<p>您的验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
			<p>验证码有效期为10分钟，请尽快完成验证。</p>
			<p>如果这不是您的操作，请忽略此邮件。</p>
			<hr>
			<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿直接回复。</p>
		</div>
	`, code)

	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(to, code string) error {
	subject := "重置密码"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>重置密码</h2>
			<p>您好，</p>
			<p>您的重置密码验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
			<p>验证码有效期为10分钟，请尽快完成密码重置。</p>
			<p>如果这不是您的操作，请忽略此邮件。</p>
			<hr>
			<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿直接回复。</p>
		</div>
	`, code)

	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendOrderConfirmation(to, orderNo, itemName string) error {
	subject := "预约成功通知"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>预约成功</h2>
			<p>您好，</p>
			<p>您的预约已成功提交！</p>
			<p><strong>订单号：</strong>%s</p>
			<p><strong>预约项目：</strong>%s</p>
			<p>请等待管理员审核，审核通过后会收到通知。</p>
			<hr>
			<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿直接回复。</p>
		</div>
	`, orderNo, itemName)

	return s.sendEmail(to, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	if config.AppConfig.Email.Host == "" || config.AppConfig.Email.User == "" {
		fmt.Printf("Email not configured. Would send to: %s, Subject: %s\n", to, subject)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		config.AppConfig.Email.Host,
		config.AppConfig.Email.Port,
		config.AppConfig.Email.User,
		config.AppConfig.Email.Password,
	)

	return d.DialAndSend(m)
}
