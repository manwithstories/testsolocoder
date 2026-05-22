package handlers

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type GuestHandler struct{}

func NewGuestHandler() *GuestHandler {
	return &GuestHandler{}
}

type GuestRequest struct {
	FirstName   string `json:"first_name" binding:"required,max=50"`
	LastName    string `json:"last_name" binding:"required,max=50"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Group       string `json:"group"`
	Relation    string `json:"relation"`
	RSVPStatus  string `json:"rsvp_status"`
	PlusOne     bool   `json:"plus_one"`
	PlusOneName string `json:"plus_one_name"`
	TableID     *uint  `json:"table_id"`
	SeatNumber  int    `json:"seat_number"`
	Notes       string `json:"notes"`
	IsVIP       bool   `json:"is_vip"`
}

type TableRequest struct {
	TableName   string `json:"table_name" binding:"required,max=100"`
	TableNumber int    `json:"table_number"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	SeatsJSON   string `json:"seats_json"`
	Notes       string `json:"notes"`
}

type SeatAssignment struct {
	SeatNumber  int   `json:"seat_number"`
	GuestID     uint  `json:"guest_id"`
}

func (h *GuestHandler) Create(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	var req GuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	fullName := req.FullName
	if fullName == "" {
		fullName = req.FirstName + " " + req.LastName
	}

	guest := models.Guest{
		WeddingID:   weddingID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		FullName:    fullName,
		Email:       req.Email,
		Phone:       req.Phone,
		Group:       req.Group,
		Relation:    req.Relation,
		RSVPStatus:  req.RSVPStatus,
		PlusOne:     req.PlusOne,
		PlusOneName: req.PlusOneName,
		TableID:     req.TableID,
		SeatNumber:  req.SeatNumber,
		Notes:       req.Notes,
		IsVIP:       req.IsVIP,
	}

	if guest.RSVPStatus == "" {
		guest.RSVPStatus = "pending"
	}

	if err := db.Create(&guest).Error; err != nil {
		response.InternalError(c, "Failed to create guest")
		return
	}

	response.Created(c, guest)
}

func (h *GuestHandler) GetList(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var guests []models.Guest
	var total int64

	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	search := c.Query("search")
	group := c.Query("group")
	rsvpStatus := c.Query("rsvp_status")

	query := db.Model(&models.Guest{}).Where("wedding_id = ?", weddingID)

	if search != "" {
		query = query.Where("full_name LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if group != "" {
		query = query.Where("`group` = ?", group)
	}
	if rsvpStatus != "" {
		query = query.Where("rsvp_status = ?", rsvpStatus)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("is_vip DESC, full_name ASC").Find(&guests)

	response.Paginated(c, guests, total, page, pageSize)
}

func (h *GuestHandler) GetByID(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var guest models.Guest
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&guest).Error; err != nil {
		response.NotFound(c, "Guest not found")
		return
	}

	response.Success(c, guest)
}

func (h *GuestHandler) Update(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	var req GuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var guest models.Guest
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&guest).Error; err != nil {
		response.NotFound(c, "Guest not found")
		return
	}

	fullName := req.FullName
	if fullName == "" {
		fullName = req.FirstName + " " + req.LastName
	}

	updates := map[string]interface{}{
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
		"full_name":     fullName,
		"email":         req.Email,
		"phone":         req.Phone,
		"group":         req.Group,
		"relation":      req.Relation,
		"rsvp_status":   req.RSVPStatus,
		"plus_one":      req.PlusOne,
		"plus_one_name": req.PlusOneName,
		"table_id":      req.TableID,
		"seat_number":   req.SeatNumber,
		"notes":         req.Notes,
		"is_vip":        req.IsVIP,
	}

	db.Model(&guest).Updates(updates)

	response.Success(c, guest)
}

func (h *GuestHandler) Delete(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var guest models.Guest
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&guest).Error; err != nil {
		response.NotFound(c, "Guest not found")
		return
	}

	db.Delete(&guest)

	response.Success(c, gin.H{"message": "Guest deleted successfully"})
}

func (h *GuestHandler) UpdateRSVPStatus(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	type RSVPRequest struct {
		RSVPStatus string `json:"rsvp_status" binding:"required,oneof=pending accepted declined"`
	}

	var req RSVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid RSVP status")
		return
	}

	db := database.GetDB()

	db.Model(&models.Guest{}).Where("id = ? AND wedding_id = ?", id, weddingID).Update("rsvp_status", req.RSVPStatus)

	response.Success(c, gin.H{"message": "RSVP status updated successfully"})
}

