package service

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"tea-platform/internal/models"
	"tea-platform/internal/repository"
	"tea-platform/internal/utils"
)

type CreateTastingRequest struct {
	TeaID           uint    `json:"tea_id" binding:"required"`
	BrewMethod      string  `json:"brew_method" binding:"omitempty,max=128"`
	WaterTemp       float64 `json:"water_temp"`
	BrewTime        int     `json:"brew_time"`
	WaterQuality    string  `json:"water_quality" binding:"omitempty,max=128"`
	TeaAmount       float64 `json:"tea_amount"`
	AromaScore      float64 `json:"aroma_score" binding:"required,min=1,max=10"`
	TasteScore      float64 `json:"taste_score" binding:"required,min=1,max=10"`
	AftertasteScore float64 `json:"aftertaste_score" binding:"required,min=1,max=10"`
	OverallScore    float64 `json:"overall_score" binding:"omitempty,min=1,max=10"`
	Notes           string  `json:"notes"`
}

type UpdateTastingRequest struct {
	TeaID           *uint    `json:"tea_id" binding:"omitempty"`
	BrewMethod      *string  `json:"brew_method" binding:"omitempty,max=128"`
	WaterTemp       *float64 `json:"water_temp"`
	BrewTime        *int     `json:"brew_time"`
	WaterQuality    *string  `json:"water_quality" binding:"omitempty,max=128"`
	TeaAmount       *float64 `json:"tea_amount"`
	AromaScore      *float64 `json:"aroma_score" binding:"omitempty,min=1,max=10"`
	TasteScore      *float64 `json:"taste_score" binding:"omitempty,min=1,max=10"`
	AftertasteScore *float64 `json:"aftertaste_score" binding:"omitempty,min=1,max=10"`
	OverallScore    *float64 `json:"overall_score" binding:"omitempty,min=1,max=10"`
	Notes           *string  `json:"notes"`
}

type TastingFilterRequest struct {
	UserID          uint    `form:"user_id"`
	TeaID           uint    `form:"tea_id"`
	MinOverallScore float64 `form:"min_overall_score"`
	MaxOverallScore float64 `form:"max_overall_score"`
	StartDate       string  `form:"start_date"`
	EndDate         string  `form:"end_date"`
	BrewMethod      string  `form:"brew_method"`
	Keyword         string  `form:"keyword"`
}

type TastingService struct {
	repo *repository.TastingRepository
	c    *gin.Context
}

func NewTastingService() *TastingService {
	return &TastingService{
		repo: repository.NewTastingRepository(),
	}
}

func (s *TastingService) SetContext(c *gin.Context) {
	s.c = c
}

func validateScore(score float64, name string) error {
	if score < 1 || score > 10 {
		return fmt.Errorf("%s 必须在 1-10 之间", name)
	}
	return nil
}

func calculateOverallScore(aroma, taste, aftertaste float64) float64 {
	return (aroma + taste + aftertaste) / 3
}

func (s *TastingService) CreateTastingRecord(userID uint, req *CreateTastingRequest) (*models.TastingRecord, error) {
	if err := validateScore(req.AromaScore, "香气评分"); err != nil {
		return nil, err
	}
	if err := validateScore(req.TasteScore, "滋味评分"); err != nil {
		return nil, err
	}
	if err := validateScore(req.AftertasteScore, "回甘评分"); err != nil {
		return nil, err
	}

	overallScore := req.OverallScore
	if overallScore == 0 {
		overallScore = calculateOverallScore(req.AromaScore, req.TasteScore, req.AftertasteScore)
	} else {
		if err := validateScore(overallScore, "综合评分"); err != nil {
			return nil, err
		}
	}

	record := &models.TastingRecord{
		UserID:          userID,
		TeaID:           req.TeaID,
		BrewMethod:      req.BrewMethod,
		WaterTemp:       req.WaterTemp,
		BrewTime:        req.BrewTime,
		WaterQuality:    req.WaterQuality,
		TeaAmount:       req.TeaAmount,
		AromaScore:      req.AromaScore,
		TasteScore:      req.TasteScore,
		AftertasteScore: req.AftertasteScore,
		OverallScore:    overallScore,
		Notes:           req.Notes,
	}

	if err := s.repo.CreateTastingRecord(record); err != nil {
		return nil, errors.New("创建品鉴记录失败")
	}

	return record, nil
}

