package services

import (
	"podcast-manager/internal/database"
	"podcast-manager/internal/models"

	"github.com/google/uuid"
)

type PlaylistService struct{}

func NewPlaylistService() *PlaylistService {
	return &PlaylistService{}
}

func (s *PlaylistService) CreatePlaylist(name, description, coverImage string) (*models.Playlist, error) {
	playlist := &models.Playlist{
		Name:        name,
		Description: description,
		CoverImage:  coverImage,
	}
	err := database.DB.Create(playlist).Error
	return playlist, err
}

func (s *PlaylistService) GetPlaylistList() ([]models.Playlist, error) {
	var playlists []models.Playlist
	err := database.DB.Order("created_at DESC").Find(&playlists).Error
	return playlists, err
}

func (s *PlaylistService) GetPlaylistByID(id uuid.UUID) (*models.Playlist, error) {
	var playlist models.Playlist
	err := database.DB.Preload("Items.Episode.Podcast").Where("id = ?", id).First(&playlist).Error
	return &playlist, err
}

func (s *PlaylistService) UpdatePlaylist(id uuid.UUID, name, description, coverImage string) (*models.Playlist, error) {
	var playlist models.Playlist
	if err := database.DB.First(&playlist, id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"name":         name,
		"description":  description,
		"cover_image":  coverImage,
	}
	if err := database.DB.Model(&playlist).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (s *PlaylistService) DeletePlaylist(id uuid.UUID) error {
	tx := database.BeginTransaction()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("playlist_id = ?", id).Delete(&models.PlaylistItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.Playlist{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *PlaylistService) AddEpisodeToPlaylist(playlistID, episodeID uuid.UUID) (*models.PlaylistItem, error) {
	var maxPosition int
	database.DB.Model(&models.PlaylistItem{}).Where("playlist_id = ?", playlistID).Select("COALESCE(MAX(position), -1)").Scan(&maxPosition)

	item := &models.PlaylistItem{
		PlaylistID: playlistID,
		EpisodeID:  episodeID,
		Position:   maxPosition + 1,
	}
	err := database.DB.Create(item).Error
	return item, err
}

func (s *PlaylistService) RemoveEpisodeFromPlaylist(playlistID, itemID uuid.UUID) error {
	return database.DB.Where("id = ? AND playlist_id = ?", itemID, playlistID).Delete(&models.PlaylistItem{}).Error
}

func (s *PlaylistService) ReorderPlaylistItems(playlistID uuid.UUID, itemIDs []uuid.UUID) error {
	tx := database.BeginTransaction()
	if tx.Error != nil {
		return tx.Error
	}

	for i, itemID := range itemIDs {
		if err := tx.Model(&models.PlaylistItem{}).
			Where("id = ? AND playlist_id = ?", itemID, playlistID).
			Update("position", i).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

func (s *StatsService) GetListeningStats(days int) ([]models.ListeningStats, error) {
	var stats []models.ListeningStats
	err := database.DB.Table("listening_histories").
		Select("DATE(start_time) as date, COALESCE(SUM(duration), 0) as total_duration, COUNT(DISTINCT episode_id) as episode_count, AVG(completion) as completion_rate").
		Where("start_time >= NOW() - INTERVAL '?' days", days).
		Group("DATE(start_time)").
		Order("date DESC").
		Scan(&stats).Error
	return stats, err
}

func (s *StatsService) GetPodcastDistribution() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := database.DB.Table("listening_histories").
		Select("podcasts.id, podcasts.title, podcasts.cover_image, COALESCE(SUM(listening_histories.duration), 0) as total_duration, COUNT(DISTINCT listening_histories.episode_id) as episode_count").
		Joins("JOIN episodes ON episodes.id = listening_histories.episode_id").
		Joins("JOIN podcasts ON podcasts.id = episodes.podcast_id").
		Group("podcasts.id, podcasts.title, podcasts.cover_image").
		Order("total_duration DESC").
		Scan(&results).Error
	return results, err
}

func (s *StatsService) GetCompletionStats() (map[string]interface{}, error) {
	var totalEpisodes int64
	database.DB.Model(&models.Episode{}).Count(&totalEpisodes)

	var completedEpisodes int64
	database.DB.Model(&models.PlaybackProgress{}).Where("completed = true").Count(&completedEpisodes)

	var inProgressEpisodes int64
	database.DB.Model(&models.PlaybackProgress{}).Where("completed = false AND current_time > 0").Count(&inProgressEpisodes)

	var totalListenedTime int64
	database.DB.Table("listening_histories").Select("COALESCE(SUM(duration), 0)").Scan(&totalListenedTime)

	return map[string]interface{}{
		"total_episodes":       totalEpisodes,
		"completed_episodes":   completedEpisodes,
		"in_progress_episodes": inProgressEpisodes,
		"completion_rate":      float64(completedEpisodes) / float64(totalEpisodes) * 100,
		"total_listened_time":  totalListenedTime,
	}, nil
}

func (s *StatsService) GetListeningHabits() (map[string]interface{}, error) {
	var hourDistribution []map[string]interface{}
	database.DB.Table("listening_histories").
		Select("EXTRACT(HOUR FROM start_time) as hour, COUNT(*) as count").
		Group("EXTRACT(HOUR FROM start_time)").
		Order("hour").
		Scan(&hourDistribution)

	var weekdayDistribution []map[string]interface{}
	database.DB.Table("listening_histories").
		Select("EXTRACT(DOW FROM start_time) as weekday, COUNT(*) as count").
		Group("EXTRACT(DOW FROM start_time)").
		Order("weekday").
		Scan(&weekdayDistribution)

	return map[string]interface{}{
		"hour_distribution":    hourDistribution,
		"weekday_distribution": weekdayDistribution,
	}, nil
}
