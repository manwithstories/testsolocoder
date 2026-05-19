package utils

import (
	"errors"
	"fmt"
	"time"

	"ticket-system/internal/database"
	"ticket-system/internal/models"

	"gorm.io/gorm"
)

func GenerateTicketNo() (string, error) {
	date := time.Now().Format("20060102")
	var ticketNo string

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		var counter models.TicketCounter
		err := database.DB.Where("date = ?", date).First(&counter).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		var newCount int
		var version int

		if errors.Is(err, gorm.ErrRecordNotFound) {
			counter = models.TicketCounter{
				Date:    date,
				Count:   1,
				Version: 1,
			}
			if err := database.DB.Create(&counter).Error; err != nil {
				if isDuplicateEntryError(err) {
					continue
				}
				return "", err
			}
			newCount = 1
		} else {
			newCount = counter.Count + 1
			version = counter.Version

			result := database.DB.Model(&models.TicketCounter{}).
				Where("id = ? AND version = ?", counter.ID, version).
				Updates(map[string]interface{}{
					"count":   newCount,
					"version": version + 1,
				})

			if result.Error != nil {
				return "", result.Error
			}

			if result.RowsAffected == 0 {
				continue
			}
		}

		ticketNo = fmt.Sprintf("TK%s%06d", date, newCount)
		return ticketNo, nil
	}

	return "", errors.New("failed to generate ticket number after max retries")
}

func isDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return len(errStr) > 10 && errStr[10:15] == "1062"
}