func (s *TastingService) UpdateTastingRecord(userID uint, recordID uint, req *UpdateTastingRequest) error {
	record, err := s.repo.GetTastingRecordByID(recordID)
	if err != nil {
		return errors.New("品鉴记录不存在")
	}

	if record.UserID != userID {
		return errors.New("无权修改该品鉴记录")
	}

	if req.TeaID != nil {
		record.TeaID = *req.TeaID
	}
	if req.BrewMethod != nil {
		record.BrewMethod = *req.BrewMethod
	}
	if req.WaterTemp != nil {
		record.WaterTemp = *req.WaterTemp
	}
	if req.BrewTime != nil {
		record.BrewTime = *req.BrewTime
	}
	if req.WaterQuality != nil {
		record.WaterQuality = *req.WaterQuality
	}
	if req.TeaAmount != nil {
		record.TeaAmount = *req.TeaAmount
	}
	if req.AromaScore != nil {
		if err := validateScore(*req.AromaScore, "香气评分"); err != nil {
			return err
		}
		record.AromaScore = *req.AromaScore
	}
	if req.TasteScore != nil {
		if err := validateScore(*req.TasteScore, "滋味评分"); err != nil {
			return err
		}
		record.TasteScore = *req.TasteScore
	}
	if req.AftertasteScore != nil {
		if err := validateScore(*req.AftertasteScore, "回甘评分"); err != nil {
			return err
		}
		record.AftertasteScore = *req.AftertasteScore
	}
	if req.Notes != nil {
		record.Notes = *req.Notes
	}

	if req.OverallScore != nil {
		if err := validateScore(*req.OverallScore, "综合评分"); err != nil {
			return err
		}
		record.OverallScore = *req.OverallScore
	} else {
		record.OverallScore = calculateOverallScore(record.AromaScore, record.TasteScore, record.AftertasteScore)
	}

	if err := s.repo.UpdateTastingRecord(record); err != nil {
		return errors.New("更新品鉴记录失败")
	}

	return nil
}

func (s *TastingService) DeleteTastingRecord(userID uint, recordID uint) error {
	record, err := s.repo.GetTastingRecordByID(recordID)
	if err != nil {
		return errors.New("品鉴记录不存在")
	}

	if record.UserID != userID {
		return errors.New("无权删除该品鉴记录")
	}

	if err := s.repo.DeleteTastingRecord(recordID); err != nil {
		return errors.New("删除品鉴记录失败")
	}

	return nil
}

func (s *TastingService) GetTastingRecordDetail(id uint) (*models.TastingRecord, error) {
	record, err := s.repo.GetTastingRecordByID(id)
	if err != nil {
		return nil, errors.New("品鉴记录不存在")
	}
	return record, nil
}

func (s *TastingService) GetTastingRecordList(page, pageSize int, req *TastingFilterRequest) ([]models.TastingRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filters := repository.TastingFilters{
		UserID:          req.UserID,
		TeaID:           req.TeaID,
		MinOverallScore: req.MinOverallScore,
		MaxOverallScore: req.MaxOverallScore,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		BrewMethod:      req.BrewMethod,
		Keyword:         req.Keyword,
	}

	return s.repo.GetTastingRecordList(page, pageSize, filters)
}

func (s *TastingService) UploadTastingImage(recordID uint, userID uint, file *multipart.FileHeader) error {
	record, err := s.repo.GetTastingRecordByID(recordID)
	if err != nil {
		return errors.New("品鉴记录不存在")
	}

	if record.UserID != userID {
		return errors.New("无权上传该品鉴记录的图片")
	}

	imageURL, err := utils.UploadFile(s.c, file, "evaluation")
	if err != nil {
		return err
	}

	image := &models.TastingImage{
		TastingRecordID: record.ID,
		ImageURL:        imageURL,
	}

	if err := s.repo.CreateTastingImage(image); err != nil {
		return errors.New("保存图片失败")
	}

	return nil
}

func (s *TastingService) GetUserTastingStats(userID uint) (map[string]interface{}, error) {
	totalCount, avgScore, err := s.repo.GetTastingStats(userID)
	if err != nil {
		return nil, errors.New("获取品鉴统计失败")
	}

	return map[string]interface{}{
		"total_count":  totalCount,
		"avg_score":    avgScore,
	}, nil
}
