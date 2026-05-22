package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
	"auction-system/config"
	"auction-system/pkg/logger"
)

func SendMail(to, subject, body string) error {
	if config.AppConfig.Mail.User == "" || config.AppConfig.Mail.Password == "" {
		logger.Warn("Mail configuration not set, skipping email to: %s", to)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.Mail.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		config.AppConfig.Mail.Host,
		config.AppConfig.Mail.Port,
		config.AppConfig.Mail.User,
		config.AppConfig.Mail.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		logger.Error("Failed to send email: %v", err)
		return err
	}

	logger.Info("Email sent successfully to: %s", to)
	return nil
}

func SendBidOutbidNotification(userEmail, auctionName, currentPrice string) {
	subject := fmt.Sprintf("您在【%s】的出价被超越了", auctionName)
	body := fmt.Sprintf(`
		<h3>出价被超越通知</h3>
		<p>您好，您在拍卖品【%s】的出价已被其他用户超越。</p>
		<p>当前最高价格：<strong>¥%s</strong></p>
		<p>请尽快登录系统重新出价！</p>
	`, auctionName, currentPrice)
	SendMail(userEmail, subject, body)
}

func SendAuctionEndingSoonNotification(userEmail, auctionName, timeLeft string) {
	subject := fmt.Sprintf("拍卖即将结束：%s", auctionName)
	body := fmt.Sprintf(`
		<h3>拍卖即将结束提醒</h3>
		<p>您好，您关注的拍卖品【%s】将在 %s 后结束。</p>
		<p>请抓紧时间出价！</p>
	`, auctionName, timeLeft)
	SendMail(userEmail, subject, body)
}

func SendBidSuccessNotification(userEmail, auctionName, price string) {
	subject := fmt.Sprintf("恭喜您成功拍得【%s】", auctionName)
	body := fmt.Sprintf(`
		<h3>竞拍成功通知</h3>
		<p>恭喜您！您已成功拍得拍卖品【%s】。</p>
		<p>成交价格：<strong>¥%s</strong></p>
		<p>请尽快完成支付！</p>
	`, auctionName, price)
	SendMail(userEmail, subject, body)
}

func SendPaymentSuccessNotification(userEmail, auctionName, orderNo string) {
	subject := fmt.Sprintf("支付成功：%s", auctionName)
	body := fmt.Sprintf(`
		<h3>支付成功通知</h3>
		<p>您好，您的订单【%s】已支付成功。</p>
		<p>拍卖品：%s</p>
		<p>请等待卖家发货！</p>
	`, orderNo, auctionName)
	SendMail(userEmail, subject, body)
}
