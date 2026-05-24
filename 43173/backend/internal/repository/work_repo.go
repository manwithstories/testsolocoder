package repository

import (
	"music-platform/internal/model"
	"music-platform/pkg/database"
	"music-platform/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type WorkRepository struct{}

func NewWorkRepository() *WorkRepository {
	return &WorkRepository{}
}

func (r *WorkRepository) Create(work *model.Work) error {
	return database.DB.Create(work).Error
}

func (r *WorkRepository) FindByID(id uint) (*model.Work, error) {
	var work model.Work
	err := database.DB.Preload("Tags").Preload("Copyright").First(&work, id).Error
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func (r *WorkRepository) FindByTitle(title string) (*model.Work, error) {
	var work model.Work
	err := database.DB.Where("title = ?", title).First(&work).Error
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func (r *WorkRepository) Update(work *model.Work) error {
	return database.DB.Save(work).Error
}

func (r *WorkRepository) UpdateWorkInfo(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&model.Work{}).Where("id = ?", id).Updates(updates).Error
}

func (r *WorkRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Work{}, id).Error
}

func (r *WorkRepository) List(page, pageSize int, keyword string, artistID uint, genre string, status int) ([]model.Work, int64, error) {
	var works []model.Work
	var total int64

	query := database.DB.Model(&model.Work{})

	if keyword != "" {
		query = query.Where("title LIKE ? OR artist_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if artistID > 0 {
		query = query.Where("artist_id = ?", artistID)
	}
	if genre != "" {
		query = query.Where("genre = ?", genre)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Tags").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&works).Error
	if err != nil {
		return nil, 0, err
	}

	return works, total, nil
}

func (r *WorkRepository) FindByArtistID(artistID uint, page, pageSize int) ([]model.Work, int64, error) {
	var works []model.Work
	var total int64

	query := database.DB.Model(&model.Work{}).Where("artist_id = ?", artistID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Tags").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&works).Error
	if err != nil {
		return nil, 0, err
	}

	return works, total, nil
}

func (r *WorkRepository) FindByAlbumID(albumID uint) ([]model.Work, error) {
	var works []model.Work
	err := database.DB.Where("album_id = ?", albumID).Preload("Tags").Order("created_at ASC").Find(&works).Error
	return works, err
}

func (r *WorkRepository) FindByMD5(md5 string) (*model.Work, error) {
	var work model.Work
	err := database.DB.Where("md5 = ?", md5).First(&work).Error
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func (r *WorkRepository) UpdatePlayCount(id uint, count int64) error {
	return database.DB.Model(&model.Work{}).Where("id = ?", id).
		UpdateColumn("play_count", gorm.Expr("play_count + ?", count)).Error
}

func (r *WorkRepository) UpdateLikeCount(id uint, count int64) error {
	return database.DB.Model(&model.Work{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", count)).Error
}

func (r *WorkRepository) UpdateCommentCount(id uint, count int64) error {
	return database.DB.Model(&model.Work{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", count)).Error
}

func (r *WorkRepository) UpdateShareCount(id uint, count int64) error {
	return database.DB.Model(&model.Work{}).Where("id = ?", id).
		UpdateColumn("share_count", gorm.Expr("share_count + ?", count)).Error
}

func (r *WorkRepository) BatchUpdateStatus(ids []uint, status model.WorkStatus) error {
	return database.DB.Model(&model.Work{}).Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":       status,
			"published_at": time.Now(),
		}).Error
}

func (r *WorkRepository) BatchCreate(works []model.Work) error {
	return database.DB.Create(&works).Error
}

func (r *WorkRepository) FindByStatus(status model.WorkStatus) ([]model.Work, error) {
	var works []model.Work
	err := database.DB.Where("status = ?", status).Find(&works).Error
	return works, err
}

func (r *WorkRepository) CreateTag(tag *model.Tag) error {
	return database.DB.Create(tag).Error
}

func (r *WorkRepository) FindTagByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := database.DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *WorkRepository) ListTags(keyword string) ([]model.Tag, error) {
	var tags []model.Tag
	query := database.DB.Model(&model.Tag{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("usage_count DESC").Find(&tags).Error
	return tags, err
}

func (r *WorkRepository) UpdateWorkStatus(id uint, status int) error {
	return database.DB.Model(&model.Work{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *WorkRepository) UpdateTagUsageCount(id uint) error {
	return database.DB.Model(&model.Tag{}).Where("id = ?", id).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
}

func (r *WorkRepository) CreateCopyright(copyright *model.Copyright) error {
	return database.DB.Create(copyright).Error
}

func (r *WorkRepository) UpdateCopyright(copyright *model.Copyright) error {
	return database.DB.Save(copyright).Error
}

func (r *WorkRepository) FindCopyrightByWorkID(workID uint) (*model.Copyright, error) {
	var copyright model.Copyright
	err := database.DB.Where("work_id = ?", workID).First(&copyright).Error
	if err != nil {
		return nil, err
	}
	return &copyright, nil
}

func (r *WorkRepository) CreateAlbum(album *model.Album) error {
	return database.DB.Create(album).Error
}

func (r *WorkRepository) FindAlbumByID(id uint) (*model.Album, error) {
	var album model.Album
	err := database.DB.Preload("Works").First(&album, id).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (r *WorkRepository) UpdateAlbum(album *model.Album) error {
	return database.DB.Save(album).Error
}

func (r *WorkRepository) DeleteAlbum(id uint) error {
	return database.DB.Delete(&model.Album{}, id).Error
}

func (r *WorkRepository) ListAlbums(page, pageSize int, keyword string, artistID uint, status int) ([]model.Album, int64, error) {
	var albums []model.Album
	var total int64

	query := database.DB.Model(&model.Album{})

	if keyword != "" {
		query = query.Where("title LIKE ? OR artist_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if artistID > 0 {
		query = query.Where("artist_id = ?", artistID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Works").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&albums).Error
	if err != nil {
		return nil, 0, err
	}

	return albums, total, nil
}

func (r *WorkRepository) UpdateAlbumWorkCount(albumID uint, count int) error {
	return database.DB.Model(&model.Album{}).Where("id = ?", albumID).
		UpdateColumn("work_count", gorm.Expr("work_count + ?", count)).Error
}
