package scheduler

import (
	"log"
	"time"

	"ticket-system/internal/database"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var Scheduler *cron.Cron

func InitScheduler() {
	Scheduler = cron.New(cron.WithSeconds())

	Scheduler.AddFunc("0 */5 * * * *", CheckSLABreach)

	Scheduler.AddFunc("0 */10 * * * *", CheckSLAAboutToExpire)

	Scheduler.Start()

	log.Println("Scheduler started: SLA checker running every 5 minutes")
}

func StopScheduler() {
	if Scheduler != nil {
		ctx := Scheduler.Stop()
		<-ctx.Done()
		log.Println("Scheduler stopped")
	}
}

func CheckSLABreach() {
	log.Println("[SLA Check] Starting SLA breach check...")

	now := time.Now()
	var tickets []models.Ticket

	err := database.DB.Preload("Customer").Preload("Assignee").
		Where("status IN ?", []string{
			models.TicketStatusOpen,
			models.TicketStatusAssigned,
			models.TicketStatusInProgress,
			models.TicketStatusPending,
		}).
		Where("is_escalated = ?", false).
		Where("(response_deadline < ? AND first_response_at IS NULL) OR (resolve_deadline < ?)", now, now).
		Find(&tickets).Error

	if err != nil {
		log.Printf("[SLA Check] Error querying breached tickets: %v", err)
		return
	}

	log.Printf("[SLA Check] Found %d breached tickets", len(tickets))

	for _, ticket := range tickets {
		if err := escalateTicket(&ticket); err != nil {
			log.Printf("[SLA Check] Failed to escalate ticket %s: %v", ticket.TicketNo, err)
		} else {
			log.Printf("[SLA Check] Ticket %s escalated successfully", ticket.TicketNo)
		}
	}

	log.Println("[SLA Check] SLA breach check completed")
}

func CheckSLAAboutToExpire() {
	log.Println("[SLA Check] Starting SLA about-to-expire check...")

	now := time.Now()
	warningThreshold := 30 * time.Minute
	warningTime := now.Add(warningThreshold)

	var tickets []models.Ticket

	err := database.DB.Preload("Customer").Preload("Assignee").
		Where("status IN ?", []string{
			models.TicketStatusOpen,
			models.TicketStatusAssigned,
			models.TicketStatusInProgress,
			models.TicketStatusPending,
		}).
		Where("is_escalated = ?", false).
		Where("first_response_at IS NULL").
		Where("response_deadline > ? AND response_deadline < ?", now, warningTime).
		Find(&tickets).Error

	if err != nil {
		log.Printf("[SLA Check] Error querying about-to-expire tickets: %v", err)
		return
	}

	log.Printf("[SLA Check] Found %d tickets about to expire", len(tickets))

	for _, ticket := range tickets {
		minutesLeft := int(time.Until(*ticket.ResponseDeadline).Minutes())
		if ticket.Assignee != nil {
			customerName := ""
			if ticket.Customer != nil {
				customerName = ticket.Customer.Name
			}
			if err := utils.SendSLAAboutToExpireNotification(
				ticket.TicketNo,
				ticket.Title,
				ticket.Priority,
				customerName,
				ticket.Assignee.Email,
				minutesLeft,
			); err != nil {
				log.Printf("[SLA Check] Failed to send warning email for ticket %s: %v", ticket.TicketNo, err)
			} else {
				log.Printf("[SLA Check] Warning email sent for ticket %s (%d min left)", ticket.TicketNo, minutesLeft)
			}
		}
	}

	log.Println("[SLA Check] SLA about-to-expire check completed")
}

func escalateTicket(ticket *models.Ticket) error {
	return database.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		updates := map[string]interface{}{
			"is_escalated": true,
			"escalated_at": &now,
			"status":       models.TicketStatusEscalated,
		}

		if err := tx.Model(ticket).Updates(updates).Error; err != nil {
			return err
		}

		var slaRecord models.SLARecord
		if err := tx.Where("ticket_id = ?", ticket.ID).First(&slaRecord).Error; err == nil {
			tx.Model(&slaRecord).Update("escalation_count", gorm.Expr("escalation_count + ?", 1))
		}

		logEntry := &models.TicketLog{
			TicketID:   ticket.ID,
			OperatorID: 1,
			Action:     "auto_escalate",
			OldValue:   ticket.Status,
			NewValue:   models.TicketStatusEscalated,
			Remark:     "SLA超时自动升级",
		}
		if err := tx.Create(logEntry).Error; err != nil {
			return err
		}

		if ticket.Assignee != nil {
			customerName := ""
			if ticket.Customer != nil {
				customerName = ticket.Customer.Name
			}

			dutyEmails := getDutyPersonEmails(tx)

			if err := utils.SendSLAEscalationNotification(
				ticket.TicketNo,
				ticket.Title,
				ticket.Priority,
				customerName,
				ticket.Assignee.Email,
				dutyEmails,
			); err != nil {
				log.Printf("[SLA Check] Warning: Failed to send escalation email: %v", err)
			}
		}

		return nil
	})
}

func getDutyPersonEmails(tx *gorm.DB) []string {
	var users []models.User
	err := tx.Where("is_on_duty = ? AND role IN ?", true, []string{models.RoleAdmin, models.RoleManager}).
		Select("email").
		Find(&users).Error

	if err != nil {
		log.Printf("[SLA Check] Failed to get duty persons: %v", err)
		return []string{}
	}

	emails := make([]string, 0, len(users))
	for _, u := range users {
		if u.Email != "" {
			emails = append(emails, u.Email)
		}
	}

	return emails
}
