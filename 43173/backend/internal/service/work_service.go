package service

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"

	"music-platform/internal/model"
	"music-platform/internal/repository"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/database"
	"music-platform/pkg/redis"
	"music-platform/pkg/utils"
)

type WorkService struct {
	workRepo     *repository.WorkRepository
	userRepo     *repository.UserRepository
	communityRepo *repository.CommunityRepository
}

func NewWorkService() *WorkService {
	return &WorkService{
		workRepo:     repository.NewWorkRepository(),
		userRepo:     repository.NewUserRepository(),
		communityRepo: repository.NewCommunityRepository(),
	}
}

type UploadWorkRequest struct {
	Title       string  `json:"title" binding:"required"`
	ArtistName  string  `json:"artist_name" binding:"required"`
	AlbumID     *uint   `json:"album_id"`
	Type        string  `json:"type"`
	Genre       string  `json:"genre"`
	SubGenre    string  `json:"sub_genre"`
	Language    string  `json:"language"`
	Description string  `json:"description"`
	Lyrics      string  `json:"lyrics"`
	Composer    string  `json:"composer"`
	Lyricist    string  `json:"lyricist"`
	Arranger    string  `json:"arranger"`
	Producer    string  `json:"producer"`
	Explicit    bool    `json:"explicit"`
	Tags        []string `json:"tags"`
	Copyright   *CopyrightInfo `json:"copyright"`
}

type CopyrightInfo struct {
	Type        string  `json:"type"`
	Owner       string  `json:"owner"`
	LicenseType string  `json:"license_type"`
	RoyaltiesRate float64 `json:"royalties_rate"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	IsExclusive bool    `json:"is_exclusive"`
	Territory   string  `json:"territory"`
	ContractNo  string  `json:"contract_no"`
	Notes       string  `json:"notes"`
}

type BatchPublishRequest struct {
	WorkIDs []uint `json:"work_ids" binding:"required"`
}

type CreateAlbumRequest struct {
	Title       string `json:"title" binding:"required"`
	ArtistName  string `json:"artist_name" binding:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	ReleaseDate string `json:"release_date"`
	Type        string `json:"type"`
	Genre       string `json:"genre"`
	RecordLabel string `json:"record_label"`
	Upc         string `json:"upc"`
}

