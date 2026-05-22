package services

import (
	"errors"
	"podcast-manager/internal/database"
	"podcast-manager/internal/models"
	"podcast-manager/internal/rss"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PodcastService struct{}

func NewPodcastService() *PodcastService {
	return &PodcastService{}
}

func (s *PodcastService) AddPodcast(feedURL string) (*models.Podcast, error) {
	var existingPodcast models.Podcast
	result := database.DB.Where("feed_url = ?", feedURL).First(&existingPodcast)
	if result.Error == nil {
		return nil, errors.New("podcast already subscribed")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	podcast, episodes, err := rss.ParseFeed(feedURL)
	if err != nil {
		return nil, err
	}

	tx := database.BeginTransaction()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(podcast).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range episodes {
		episodes[i].PodcastID = podcast.ID
	}

	if len(episodes) > 0 {
		if err := tx.Create(&episodes).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return podcast, nil
}

func (s *PodcastService) GetPodcastList(page, perPage int, search string) ([]models.Podcast, int64, error) {
	var podcasts []models.Podcast
	var total int64

	query := database.DB.Model(&models.Podcast{})
	if search != "" {
		query = query.Where("title ILIKE ? OR author ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&podcasts).Error

	return podcasts, total, err
}

func (s *PodcastService) GetPodcastByID(id uuid.UUID) (*models.Podcast, error) {
	var podcast models.Podcast
	err := database.DB.Where("id = ?", id).First(&podcast).Error
	if err != nil {
		return nil, err
	}
	return &podcast, nil
}

func (s *PodcastService) UpdatePodcast(id uuid.UUID, updates map[string]interface{}) (*models.Podcast, error) {
	podcast, err := s.GetPodcastByID(id)
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(podcast).Updates(updates).Error; err != nil {
		return nil, err
	}

	return podcast, nil
}

func (s *PodcastService) DeletePodcast(id uuid.UUID) error {
	tx := database.BeginTransaction()
	if tx.Error != nil {
		return tx.Error
	}

	var episodes []models.Episode
	if err := tx.Where("podcast_id = ?", id).Find(&episodes).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, episode := range episodes {
		if err := tx.Where("episode_id = ?", episode.ID).Delete(&models.PlaybackProgress{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where("episode_id = ?", episode.ID).Delete(&models.Note{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where("episode_id = ?", episode.ID).Delete(&models.ListeningHistory{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where("episode_id = ?", episode.ID).Delete(&models.PlaylistItem{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("podcast_id = ?", id).Delete(&models.Episode{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.Podcast{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *PodcastService) RefreshPodcast(id uuid.UUID) (*models.Podcast, int, error) {
	podcast, err := s.GetPodcastByID(id)
	if err != nil {
		return nil, 0, err
	}

	parsedPodcast, episodes, err := rss.ParseFeed(podcast.FeedURL)
	if err != nil {
		return nil, 0, err
	}

	tx := database.BeginTransaction()
	if tx.Error != nil {
		return nil, 0, tx.Error
	}

	podcast.LastChecked = time.Now()
	if parsedPodcast.LastUpdated.After(podcast.LastUpdated) {
		podcast.LastUpdated = parsedPodcast.LastUpdated
	}
	if err := tx.Save(podcast).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	newEpisodesCount := 0
	for _, episode := range episodes {
		var existingEpisode models.Episode
		result := tx.Where("guid = ?", episode.GUID).First(&existingEpisode)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			episode.PodcastID = podcast.ID
			if err := tx.Create(&episode).Error; err != nil {
				tx.Rollback()
				return nil, 0, err
			}
			newEpisodesCount++
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return podcast, newEpisodesCount, nil
}

func (s *PodcastService) GetPodcastStats(id uuid.UUID) (*models.PodcastStats, error) {
	var stats models.PodcastStats
	stats.PodcastID = id

	var totalEpisodes int64
	database.DB.Model(&models.Episode{}).Where("podcast_id = ?", id).Count(&totalEpisodes)
	stats.TotalEpisodes = totalEpisodes

	var unplayedCount int64
	database.DB.Table("episodes").
		Where("podcast_id = ? AND id NOT IN (SELECT episode_id FROM playback_progresses WHERE completed = true)", id).
		Count(&unplayedCount)
	stats.UnplayedCount = unplayedCount

	var completedCount int64
	database.DB.Table("playback_progresses").
		Joins("JOIN episodes ON episodes.id = playback_progresses.episode_id").
		Where("episodes.podcast_id = ? AND playback_progresses.completed = true", id).
		Count(&completedCount)
	stats.CompletedCount = completedCount

	var totalListened int64
	database.DB.Table("listening_histories").
		Joins("JOIN episodes ON episodes.id = listening_histories.episode_id").
		Where("episodes.podcast_id = ?", id).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&totalListened)
	stats.TotalListened = totalListened

	return &stats, nil
}
