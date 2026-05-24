package service

import (
	"context"
	"fmt"
	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductionService struct {
	repo      *repository.ProductionRepo
	orderRepo *repository.OrderRepo
	redis     *redis.Client
}

func NewProductionService(repo *repository.ProductionRepo, orderRepo *repository.OrderRepo, rdb *redis.Client) *ProductionService {
	return &ProductionService{repo: repo, orderRepo: orderRepo, redis: rdb}
}

func (s *ProductionService) Create(ctx context.Context, orderID int64, operatorID int64) (*model.Production, error) {
	existing, _ := s.repo.GetByOrderID(ctx, orderID)
	if existing != nil {
		return nil, fmt.Errorf("该订单已有生产记录")
	}
	p := &model.Production{
		OrderID:         orderID,
		Status:          model.ProductionStatusQueued,
		ProgressPercent: 0,
		OperatorID:      operatorID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductionService) UpdateStatus(ctx context.Context, id int64, req *dto.UpdateProductionStatusRequest, operatorID int64) (*model.Production, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("生产记录不存在")
	}
	if !req.Status.Valid() {
		return nil, fmt.Errorf("无效的生产状态")
	}
	if p.Status == model.ProductionStatusCompleted {
		return nil, fmt.Errorf("已完成的生产记录不能修改")
	}
	p.Status = req.Status
	p.ProgressPercent = req.Status.ProgressPercent()
	p.OperatorID = operatorID
	p.Remark = req.Remark
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	if s.redis != nil {
		msg := fmt.Sprintf(`{"order_id":%d,"status":"%s","progress":%d}`, p.OrderID, p.Status, p.ProgressPercent)
		channel := fmt.Sprintf("production_updates:%d", p.OrderID)
		s.redis.Publish(ctx, channel, msg)
		s.redis.Publish(ctx, "production_updates", msg)
	}
	return p, nil
}

func (s *ProductionService) GetByID(ctx context.Context, id int64) (*model.Production, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductionService) GetByOrderID(ctx context.Context, orderID int64) (*model.Production, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}

func (s *ProductionService) List(ctx context.Context, filter repository.ProductionListFilter) ([]model.Production, int64, error) {
	return s.repo.List(ctx, filter)
}

func ToProductionResponse(p *model.Production, orderNo string) dto.ProductionResponse {
	return dto.ProductionResponse{
		ID:              p.ID,
		OrderID:         p.OrderID,
		OrderNo:         orderNo,
		Status:          p.Status,
		ProgressPercent: p.ProgressPercent,
		OperatorID:      p.OperatorID,
		Remark:          p.Remark,
		CreatedAt:       p.CreatedAt.Format(time.DateTime),
		UpdatedAt:       p.UpdatedAt.Format(time.DateTime),
	}
}
