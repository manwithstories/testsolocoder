package services

import (
	"errors"
	"podcast-manager/internal/database"
	"podcast-manager/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeService struct{}

func NewEpisodeService() *EpisodeService {
	return &EpisodeService{}
}

func (s *EpisodeService) GetEpisodeList(podcastID uuid.UUID, page, perPage int, search string) ([]models.Episode, int64, error) {
	var episodes []models.Episode
	var total int64

	query := database.DB.Model(&models.Episode{})
	if podcastID != uuid.Nil {
		query = query.Where("podcast_id = ?", podcastID)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("pub_date DESC").
		Offset(offset).
		Limit(perPage).
		Preload("Podcast").
		Find(&episodes).Error

	return episodes, total, err
}

func (s *EpisodeService) GetEpisodeByID(id uuid.UUID) (*models.Episode, error) {
	var episode models.Episode
	err := database.DB.Preload("Podcast").Where("id = ?", id).First(&episode).Error
	if err != nil {
		return nil, err
	}

	if episode.IsNew {
		database.DB.Model(&episode).Update("is_new", false)
	}

	return &episode, nil
}

type PlaybackService struct{}

func NewPlaybackService() *PlaybackService {
	return &PlaybackService{}
}

func (s *PlaybackService) UpdatePlaybackProgress(episodeID uuid.UUID, currentTime float64, duration int) (*models.PlaybackProgress, error) {
	tx := database.BeginTransaction()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var progress models.PlaybackProgress
	result := tx.Where("episode_id = ?", episodeID).First(&progress)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, result.Error
	}

	isNew := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if isNew {
		progress = models.PlaybackProgress{
			EpisodeID:    episodeID,
			CurrentTime:  currentTime,
			LastPlayedAt: time.Now(),
			PlayCount:    1,
		}
		if duration > 0 && currentTime >= float64(duration)*0.9 {
			progress.Completed = true
			progress.CompletedAt = time.Now()
		}
		if err := tx.Create(&progress).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		updates := map[string]interface{}{
			"current_time":   currentTime,
			"last_played_at": time.Now(),
		}
		if duration > 0 && currentTime >= float64(duration)*0.9 && !progress.Completed {
			updates["completed"] = true
			updates["completed_at"] = time.Now()
		}
		if progress.LastPlayedAt.IsZero() || time.Since(progress.LastPlayedAt) > time.Hour {
			updates["play_count"] = progress.PlayCount + 1
		}
		if err := tx.Model(&progress).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	database.DB.First(&progress, progress.ID)
	return &progress, nil
}

func (s *PlaybackService) GetPlaybackProgress(episodeID uuid.UUID) (*models.PlaybackProgress, error) {
	var progress models.PlaybackProgress
	err := database.DB.Where("episode_id = ?", episodeID).First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.PlaybackProgress{
				EpisodeID:   episodeID,
				CurrentTime: 0,
				Completed: false,
				PlayCount: 0,
			}, nil
		}
		return nil, err
	}
	return &progress, nil
}

func (s *PlaybackService) MarkAsCompleted(episodeID uuid.UUID) error {
	return database.DB.Model(&models.PlaybackProgress{}).
		Where("episode_id = ?", episodeID).
		Updates(map[string]interface{}{
			"completed":    true,
			"completed_at": time.Now(),
		}).Error
}

func (s *HistoryService) AddListeningHistory(episodeID uuid.UUID, startTime, endTime time.Time, duration int, completion float64) (*models.ListeningHistory, error) {
	history := &models.ListeningHistory{
		EpisodeID: episodeID,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
		Completion: completion,
	}
	err := database.DB.Create(history).Error
	return history, err
}

type HistoryService struct{}

func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

func (s *HistoryService) GetListeningHistory(page, perPage int, podcastID uuid.UUID, startDate, endDate string, completed *bool) ([]models.ListeningHistory, int64, error) {
	var histories []models.ListeningHistory
	var total int64

	query := database.DB.Model(&models.ListeningHistory{})
	if podcastID != uuid.Nil {
		query = query.Joins("JOIN episodes ON episodes.id = listening_histories.episode_id").
			Where("episodes.podcast_id = ?", podcastID)
	}
	if startDate != "" {
		query = query.Where("DATE(listening_histories.start_time) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(listening_histories.start_time) <= ?", endDate)
	}
	if completed != nil {
		if *completed {
			query = query.Where("completion >= 0.9")
		} else {
			query = query.Where("completion < 0.9")
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("start_time DESC").
		Offset(offset).
		Limit(perPage).
		Preload("Episode.Podcast").
		Find(&histories).Error

	return histories, total, err
}

type NoteService struct{}

func NewNoteService() *NoteService {
	return &NoteService{}
}

func (s *NoteService) AddNote(episodeID uuid.UUID, timestamp float64, content string, tags []string) (*models.Note, error) {
	note := &models.Note{
		EpisodeID: episodeID,
		Timestamp:  timestamp,
		Content:    content,
		Tags:       tags,
	}
	err := database.DB.Create(note).Error
	return note, err
}

func (s *NoteService) GetNotes(episodeID uuid.UUID, search string, tag string) ([]models.Note, error) {
	var notes []models.Note
	query := database.DB.Where("episode_id = ?", episodeID)
	if search != "" {
		query = query.Where("content ILIKE ?", "%"+search+"%")
	}
	if tag != "" {
		query = query.Where("tags @> ARRAY[?]::text[]", tag)
	}
	err := query.Order("timestamp ASC").Find(&notes).Error
	return notes, err
}

func (s *NoteService) UpdateNote(id uuid.UUID, content string, tags []string) (*models.Note, error) {
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"content": content,
		"tags":    tags,
	}
	if err := database.DB.Model(&note).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *NoteService) DeleteNote(id uuid.UUID) error {
	return database.DB.Delete(&models.Note{}, id).Error
}

func (s *NoteService) SearchNotes(search string) ([]models.Note, error) {
	var notes []models.Note
	err := database.DB.Where("content ILIKE ?", "%"+search+"%").
		Preload("Episode.Podcast").
		Order("created_at DESC").
		Find(&notes).Error
	return notes, err
}
