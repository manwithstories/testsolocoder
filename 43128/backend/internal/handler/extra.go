package handler

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"event-platform/internal/database"
	"event-platform/internal/model"
	"event-platform/internal/queue"
	"event-platform/internal/service"
	"event-platform/pkg/excel"
	"event-platform/pkg/pdf"
	"event-platform/pkg/response"
	"event-platform/pkg/retry"

	"github.com/gin-gonic/gin"
)

type ScoreImportHandler struct {
	svc *service.ScoreService
}

func NewScoreImportHandler(svc *service.ScoreService) *ScoreImportHandler {
	return &ScoreImportHandler{svc: svc}
}

func (h *ScoreImportHandler) Template(c *gin.Context) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="score_template.xlsx"`)
	var buf bytes.Buffer
	if err := excel.GenerateScoreTemplate(&buf); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func (h *ScoreImportHandler) Import(c *gin.Context) {
	eid, _ := strconv.ParseUint(c.PostForm("event_id"), 10, 64)
	iid, _ := strconv.ParseUint(c.PostForm("event_item_id"), 10, 64)
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "缺少文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		response.BadRequest(c, "读取文件失败")
		return
	}
	defer f.Close()
	rows, err := excel.ParseScoreFile(f)
	if err != nil {
		response.BadRequest(c, "解析Excel失败: "+err.Error())
		return
	}
	var list []model.Score
	for _, r := range rows {
		list = append(list, model.Score{
			EventID:     uint(eid),
			EventItemID: uint(iid),
			UserID:      r.UserID,
			Score:       r.Score,
			TimeUsed:    r.TimeUsed,
			Remarks:     r.Remarks,
			IsValid:     true,
		})
	}
	if err := h.svc.Entry(list); err != nil {
		response.Fail(c, 400, 40001, err.Error())
		return
	}
	response.OK(c, gin.H{"imported": len(list)})
}

type CertificateHandler struct {
	certSvc  *service.CertificateService
	scoreSvc *service.ScoreService
	userSvc  *service.UserService
	eventSvc *service.EventService
	msgQueue *queue.Queue
	certDir  string
}

func NewCertificateHandler(
	certSvc *service.CertificateService,
	scoreSvc *service.ScoreService,
	userSvc *service.UserService,
	eventSvc *service.EventService,
	msgQueue *queue.Queue,
	certDir string,
) *CertificateHandler {
	return &CertificateHandler{certSvc: certSvc, scoreSvc: scoreSvc, userSvc: userSvc, eventSvc: eventSvc, msgQueue: msgQueue, certDir: certDir}
}

func (h *CertificateHandler) MyList(c *gin.Context) {
	uid, _ := c.Get("user_id")
	list, err := h.certSvc.ListByUser(uid.(uint))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *CertificateHandler) Generate(c *gin.Context) {
	scoreID, _ := strconv.ParseUint(c.Param("score_id"), 10, 64)
	_ = scoreID
	response.OK(c, nil)
}

func (h *CertificateHandler) Download(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	cert, err := h.certSvc.Get(uint(id))
	if err != nil || cert == nil || cert.FilePath == "" {
		response.NotFound(c, "证书不存在或未生成")
		return
	}
	c.File(cert.FilePath)
}

func (h *CertificateHandler) GenerateCert(certID uint) error {
	return retry.Do(func() error {
		cert, err := h.certSvc.Get(certID)
		if err != nil {
			return err
		}
		user, _ := h.userSvc.GetByID(cert.UserID)
		event, _ := h.eventSvc.Get(cert.EventID)
		item, _ := h.eventSvc.Item(cert.EventItemID)
		name := "参与者"
		eventName := ""
		itemName := ""
		if user != nil {
			name = user.RealName
		}
		if event != nil {
			eventName = event.Name
		}
		if item != nil {
			itemName = item.Name
		}
		path, err := pdf.GenerateCertificate(h.certDir, cert.CertificateNo, eventName, itemName, name, cert.Rank, cert.Score)
		if err != nil {
			return err
		}
		cert.FilePath = path
		now := time.Now()
		cert.GeneratedAt = &now
		cert.Status = "generated"
		return h.certSvc.Update(cert)
	}, 3, 2*time.Second)
}

type StatsHandler struct {
	regSvc   *service.RegistrationService
	scoreSvc *service.ScoreService
	eventSvc *service.EventService
}

func NewStatsHandler(
	regSvc *service.RegistrationService,
	scoreSvc *service.ScoreService,
	eventSvc *service.EventService,
) *StatsHandler {
	return &StatsHandler{regSvc: regSvc, scoreSvc: scoreSvc, eventSvc: eventSvc}
}

func (h *StatsHandler) Overview(c *gin.Context) {
	cacheKey := "stats:overview"
	if database.RDB != nil {
		if v, err := database.RDB.Get(c.Request.Context(), cacheKey).Result(); err == nil {
			c.Data(200, "application/json", []byte(v))
			return
		}
	}
	eid, _ := strconv.ParseUint(c.Query("event_id"), 10, 64)
	resp := gin.H{
		"event_id":     eid,
		"generated_at": time.Now().Format(time.RFC3339),
	}
	if database.RDB != nil {
		b, _ := json.Marshal(resp)
		database.RDB.Set(c.Request.Context(), cacheKey, b, 5*time.Minute)
	}
	response.OK(c, resp)
}

func (h *StatsHandler) Export(c *gin.Context) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="stats.xlsx"`)
	var buf bytes.Buffer
	headers := []string{"赛事", "项目", "参赛人数", "平均分"}
	rows := [][]interface{}{
		{"示例赛事", "示例项目", 100, 85.5},
	}
	if err := excel.ExportStats(&buf, headers, rows); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
