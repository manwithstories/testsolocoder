package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"smart-energy-platform/models"
	"time"
)

func SendNotification(userID uint, notifType, title, content string) {
	notification := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Content: content,
	}

	if err := models.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification: %v", err)
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err == nil {
		if user.Email != "" && Cfg != nil && Cfg.EmailSMTPHost != "" {
			go sendEmailNotification(user.Email, title, content)
		}
	}
}

func SendFamilyNotification(familyID uint, notifType, title, content string) {
	var members []models.FamilyMember
	models.DB.Where("family_id = ? AND status = ?", familyID, 1).Find(&members)

	for _, member := range members {
		SendNotification(member.UserID, notifType, title, content)
	}
}

func sendEmailNotification(to, subject, body string) {
	if Cfg == nil {
		return
	}

	smtpHost := Cfg.EmailSMTPHost
	smtpPort := Cfg.EmailSMTPPort
	from := Cfg.EmailFrom
	password := Cfg.EmailPassword

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), tlsconfig)
	if err != nil {
		log.Printf("Email TLS connection error: %v", err)
		return
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("Email client creation error: %v", err)
		return
	}
	defer client.Close()

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err := client.Auth(auth); err != nil {
		log.Printf("Email auth error: %v", err)
		return
	}

	if err := client.Mail(from); err != nil {
		log.Printf("Email MAIL error: %v", err)
		return
	}

	if err := client.Rcpt(to); err != nil {
		log.Printf("Email RCPT error: %v", err)
		return
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("Email DATA error: %v", err)
		return
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, body)

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Printf("Email write error: %v", err)
		return
	}

	err = w.Close()
	if err != nil {
		log.Printf("Email close error: %v", err)
		return
	}

	client.Quit()
	log.Printf("Email sent to: %s", to)
}

func CheckDeviceOffline() {
	var devices []models.Device
	models.DB.Where("status = ?", "online").Find(&devices)

	now := time.Now()
	for _, device := range devices {
		if device.LastOnlineTime != nil && now.Sub(*device.LastOnlineTime) > 30*time.Minute {
			device.Status = "offline"
			models.DB.Save(&device)
			UpdateDeviceStatus(device.ID, "offline")

			SendFamilyNotification(device.FamilyID, "alert",
				"Device Offline Alert",
				fmt.Sprintf("Device '%s' has gone offline.", device.Name))
		}
	}
}

func ScheduleOfflineChecker() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		CheckDeviceOffline()
	}
}
