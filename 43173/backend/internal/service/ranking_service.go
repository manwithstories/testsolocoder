package service

import (
	"encoding/json"
	"strconv"
	"time"

	"music-platform/internal/model"
	"music-platform/internal/repository"
	"music-platform/pkg/redis"
	"music-platform/pkg/utils"
)

type RankingService struct {
	workRepo *repository.WorkRepository
}

func NewRankingService() *RankingService {
	return &RankingService{
		workRepo: repository.NewWorkRepository(),
	}
}

type RankingItem struct {
	Rank      int         `json:"rank"`
	WorkID    uint        `json:"work_id"`
	Title     string      `json:"title"`
	ArtistName string     `json:"artist_name"`
	CoverURL  string      `json:"cover_url"`
	Score     float64     `json:"score"`
	Work      *model.Work `json:"work,omitempty"`
}

const (
	RankingTypeDaily   = "daily"
	RankingTypeWeekly  = "weekly"
	RankingTypeMonthly = "monthly"

	RankingCategoryPlays = "plays"
	RankingCategoryHot   = "hot"
	RankingCategoryLikes = "likes"
)

func (s *RankingService) GetRanking(rankingType string, category string, limit int) ([]RankingItem, error) {
	var key string
	now := time.Now()

	switch rankingType {
	case RankingTypeDaily:
		dateKey := utils.GetDateKey(now)
		key = "ranking:daily:" + category + ":" + dateKey
	case RankingTypeWeekly:
		weekKey := utils.GetWeekKey(now)
		key = "ranking:weekly:" + category + ":" + weekKey
	case RankingTypeMonthly:
		monthKey := utils.GetMonthKey(now)
		key = "ranking:monthly:" + category + ":" + monthKey
	default:
		dateKey := utils.GetDateKey(now)
		key = "ranking:daily:" + category + ":" + dateKey
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	results, err := redis.ZRevRangeWithScores(key, 0, int64(limit-1))
	if err != nil {
		return nil, err
	}

	var items []RankingItem
	for i, result := range results {
		workID, _ := strconv.ParseUint(result.Member, 10, 64)
		work, _ := s.workRepo.FindByID(uint(workID))

		item := RankingItem{
			Rank:  i + 1,
			WorkID: uint(workID),
			Score: result.Score,
		}

		if work != nil {
			item.Title = work.Title
			item.ArtistName = work.ArtistName
			item.CoverURL = work.CoverURL
			item.Work = work
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *RankingService) GetWorkRanking(workID uint, rankingType string, category string) (int64, float64, error) {
	var key string
	now := time.Now()

	switch rankingType {
	case RankingTypeDaily:
		dateKey := utils.GetDateKey(now)
		key = "ranking:daily:" + category + ":" + dateKey
	case RankingTypeWeekly:
		weekKey := utils.GetWeekKey(now)
		key = "ranking:weekly:" + category + ":" + weekKey
	case RankingTypeMonthly:
		monthKey := utils.GetMonthKey(now)
		key = "ranking:monthly:" + category + ":" + monthKey
	default:
		dateKey := utils.GetDateKey(now)
		key = "ranking:daily:" + category + ":" + dateKey
	}

	rank, err := redis.ZRank(key, strconv.FormatUint(uint64(workID), 10))
	if err != nil {
		return 0, 0, err
	}

	score, err := redis.ZScore(key, strconv.FormatUint(uint64(workID), 10))
	if err != nil {
		return 0, 0, err
	}

	return rank + 1, score, nil
}

func (s *RankingService) IncrementPlay(workID uint) {
	now := time.Now()
	dateKey := utils.GetDateKey(now)
	weekKey := utils.GetWeekKey(now)
	monthKey := utils.GetMonthKey(now)

	workIDStr := strconv.FormatUint(uint64(workID), 10)

	redis.ZIncrBy("ranking:daily:plays:"+dateKey, 1, workIDStr)
	redis.ZIncrBy("ranking:weekly:plays:"+weekKey, 1, workIDStr)
	redis.ZIncrBy("ranking:monthly:plays:"+monthKey, 1, workIDStr)

	_ = s.workRepo.UpdatePlayCount(workID, 1)
}

func (s *RankingService) IncrementLike(workID uint) {
	now := time.Now()
	dateKey := utils.GetDateKey(now)
	weekKey := utils.GetWeekKey(now)
	monthKey := utils.GetMonthKey(now)

	workIDStr := strconv.FormatUint(uint64(workID), 10)

	redis.ZIncrBy("ranking:daily:likes:"+dateKey, 1, workIDStr)
	redis.ZIncrBy("ranking:weekly:likes:"+weekKey, 1, workIDStr)
	redis.ZIncrBy("ranking:monthly:likes:"+monthKey, 1, workIDStr)

	_ = s.workRepo.UpdateLikeCount(workID, 1)
}

func (s *RankingService) UpdateHotScore(workID uint, playDuration int, totalDuration int) {
	if totalDuration <= 0 {
		return
	}

	completionRate := float64(playDuration) / float64(totalDuration)
	if completionRate > 1 {
		completionRate = 1
	}

	hotScore := completionRate * 10

	now := time.Now()
	dateKey := utils.GetDateKey(now)
	weekKey := utils.GetWeekKey(now)
	monthKey := utils.GetMonthKey(now)

	workIDStr := strconv.FormatUint(uint64(workID), 10)

	redis.ZIncrBy("ranking:daily:hot:"+dateKey, hotScore, workIDStr)
	redis.ZIncrBy("ranking:weekly:hot:"+weekKey, hotScore, workIDStr)
	redis.ZIncrBy("ranking:monthly:hot:"+monthKey, hotScore, workIDStr)
}

func (s *RankingService) GetCacheRanking(key string, limit int) ([]RankingItem, error) {
	cacheKey := "cache:ranking:" + key

	cached, err := redis.Get(cacheKey)
	if err == nil && cached != "" {
		var items []RankingItem
		_ = json.Unmarshal([]byte(cached), &items)
		return items, nil
	}

	results, err := redis.ZRevRangeWithScores(key, 0, int64(limit-1))
	if err != nil {
		return nil, err
	}

	var items []RankingItem
	for i, result := range results {
		workID, _ := strconv.ParseUint(result.Member, 10, 64)
		work, _ := s.workRepo.FindByID(uint(workID))

		item := RankingItem{
			Rank:   i + 1,
			WorkID: uint(workID),
			Score:  result.Score,
		}

		if work != nil {
			item.Title = work.Title
			item.ArtistName = work.ArtistName
			item.CoverURL = work.CoverURL
		}

		items = append(items, item)
	}

	data, _ := json.Marshal(items)
	_ = redis.Set(cacheKey, string(data), 5*time.Minute)

	return items, nil
}
