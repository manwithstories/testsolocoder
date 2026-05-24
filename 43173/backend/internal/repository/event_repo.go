package repository

import (
	"music-platform/internal/model"
	"music-platform/pkg/database"
	"music-platform/pkg/utils"

	"gorm.io/gorm"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) Create(event *model.Event) error {
	return database.DB.Create(event).Error
}

func (r *EventRepository) FindByID(id uint) (*model.Event, error) {
	var event model.Event
	err := database.DB.Preload("User").First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(event *model.Event) error {
	return database.DB.Save(event).Error
}

func (r *EventRepository) UpdateEventInfo(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&model.Event{}).Where("id = ?", id).Updates(updates).Error
}

func (r *EventRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Event{}, id).Error
}

func (r *EventRepository) List(page, pageSize int, keyword string, artistID uint, city string, status int) ([]model.Event, int64, error) {
	var events []model.Event
	var total int64

	query := database.DB.Model(&model.Event{})

	if keyword != "" {
		query = query.Where("title LIKE ? OR venue LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if artistID > 0 {
		query = query.Where("artist_id = ?", artistID)
	}
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *EventRepository) UpdateSoldTickets(eventID uint, count int) error {
	return database.DB.Model(&model.Event{}).Where("id = ?", eventID).
		UpdateColumn("sold_tickets", gorm.Expr("sold_tickets + ?", count)).Error
}

func (r *EventRepository) UpdateViewCount(id uint) error {
	return database.DB.Model(&model.Event{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *EventRepository) UpdateLikeCount(id uint, count int64) error {
	return database.DB.Model(&model.Event{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", count)).Error
}

func (r *EventRepository) CreateTicket(ticket *model.Ticket) error {
	return database.DB.Create(ticket).Error
}

func (r *EventRepository) BatchCreateTickets(tickets []model.Ticket) error {
	return database.DB.Create(&tickets).Error
}

func (r *EventRepository) FindTicketByID(id uint) (*model.Ticket, error) {
	var ticket model.Ticket
	err := database.DB.Preload("Order").Preload("Event").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *EventRepository) GetTicketsByOrderID(orderID uint) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := database.DB.Where("order_id = ?", orderID).Find(&tickets).Error
	return tickets, err
}

func (r *EventRepository) GetTicketsByUserID(userID uint, page, pageSize int) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	query := database.DB.Model(&model.Ticket{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Event").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tickets).Error
	return tickets, total, err
}

func (r *EventRepository) UpdateTicketStatus(id uint, status model.TicketStatus) error {
	return database.DB.Model(&model.Ticket{}).Where("id = ?", id).Update("status", status).Error
}

func (r *EventRepository) CreateOrder(order *model.Order) error {
	return database.DB.Create(order).Error
}

func (r *EventRepository) FindOrderByID(id uint) (*model.Order, error) {
	var order model.Order
	err := database.DB.Preload("User").Preload("Event").Preload("Tickets").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *EventRepository) FindOrderByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := database.DB.Where("order_no = ?", orderNo).Preload("User").Preload("Event").Preload("Tickets").First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *EventRepository) UpdateOrderStatus(id uint, status model.OrderStatus) error {
	return database.DB.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *EventRepository) GetOrdersByUserID(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := database.DB.Model(&model.Order{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("Event").Preload("Tickets").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *EventRepository) GetOrdersByArtistID(artistID uint, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := database.DB.Model(&model.Order{}).Where("artist_id = ?", artistID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Preload("User").Preload("Event").Preload("Tickets").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}
