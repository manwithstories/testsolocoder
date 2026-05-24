package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
)

// OrderService 订单业务逻辑层
type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	db          *gorm.DB
}

// NewOrderService 创建订单服务
func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	db *gorm.DB,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		db:          db,
	}
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return fmt.Sprintf("OD%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
}

// toCustomOptionModels 将 DTO 定制选项转为模型
func toCustomOptionModels(opts []dto.CustomOptionDTO) []model.CustomOption {
	list := make([]model.CustomOption, 0, len(opts))
	for _, o := range opts {
		list = append(list, model.CustomOption{
			OptionType:  o.OptionType,
			OptionValue: o.OptionValue,
			PriceAdjust: o.PriceAdjust,
		})
	}
	return list
}

// computeItemSubtotal 计算单项小计：(基价 + 选项加价合计) × 数量
func computeItemSubtotal(basePrice float64, options []model.CustomOption, quantity int) float64 {
	var optionSum float64
	for _, o := range options {
		optionSum += o.PriceAdjust
	}
	return (basePrice + optionSum) * float64(quantity)
}

// CreateInquiry 业主发起询价
func (s *OrderService) CreateInquiry(ownerID uint, req *dto.CreateInquiryRequest) (*model.Order, []*model.OrderItem, error) {
	if len(req.Items) == 0 {
		return nil, nil, errors.New("询价项不能为空")
	}

	// 校验产品是否存在并获取信息
	productIDs := make([]uint, 0, len(req.Items))
	for _, it := range req.Items {
		productIDs = append(productIDs, it.ProductID)
	}
	// 校验所有产品属于同一厂商
	var manufacturerID uint
	productMap := make(map[uint]*model.Product)
	for _, it := range req.Items {
		p, err := s.productRepo.GetByID(it.ProductID)
		if err != nil {
			return nil, nil, err
		}
		if p == nil {
			return nil, nil, fmt.Errorf("产品 %d 不存在", it.ProductID)
		}
		if manufacturerID == 0 {
			manufacturerID = p.ManufacturerID
		} else if manufacturerID != p.ManufacturerID {
			return nil, nil, errors.New("询价单中的产品必须属于同一厂商")
		}
		productMap[p.ID] = p
	}

	// 构建订单项
	items := make([]*model.OrderItem, 0, len(req.Items))
	var totalAmount float64
	for _, it := range req.Items {
		p := productMap[it.ProductID]
		opts := toCustomOptionModels(it.CustomOptions)
		subtotal := computeItemSubtotal(p.BasePrice, opts, it.Quantity)
		item := &model.OrderItem{
			ProductID:   p.ID,
			ProductName: p.Name,
			BasePrice:   p.BasePrice,
			Quantity:    it.Quantity,
			Subtotal:    subtotal,
		}
		if err := item.SetCustomOptions(opts); err != nil {
			return nil, nil, err
		}
		items = append(items, item)
		totalAmount += subtotal
	}

	// 生成订单号（保证唯一）
	orderNo := generateOrderNo()
	for {
		exist, err := s.orderRepo.GetByOrderNo(orderNo)
		if err != nil {
			return nil, nil, err
		}
		if exist == nil {
			break
		}
		orderNo = generateOrderNo()
	}

	order := &model.Order{
		OrderNo:        orderNo,
		OwnerID:        ownerID,
		DesignerID:     req.DesignerID,
		ManufacturerID: manufacturerID,
		TotalAmount:    totalAmount,
		Discount:       0,
		FinalAmount:    totalAmount,
		Status:         model.OrderStatusInquiry,
		Address:        req.Address,
		ContactName:    req.ContactName,
		ContactPhone:   req.ContactPhone,
		Remark:         req.Remark,
	}

	history := &model.OrderHistory{
		Status:       model.OrderStatusInquiry,
		OperatorID:   ownerID,
		OperatorRole: model.RoleOwner,
		Remark:       "业主发起询价",
	}

	if err := s.orderRepo.CreateInquiry(order, items, history); err != nil {
		return nil, nil, err
	}
	return order, items, nil
}

// Quote 厂商报价
func (s *OrderService) Quote(orderID uint, manufacturerID uint, req *dto.QuoteRequest) (*model.Order, []*model.OrderItem, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, nil, err
	}
	if order == nil {
		return nil, nil, errors.New("订单不存在")
	}
	if order.ManufacturerID != manufacturerID {
		return nil, nil, errors.New("无权报价他人的订单")
	}
	if order.Status != model.OrderStatusInquiry {
		return nil, nil, errors.New("当前订单状态无法报价")
	}

	items, err := s.orderRepo.GetItems(order.ID)
	if err != nil {
		return nil, nil, err
	}
	if len(items) == 0 {
		return nil, nil, errors.New("订单无订单项")
	}

	// 构建 id -> item 映射
	itemMap := make(map[uint]*model.OrderItem, len(items))
	for _, it := range items {
		itemMap[it.ID] = it
	}

	var totalAmount float64
	updatedItems := make([]*model.OrderItem, 0, len(req.ItemPrices))
	for _, q := range req.ItemPrices {
		it, ok := itemMap[q.OrderItemID]
		if !ok {
			return nil, nil, fmt.Errorf("订单项 %d 不存在", q.OrderItemID)
		}
		it.BasePrice = q.BasePrice
		opts, err := it.ParseCustomOptions()
		if err != nil {
			return nil, nil, err
		}
		it.Subtotal = computeItemSubtotal(q.BasePrice, opts, it.Quantity)
		updatedItems = append(updatedItems, it)
		totalAmount += it.Subtotal
	}

	// 校验折扣不超过总价
	if req.Discount < 0 {
		return nil, nil, errors.New("折扣不能为负数")
	}
	if req.Discount > totalAmount {
		return nil, nil, errors.New("折扣不能超过总价")
	}

	order.TotalAmount = totalAmount
	order.Discount = req.Discount
	order.FinalAmount = totalAmount - req.Discount
	order.Status = model.OrderStatusQuoted

	history := &model.OrderHistory{
		Status:       model.OrderStatusQuoted,
		OperatorID:   manufacturerID,
		OperatorRole: model.RoleManufacturer,
		Remark:       req.Remark,
	}

	if err := s.orderRepo.Quote(order, updatedItems, history); err != nil {
		return nil, nil, err
	}
	return order, updatedItems, nil
}

