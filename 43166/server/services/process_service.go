package services

import (
	"errors"
	"time"

	"business-registration-platform/database"
	"business-registration-platform/models"
)

type ProcessService struct{}

func NewProcessService() *ProcessService {
	return &ProcessService{}
}

type UpdateStepRequest struct {
	Status          string `json:"status"`
	Remark          string `json:"remark"`
	Description     string `json:"description"`
	CertificateFile string `json:"certificateFile"`
}

func (s *ProcessService) GetProcessSteps(applicationID uint) ([]models.ProcessStep, error) {
	var steps []models.ProcessStep
	if err := database.DB.Where("application_id = ?", applicationID).
		Order("step_order ASC").
		Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

func (s *ProcessService) UpdateProcessStep(stepID uint, handlerID uint, req *UpdateStepRequest) error {
	var step models.ProcessStep
	if err := database.DB.First(&step, stepID).Error; err != nil {
		return errors.New("process step not found")
	}

	if req.Status != "" {
		step.Status = models.ProcessStepStatus(req.Status)
	}
	if req.Remark != "" {
		step.Remark = req.Remark
	}
	if req.Description != "" {
		step.Description = req.Description
	}
	if req.CertificateFile != "" {
		step.CertificateFile = req.CertificateFile
	}

	handlerIDCopy := handlerID
	step.HandlerID = &handlerIDCopy
	now := time.Now()

	if step.Status == models.StepStatusCompleted {
		step.CompletedAt = &now
	} else if step.Status == models.StepStatusInProgress {
		if step.StartedAt == nil {
			step.StartedAt = &now
		}
	}

	if err := database.DB.Save(&step).Error; err != nil {
		return err
	}

	return s.updateApplicationProgress(step.ApplicationID)
}

func (s *ProcessService) updateApplicationProgress(applicationID uint) error {
	var steps []models.ProcessStep
	if err := database.DB.Where("application_id = ?", applicationID).Find(&steps).Error; err != nil {
		return err
	}

	if len(steps) == 0 {
		return nil
	}

	completedCount := 0
	var currentStep string
	for _, step := range steps {
		if step.Status == models.StepStatusCompleted {
			completedCount++
		}
		if step.Status == models.StepStatusInProgress {
			currentStep = string(step.StepType)
		}
	}

	progress := (completedCount * 100) / len(steps)

	updates := map[string]interface{}{
		"progress_percent": progress,
	}

	if currentStep != "" {
		updates["current_step"] = currentStep
	}

	if completedCount == len(steps) {
		updates["status"] = models.AppStatusCompleted
		now := time.Now()
		updates["completed_at"] = &now
	}

	return database.DB.Model(&models.Application{}).Where("id = ?", applicationID).Updates(updates).Error
}

func (s *ProcessService) StartStep(stepID uint, handlerID uint) error {
	var step models.ProcessStep
	if err := database.DB.First(&step, stepID).Error; err != nil {
		return errors.New("process step not found")
	}

	if step.Status != models.StepStatusPending {
		return errors.New("step is not in pending status")
	}

	now := time.Now()
	step.Status = models.StepStatusInProgress
	step.StartedAt = &now
	handlerIDCopy := handlerID
	step.HandlerID = &handlerIDCopy

	return database.DB.Save(&step).Error
}

func (s *ProcessService) CompleteStep(stepID uint, handlerID uint, certificateFile, remark string) error {
	var step models.ProcessStep
	if err := database.DB.First(&step, stepID).Error; err != nil {
		return errors.New("process step not found")
	}

	now := time.Now()
	step.Status = models.StepStatusCompleted
	step.CompletedAt = &now
	handlerIDCopy := handlerID
	step.HandlerID = &handlerIDCopy
	step.CertificateFile = certificateFile
	step.Remark = remark

	if err := database.DB.Save(&step).Error; err != nil {
		return err
	}

	return s.updateApplicationProgress(step.ApplicationID)
}

func (s *ProcessService) SkipStep(stepID uint, handlerID uint, remark string) error {
	var step models.ProcessStep
	if err := database.DB.First(&step, stepID).Error; err != nil {
		return errors.New("process step not found")
	}

	now := time.Now()
	step.Status = models.StepStatusSkipped
	step.CompletedAt = &now
	handlerIDCopy := handlerID
	step.HandlerID = &handlerIDCopy
	step.Remark = remark

	if err := database.DB.Save(&step).Error; err != nil {
		return err
	}

	return s.updateApplicationProgress(step.ApplicationID)
}
