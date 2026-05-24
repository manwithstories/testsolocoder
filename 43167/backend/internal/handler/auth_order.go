package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"watchplatform/internal/app"
	"watchplatform/internal/config"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type AuthOrderCreateReq struct {
	WatchID uint   `json:"watch_id"`
	Note    string `json:"note"`
}

func CreateAuthOrder(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req AuthOrderCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	order := model.AuthOrder{
		UserID: u.ID,
		WatchID: req.WatchID,
		Status: model.AuthPending,
		Note:   req.Note,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	form, err := c.MultipartForm()
	if err == nil && len(form.File["photos"]) > 0 {
		for _, f := range form.File["photos"] {
			if f.Size > config.Cfg.MaxFileSize {
				continue
			}
			ext := strings.ToLower(filepath.Ext(f.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				continue
			}
			dir := filepath.Join(config.Cfg.UploadDir, "auth", strconv.FormatUint(uint64(order.ID), 10))
			name := uuid.New().String() + ext
			dst := filepath.Join(dir, name)
			if err := c.SaveUploadedFile(f, dst); err != nil {
				continue
			}
			url := "/uploads/auth/" + strconv.FormatUint(uint64(order.ID), 10) + "/" + name
			_ = database.DB.Create(&model.AuthPhoto{AuthOrderID: order.ID, URL: url})
		}
	}
	pushMessage(u.ID, "new_auth_order", "新的鉴定申请", "您提交了一个新的鉴定申请", "auth_order", order.ID)
	app.OK(c, order)
}

func ListAuthOrders(c *gin.Context) {
	u := middleware.CurrentUser(c)
	db := database.DB.Model(&model.AuthOrder{})
	if u.Role == model.RoleAppraiser {
		if c.Query("mine") == "1" {
			db = db.Where("appraiser_id = ?", u.ID)
		}
	} else {
		db = db.Where("user_id = ?", u.ID)
	}
	var total int64
	db.Count(&total)
	var list []model.AuthOrder
	db.Preload("Photos").Preload("Report").Order("created_at desc").Find(&list)
	app.OK(c, gin.H{"total": total, "list": list})
}

func GetAuthOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var o model.AuthOrder
	if err := database.DB.Preload("Photos").Preload("Report").First(&o, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "鉴定申请不存在")
		return
	}
	app.OK(c, o)
}

func AssignAuthOrder(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var o model.AuthOrder
	if err := database.DB.First(&o, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "鉴定申请不存在")
		return
	}
	if o.Status != model.AuthPending {
		app.Fail(c, http.StatusBadRequest, "状态不允许")
		return
	}
	if err := database.DB.Model(&o).Updates(map[string]interface{}{
		"appraiser_id": u.ID,
		"status":       model.AuthAssigned,
	}).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(o.UserID, "auth_assigned", "鉴定已分配", "您的鉴定申请已被受理", "auth_order", o.ID)
	app.OK(c, nil)
}

type AuthReportReq struct {
	Conclusion     string  `json:"conclusion" binding:"required,oneof=authentic fake uncertain"`
	Authentic      bool    `json:"authentic"`
	Details        string  `json:"details"`
	EstimatedValue float64 `json:"estimated_value"`
}

func SubmitAuthReport(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var o model.AuthOrder
	if err := database.DB.First(&o, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "鉴定申请不存在")
		return
	}
	if o.AppraiserID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权审核")
		return
	}
	if o.Status != model.AuthAssigned {
		app.Fail(c, http.StatusBadRequest, "状态不允许")
		return
	}
	var req AuthReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	var report model.AuthReport
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		report = model.AuthReport{
			AuthOrderID:    o.ID,
			AppraiserID:    u.ID,
			Conclusion:     req.Conclusion,
			Authentic:      req.Authentic,
			Details:        req.Details,
			EstimatedValue: req.EstimatedValue,
		}
		if err := tx.Create(&report).Error; err != nil {
			return err
		}
		pdfName := "report_" + strconv.FormatUint(uint64(report.ID), 10) + ".pdf"
		pdfDir := filepath.Join(config.Cfg.UploadDir, "reports")
		_ = os.MkdirAll(pdfDir, 0o755)
		pdfPath := filepath.Join(pdfDir, pdfName)
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.Cell(0, 12, "Watch Authentication Report")
		pdf.Ln(12)
		pdf.SetFont("Helvetica", "", 12)
		pdf.Cell(0, 10, "Report ID: "+strconv.FormatUint(uint64(report.ID), 10))
		pdf.Ln(8)
		pdf.Cell(0, 10, "Conclusion: "+req.Conclusion)
		pdf.Ln(8)
		pdf.Cell(0, 10, "Estimated Value: "+strconv.FormatFloat(req.EstimatedValue, 'f', 2, 64))
		pdf.Ln(8)
		pdf.MultiCell(0, 8, "Details:\n"+req.Details, "", "", false)
		if err := pdf.OutputFileAndClose(pdfPath); err != nil {
			return err
		}
		if err := tx.Model(&report).Update("pdf_path", "/uploads/reports/"+pdfName).Error; err != nil {
			return err
		}
		if err := tx.Model(&o).Update("status", model.AuthReported).Error; err != nil {
			return err
		}
		if o.WatchID != 0 && req.Authentic {
			if err := tx.Model(&model.Watch{}).Where("id = ?", o.WatchID).Update("authed", true).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(o.UserID, "auth_report", "鉴定报告已出具", "您的鉴定申请已出具报告", "auth_order", o.ID)
	app.OK(c, report)
}

func DownloadReportPDF(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var o model.AuthOrder
	if err := database.DB.Preload("Report").First(&o, id).Error; err != nil || o.Report == nil {
		app.Fail(c, http.StatusNotFound, "报告不存在")
		return
	}
	if o.Report.PDFPath == "" {
		app.Fail(c, http.StatusNotFound, "PDF 未生成")
		return
	}
	full := filepath.Join(config.Cfg.UploadDir, strings.TrimPrefix(o.Report.PDFPath, "/uploads/"))
	c.FileAttachment(full, "report.pdf")
}