func (s *WorkService) UploadWork(userID uint, req *UploadWorkRequest, audioFile io.Reader, coverFile io.Reader, audioFilename string, coverFilename string) (*model.Work, error) {
	artistInfo, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	audioExt := utils.GetFileExt(audioFilename)
	if !utils.IsValidAudioExt(audioExt) {
		return nil, apperrors.ErrInvalidAudioFormat
	}

	audioDir := "./uploads/audio"
	coverDir := "./uploads/images"

	_ = os.MkdirAll(audioDir, 0755)
	_ = os.MkdirAll(coverDir, 0755)

	audioUUID := utils.GenerateUUID()
	audioNewName := audioUUID + audioExt
	audioPath := filepath.Join(audioDir, audioNewName)

	var fileMD5 string
	{
		hash := md5.New()
		tee := io.TeeReader(audioFile, hash)
		tmpFile, err := os.Create(audioPath)
		if err != nil {
			return nil, apperrors.ErrUploadFailed
		}
		_, _ = io.Copy(tmpFile, tee)
		tmpFile.Close()
		fileMD5 = hex.EncodeToString(hash.Sum(nil))
	}

	existingWork, _ := s.workRepo.FindByMD5(fileMD5)
	if existingWork != nil {
		_ = os.Remove(audioPath)
		return nil, apperrors.ErrWorkAlreadyExists
	}

	var coverURL string
	if coverFile != nil {
		coverExt := utils.GetFileExt(coverFilename)
		if utils.IsValidImageExt(coverExt) {
			coverUUID := utils.GenerateUUID()
			coverNewName := coverUUID + coverExt
			coverPath := filepath.Join(coverDir, coverNewName)

			coverFileSave, err := os.Create(coverPath)
			if err == nil {
				_, _ = io.Copy(coverFileSave, coverFile)
				coverFileSave.Close()
				coverURL = "/uploads/images/" + coverNewName
			}
		}
	}

	workType := model.WorkTypeSingle
	if req.Type != "" {
		workType = model.WorkType(req.Type)
	}

	work := &model.Work{
		UserID:       userID,
		ArtistID:     artistInfo.ID,
		Title:        req.Title,
		ArtistName:   req.ArtistName,
		AlbumID:      req.AlbumID,
		Type:         workType,
		Genre:        req.Genre,
		SubGenre:     req.SubGenre,
		Language:     req.Language,
		Description:  req.Description,
		Lyrics:       req.Lyrics,
		Composer:     req.Composer,
		Lyricist:     req.Lyricist,
		Arranger:     req.Arranger,
		Producer:     req.Producer,
		CoverURL:     coverURL,
		AudioURL:     "/uploads/audio/" + audioNewName,
		AudioFormat:  audioExt[1:],
		MD5:          fileMD5,
		Status:       model.WorkStatusDraft,
		Explicit:     req.Explicit,
		IsPublic:     true,
	}

	err = s.workRepo.Create(work)
	if err != nil {
		_ = os.Remove(audioPath)
		if coverURL != "" {
			_ = os.Remove(filepath.Join(coverDir, filepath.Base(coverURL)))
		}
		return nil, err
	}

	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}
		tag, err := s.workRepo.FindTagByName(tagName)
		if err != nil {
			tag = &model.Tag{
				Name: tagName,
				Type: "genre",
			}
			_ = s.workRepo.CreateTag(tag)
		}
		_ = s.workRepo.UpdateTagUsageCount(tag.ID)

		database.DB.Model(work).Association("Tags").Append(tag)
	}

	if req.Copyright != nil {
		copyright := &model.Copyright{
			WorkID:        work.ID,
			CopyrightType: req.Copyright.Type,
			Owner:         req.Copyright.Owner,
			LicenseType:   req.Copyright.LicenseType,
			RoyaltiesRate: req.Copyright.RoyaltiesRate,
			IsExclusive:   req.Copyright.IsExclusive,
			Territory:     req.Copyright.Territory,
			ContractNo:    req.Copyright.ContractNo,
			Notes:         req.Copyright.Notes,
		}
		if req.Copyright.StartDate != "" {
			startDate, _ := time.Parse("2006-01-02", req.Copyright.StartDate)
			copyright.StartDate = startDate
		}
		if req.Copyright.EndDate != "" {
			endDate, _ := time.Parse("2006-01-02", req.Copyright.EndDate)
			copyright.EndDate = &endDate
		}
		_ = s.workRepo.CreateCopyright(copyright)
	}

	_ = s.userRepo.UpdateArtistInfo(artistInfo)

	return work, nil
}

func (s *WorkService) GetWorkByID(id uint) (*model.Work, error) {
	work, err := s.workRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.ErrWorkNotFound
	}
	return work, nil
}

func (s *WorkService) ListWorks(page, pageSize int, keyword string, artistID uint, genre string, status int) ([]model.Work, int64, error) {
	return s.workRepo.List(page, pageSize, keyword, artistID, genre, status)
}

func (s *WorkService) GetArtistWorks(artistID uint, page, pageSize int) ([]model.Work, int64, error) {
	return s.workRepo.FindByArtistID(artistID, page, pageSize)
}

func (s *WorkService) UpdateWork(workID uint, userID uint, updates map[string]interface{}) error {
	work, err := s.workRepo.FindByID(workID)
	if err != nil {
		return apperrors.ErrWorkNotFound
	}

	if work.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.workRepo.UpdateWorkInfo(workID, updates)
}

func (s *WorkService) DeleteWork(workID uint, userID uint) error {
	work, err := s.workRepo.FindByID(workID)
	if err != nil {
		return apperrors.ErrWorkNotFound
	}

	if work.UserID != userID {
		return apperrors.ErrForbidden
	}

	if work.AudioURL != "" {
		_ = os.Remove("." + work.AudioURL)
	}
	if work.CoverURL != "" {
		_ = os.Remove("." + work.CoverURL)
	}

	return s.workRepo.Delete(workID)
}

func (s *WorkService) BatchPublish(req *BatchPublishRequest) error {
	for _, workID := range req.WorkIDs {
		work, err := s.workRepo.FindByID(workID)
		if err != nil {
			return apperrors.ErrWorkNotFound
		}
		if work.Status != model.WorkStatusDraft && work.Status != model.WorkStatusRejected {
			return apperrors.NewAppError(2005, "作品状态不允许发布")
		}
	}

	return s.workRepo.BatchUpdateStatus(req.WorkIDs, model.WorkStatusPublished)
}