func (h *GuestHandler) ImportGuests(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Please upload a file")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.BadRequest(c, "Failed to open file")
		return
	}
	defer f.Close()

	var guests []models.Guest

	filename := file.Filename
	if strings.HasSuffix(filename, ".csv") {
		reader := csv.NewReader(f)
		reader.Read()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			guest := models.Guest{
				WeddingID: weddingID,
				FirstName: record[0],
				LastName:  record[1],
				FullName:  record[0] + " " + record[1],
				Email:     record[2],
				Phone:     record[3],
				Group:     record[4],
			}
			if len(record) > 5 && record[5] != "" {
				guest.RSVPStatus = record[5]
			} else {
				guest.RSVPStatus = "pending"
			}

			guests = append(guests, guest)
		}
	} else if strings.HasSuffix(filename, ".xlsx") || strings.HasSuffix(filename, ".xls") {
		xlsx, err := excelize.OpenReader(f)
		if err != nil {
			response.BadRequest(c, "Failed to read Excel file")
			return
		}

		sheetName := xlsx.GetSheetName(0)
		rows, err := xlsx.GetRows(sheetName)
		if err != nil {
			response.BadRequest(c, "Failed to read Excel rows")
			return
		}

		for i, row := range rows {
			if i == 0 {
				continue
			}
			guest := models.Guest{
				WeddingID: weddingID,
				RSVPStatus: "pending",
			}
			if len(row) > 0 {
				guest.FirstName = row[0]
			}
			if len(row) > 1 {
				guest.LastName = row[1]
			}
			guest.FullName = guest.FirstName + " " + guest.LastName
			if len(row) > 2 {
				guest.Email = row[2]
			}
			if len(row) > 3 {
				guest.Phone = row[3]
			}
			if len(row) > 4 {
				guest.Group = row[4]
			}
			if len(row) > 5 && row[5] != "" {
				guest.RSVPStatus = row[5]
			}

			guests = append(guests, guest)
		}
	} else {
		response.BadRequest(c, "Unsupported file format. Please upload CSV or Excel file")
		return
	}

	if len(guests) > 0 {
		db := database.GetDB()
		db.Create(&guests)
	}

	response.Success(c, gin.H{
		"message": "Guests imported successfully",
		"count":   len(guests),
	})
}

func (h *GuestHandler) ExportGuests(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var guests []models.Guest
	db.Where("wedding_id = ?", weddingID).Order("full_name ASC").Find(&guests)

	f := excelize.NewFile()
	sheet := "Guests"
	f.NewSheet(sheet)

	headers := []string{"First Name", "Last Name", "Email", "Phone", "Group", "Relation", "RSVP Status", "Table", "Seat", "VIP"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for i, guest := range guests {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), guest.FirstName)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), guest.LastName)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), guest.Email)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), guest.Phone)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), guest.Group)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), guest.Relation)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), guest.RSVPStatus)
		if guest.TableID != nil {
			f.SetCellValue(sheet, "H"+strconv.Itoa(row), *guest.TableID)
		}
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), guest.SeatNumber)
		f.SetCellValue(sheet, "J"+strconv.Itoa(row), guest.IsVIP)
	}

	f.DeleteSheet("Sheet1")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=guests.xlsx")
	f.Write(c.Writer)
}

func (h *GuestHandler) CreateTable(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	var req TableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	table := models.GuestTable{
		WeddingID:   weddingID,
		TableName:   req.TableName,
		TableNumber: req.TableNumber,
		Capacity:    req.Capacity,
		SeatsJSON:   req.SeatsJSON,
		Notes:       req.Notes,
	}

	if err := db.Create(&table).Error; err != nil {
		response.InternalError(c, "Failed to create table")
		return
	}

	response.Created(c, table)
}

func (h *GuestHandler) GetTables(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var tables []models.GuestTable
	db.Where("wedding_id = ?", weddingID).Order("table_number ASC").Find(&tables)

	var result []gin.H
	for _, table := range tables {
		var guests []models.Guest
		db.Where("wedding_id = ? AND table_id = ?", weddingID, table.ID).Order("seat_number ASC").Find(&guests)

		result = append(result, gin.H{
			"id":           table.ID,
			"table_name":   table.TableName,
			"table_number": table.TableNumber,
			"capacity":     table.Capacity,
			"notes":        table.Notes,
			"created_at":   table.CreatedAt,
			"updated_at":   table.UpdatedAt,
			"table_guests": guests,
		})
	}

	response.Success(c, result)
}

func (h *GuestHandler) UpdateTable(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	var req TableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var table models.GuestTable
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&table).Error; err != nil {
		response.NotFound(c, "Table not found")
		return
	}

	db.Model(&table).Updates(map[string]interface{}{
		"table_name":   req.TableName,
		"table_number": req.TableNumber,
		"capacity":     req.Capacity,
		"seats_json":   req.SeatsJSON,
		"notes":        req.Notes,
	})

	response.Success(c, table)
}

func (h *GuestHandler) DeleteTable(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var table models.GuestTable
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&table).Error; err != nil {
		response.NotFound(c, "Table not found")
		return
	}

	db.Model(&models.Guest{}).Where("table_id = ?", id).Update("table_id", nil)

	db.Delete(&table)

	response.Success(c, gin.H{"message": "Table deleted successfully"})
}

func (h *GuestHandler) AssignSeat(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	type SeatAssignRequest struct {
		GuestID    uint `json:"guest_id" binding:"required"`
		TableID    uint `json:"table_id" binding:"required"`
		SeatNumber int  `json:"seat_number" binding:"required,min=1"`
	}

	var req SeatAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var guest models.Guest
	if err := db.Where("id = ? AND wedding_id = ?", req.GuestID, weddingID).First(&guest).Error; err != nil {
		response.NotFound(c, "Guest not found")
		return
	}

	var table models.GuestTable
	if err := db.Where("id = ? AND wedding_id = ?", req.TableID, weddingID).First(&table).Error; err != nil {
		response.NotFound(c, "Table not found")
		return
	}

	db.Model(&models.Guest{}).Where("id = ?", req.GuestID).Updates(map[string]interface{}{
		"table_id":    req.TableID,
		"seat_number": req.SeatNumber,
	})

	response.Success(c, gin.H{"message": "Seat assigned successfully"})
}
