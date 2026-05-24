package utils

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

type EmailService struct {
	config EmailConfig
	dialer *gomail.Dialer
}

func NewEmailService(config EmailConfig) *EmailService {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.Username, config.Password)
	return &EmailService{
		config: config,
		dialer: dialer,
	}
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := es.dialer.DialAndSend(m); err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return nil
	}

	return fmt.Errorf("邮件发送失败，已重试%d次: %w", maxRetries, lastErr)
}

func (es *EmailService) SendInterviewInvitation(to, applicantName, jobTitle, interviewer string, scheduledAt time.Time, location string) error {
	subject := fmt.Sprintf("面试邀请 - %s职位", jobTitle)
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">面试邀请</h2>
			<p>尊敬的 %s：</p>
			<p>感谢您投递 <strong>%s</strong> 职位。我们非常高兴地邀请您参加面试。</p>
			<div style="background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 15px 0;">
				<p><strong>面试时间：</strong>%s</p>
				<p><strong>面试官：</strong>%s</p>
				<p><strong>面试地点：</strong>%s</p>
			</div>
			<p>请准时参加，如有变动请及时联系我们。</p>
			<p style="color: #666; font-size: 12px;">此邮件由招聘系统自动发送，请勿回复。</p>
		</div>
	`, applicantName, jobTitle, scheduledAt.Format("2006-01-02 15:04"), interviewer, location)

	return es.SendEmail(to, subject, body)
}

func (es *EmailService) SendApplicationStatusUpdate(to, applicantName, jobTitle, status string) error {
	subject := fmt.Sprintf("投递状态更新 - %s", jobTitle)
	statusText := map[string]string{
		"viewed":     "已查看",
		"interested": "HR感兴趣",
		"interview":  "面试中",
		"accepted":   "已录用",
		"rejected":   "未通过",
	}[status]

	if statusText == "" {
		statusText = status
	}

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">投递状态更新</h2>
			<p>尊敬的 %s：</p>
			<p>您投递的 <strong>%s</strong> 职位状态已更新为：<strong>%s</strong></p>
			<p>如有疑问，请联系我们。</p>
			<p style="color: #666; font-size: 12px;">此邮件由招聘系统自动发送，请勿回复。</p>
		</div>
	`, applicantName, jobTitle, statusText)

	return es.SendEmail(to, subject, body)
}

func (es *EmailService) SendWelcomeEmail(to, name string) error {
	subject := "欢迎加入招聘平台"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">欢迎加入！</h2>
			<p>尊敬的 %s：</p>
			<p>欢迎注册招聘平台！您的账户已成功创建。</p>
			<p>您可以登录平台开始您的求职/招聘之旅。</p>
			<p style="color: #666; font-size: 12px;">此邮件由招聘系统自动发送，请勿回复。</p>
		</div>
	`, name)

	return es.SendEmail(to, subject, body)
}
