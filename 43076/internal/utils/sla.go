package utils

import (
	"time"

	"ticket-system/internal/config"
	"ticket-system/internal/models"
)

func CalculateSLADeadlines(priority string, createdAt time.Time) (responseDeadline, resolveDeadline time.Time) {
	slaConfig, ok := config.AppConfig.SLA.GetPriorityConfig(priority)
	if !ok {
		slaConfig, _ = config.AppConfig.SLA.GetPriorityConfig(models.TicketPriorityMedium)
	}

	responseDeadline = createdAt.Add(time.Duration(slaConfig.ResponseMinutes) * time.Minute)
	resolveDeadline = createdAt.Add(time.Duration(slaConfig.ResolutionHours) * time.Hour)

	return
}

func IsResponseSLABreached(createdAt, firstResponseAt time.Time, priority string) bool {
	slaConfig, ok := config.AppConfig.SLA.GetPriorityConfig(priority)
	if !ok {
		slaConfig, _ = config.AppConfig.SLA.GetPriorityConfig(models.TicketPriorityMedium)
	}

	responseDeadline := createdAt.Add(time.Duration(slaConfig.ResponseMinutes) * time.Minute)
	return firstResponseAt.After(responseDeadline)
}

func IsResolveSLABreached(createdAt, resolvedAt time.Time, priority string) bool {
	slaConfig, ok := config.AppConfig.SLA.GetPriorityConfig(priority)
	if !ok {
		slaConfig, _ = config.AppConfig.SLA.GetPriorityConfig(models.TicketPriorityMedium)
	}

	resolveDeadline := createdAt.Add(time.Duration(slaConfig.ResolutionHours) * time.Hour)
	return resolvedAt.After(resolveDeadline)
}

func ShouldEscalate(createdAt time.Time, priority string) bool {
	slaConfig, ok := config.AppConfig.SLA.GetPriorityConfig(priority)
	if !ok {
		slaConfig, _ = config.AppConfig.SLA.GetPriorityConfig(models.TicketPriorityMedium)
	}

	escalationTime := createdAt.Add(time.Duration(slaConfig.EscalationHours) * time.Hour)
	return time.Now().After(escalationTime)
}