// ConfirmOrder 业主确认下单
func (s *OrderService) ConfirmOrder(orderID uint, ownerID uint, remark string) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}
	if order.OwnerID != ownerID {
		return nil, errors.New("无权操作他人的订单")
	}
	if order.Status != model.OrderStatusQuoted {
		return nil, errors.New("当前订单状态无法确认")
	}

	order.Status = model.OrderStatusConfirmed
	order.UpdatedAt = time.Now()

	history := &model.OrderHistory{
		Status:       model.OrderStatusConfirmed,
		OperatorID:   ownerID,
		OperatorRole: model.RoleOwner,
		Remark:       remark,
	}

	if err := s.orderRepo.UpdateStatus(order, history); err != nil {
		return nil, err
	}
	return order, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderID uint, userID uint, role string, remark string) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	// 权限校验：业主可取消自己的订单，厂商可取消自己的订单
	switch role {
	case model.RoleOwner:
		if order.OwnerID != userID {
			return nil, errors.New("无权操作他人的订单")
		}
	case model.RoleManufacturer:
		if order.ManufacturerID != userID {
			return nil, errors.New("无权操作他人的订单")
		}
	default:
		return nil, errors.New("无权取消订单")
	}

	if !model.CanTransitTo(order.Status, model.OrderStatusCancelled) {
		return nil, errors.New("当前订单状态无法取消")
	}

	order.Status = model.OrderStatusCancelled
	order.UpdatedAt = time.Now()

	history := &model.OrderHistory{
		Status:       model.OrderStatusCancelled,
		OperatorID:   userID,
		OperatorRole: role,
		Remark:       remark,
	}

	if err := s.orderRepo.UpdateStatus(order, history); err != nil {
		return nil, err
	}
	return order, nil
}

// UpdateStatus 更新订单状态（用于生产/发货/完成等）
func (s *OrderService) UpdateStatus(orderID uint, userID uint, role string, newStatus string, remark string) (*model.Order, error) {
	if !model.ValidOrderStatus(newStatus) {
		return nil, errors.New("目标状态不合法")
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	// 权限校验
	switch role {
	case model.RoleManufacturer:
		if order.ManufacturerID != userID {
			return nil, errors.New("无权操作他人的订单")
		}
		// 厂商可推进：paid -> producing -> shipped -> completed
		switch newStatus {
		case model.OrderStatusProducing, model.OrderStatusShipped, model.OrderStatusCompleted:
		default:
			return nil, errors.New("厂商无权执行该状态变更")
		}
	case model.RoleOwner:
		if order.OwnerID != userID {
			return nil, errors.New("无权操作他人的订单")
		}
		// 业主可推进：confirmed -> paid
		if newStatus != model.OrderStatusPaid {
			return nil, errors.New("业主无权执行该状态变更")
		}
	default:
		return nil, errors.New("无权执行状态变更")
	}

	if !model.CanTransitTo(order.Status, newStatus) {
		return nil, fmt.Errorf("订单状态无法从 %s 转为 %s", order.Status, newStatus)
	}

	order.Status = newStatus
	order.UpdatedAt = time.Now()

	history := &model.OrderHistory{
		Status:       newStatus,
		OperatorID:   userID,
		OperatorRole: role,
		Remark:       remark,
	}

	if err := s.orderRepo.UpdateStatus(order, history); err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID 获取订单详情（含订单项、历史）
func (s *OrderService) GetByID(orderID uint, userID uint, role string) (*model.Order, []*model.OrderItem, []*model.OrderHistory, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, nil, nil, err
	}
	if order == nil {
		return nil, nil, nil, errors.New("订单不存在")
	}

	// 权限校验：业主查看自己的、厂商查看自己的、设计师查看与自己相关的
	switch role {
	case model.RoleOwner:
		if order.OwnerID != userID {
			return nil, nil, nil, errors.New("无权查看他人的订单")
		}
	case model.RoleManufacturer:
		if order.ManufacturerID != userID {
			return nil, nil, nil, errors.New("无权查看他人的订单")
		}
	case model.RoleDesigner:
		if order.DesignerID != userID {
			return nil, nil, nil, errors.New("无权查看他人的订单")
		}
	}

	items, err := s.orderRepo.GetItems(order.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	histories, err := s.orderRepo.GetHistories(order.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	return order, items, histories, nil
}

// List 分页查询订单列表
func (s *OrderService) List(params *dto.OrderListRequest) ([]*model.Order, int64, error) {
	p := &repository.OrderListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
		Status:   params.Status,
		Role:     params.Role,
		UserID:   params.UserID,
	}
	return s.orderRepo.List(p)
}

// GetHistories 获取订单历史
func (s *OrderService) GetHistories(orderID uint) ([]*model.OrderHistory, error) {
	return s.orderRepo.GetHistories(orderID)
}