func (s *WorkService) CreateAlbum(userID uint, req *CreateAlbumRequest) (*model.Album, error) {
	artistInfo, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	releaseDate, _ := time.Parse("2006-01-02", req.ReleaseDate)

	albumType := model.WorkTypeAlbum
	if req.Type != "" {
		albumType = model.WorkType(req.Type)
	}

	album := &model.Album{
		UserID:      userID,
		ArtistID:    artistInfo.ID,
		Title:       req.Title,
		ArtistName:  req.ArtistName,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		ReleaseDate: releaseDate,
		Type:        albumType,
		Genre:       req.Genre,
		RecordLabel: req.RecordLabel,
		Upc:         req.Upc,
		Status:      model.WorkStatusDraft,
	}

	err = s.workRepo.CreateAlbum(album)
	if err != nil {
		return nil, err
	}

	return album, nil
}

func (s *WorkService) GetAlbumByID(id uint) (*model.Album, error) {
	album, err := s.workRepo.FindAlbumByID(id)
	if err != nil {
		return nil, apperrors.ErrAlbumNotFound
	}
	return album, nil
}

func (s *WorkService) ListAlbums(page, pageSize int, keyword string, artistID uint, status int) ([]model.Album, int64, error) {
	return s.workRepo.ListAlbums(page, pageSize, keyword, artistID, status)
}

func (s *WorkService) UpdateAlbum(albumID uint, userID uint, updates map[string]interface{}) error {
	album, err := s.workRepo.FindAlbumByID(albumID)
	if err != nil {
		return apperrors.ErrAlbumNotFound
	}

	if album.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.workRepo.UpdateAlbum(album)
}

func (s *WorkService) DeleteAlbum(albumID uint, userID uint) error {
	album, err := s.workRepo.FindAlbumByID(albumID)
	if err != nil {
		return apperrors.ErrAlbumNotFound
	}

	if album.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.workRepo.DeleteAlbum(albumID)
}

func (s *WorkService) AddWorkToAlbum(albumID uint, workID uint, userID uint) error {
	album, err := s.workRepo.FindAlbumByID(albumID)
	if err != nil {
		return apperrors.ErrAlbumNotFound
	}

	if album.UserID != userID {
		return apperrors.ErrForbidden
	}

	work, err := s.workRepo.FindByID(workID)
	if err != nil {
		return apperrors.ErrWorkNotFound
	}

	work.AlbumID = &albumID
	err = s.workRepo.Update(work)
	if err != nil {
		return err
	}

	_ = s.workRepo.UpdateAlbumWorkCount(albumID, 1)

	return nil
}

func (s *WorkService) ListTags(keyword string) ([]model.Tag, error) {
	return s.workRepo.ListTags(keyword)
}

func (s *WorkService) UpdateWorkStatus(workID uint, status int) error {
	return s.workRepo.UpdateWorkStatus(workID, status)
}

func (s *WorkService) RecordPlay(workID uint, userID uint, duration int, ip string) error {
	work, err := s.workRepo.FindByID(workID)
	if err != nil {
		return apperrors.ErrWorkNotFound
	}

	_ = s.workRepo.UpdatePlayCount(workID, 1)

	dateKey := utils.GetDateKey(time.Now())
	weekKey := utils.GetWeekKey(time.Now())
	monthKey := utils.GetMonthKey(time.Now())

	redis.ZIncrBy("ranking:daily:plays:"+dateKey, 1, workID)
	redis.ZIncrBy("ranking:weekly:plays:"+weekKey, 1, workID)
	redis.ZIncrBy("ranking:monthly:plays:"+monthKey, 1, workID)

	hotScore := 1.0
	if duration > 0 {
		hotScore = float64(duration) / float64(work.Duration)
		if hotScore > 1 {
			hotScore = 1
		}
	}
	redis.ZIncrBy("ranking:daily:hot:"+dateKey, hotScore, workID)
	redis.ZIncrBy("ranking:weekly:hot:"+weekKey, hotScore, workID)
	redis.ZIncrBy("ranking:monthly:hot:"+monthKey, hotScore, workID)

	record := &model.PlayRecord{
		UserID:   userID,
		WorkID:   workID,
		ArtistID: work.ArtistID,
		Duration: duration,
		IP:       ip,
		Source:   "web",
	}
	_ = s.communityRepo.CreatePlayRecord(record)

	return nil
}
