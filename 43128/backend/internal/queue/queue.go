package queue

import (
	"context"
	"encoding/json"
	"time"

	"event-platform/internal/database"
	applog "event-platform/internal/logger"
	"event-platform/internal/model"
	"event-platform/internal/repository"
)

const (
	ChannelMessages    = "event:messages"
	ChannelCertificates = "event:certificates"
)

type MessagePayload struct {
	UserIDs []uint          `json:"user_ids"`
	Type    model.MessageType `json:"type"`
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Extra   string          `json:"extra,omitempty"`
}

type CertPayload struct {
	CertificateID uint `json:"certificate_id"`
}

type Queue struct {
	msgRepo *repository.MessageRepo
	certRepo *repository.CertificateRepo
}

func New(msgRepo *repository.MessageRepo, certRepo *repository.CertificateRepo) *Queue {
	return &Queue{msgRepo: msgRepo, certRepo: certRepo}
}

func (q *Queue) PublishMessage(ctx context.Context, p MessagePayload) error {
	b, _ := json.Marshal(p)
	return database.RDB.Publish(ctx, ChannelMessages, b).Err()
}

func (q *Queue) PublishCertificate(ctx context.Context, p CertPayload) error {
	b, _ := json.Marshal(p)
	return database.RDB.Publish(ctx, ChannelCertificates, b).Err()
}

func (q *Queue) ConsumeMessages(ctx context.Context) {
	sub := database.RDB.Subscribe(ctx, ChannelMessages)
	defer sub.Close()
	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var p MessagePayload
			if err := json.Unmarshal([]byte(msg.Payload), &p); err != nil {
				applog.Errorf("consume message parse error: %v", err)
				continue
			}
			var list []*model.Message
			for _, uid := range p.UserIDs {
				list = append(list, &model.Message{
					UserID:  uid,
					Type:    p.Type,
					Title:   p.Title,
					Content: p.Content,
					Extra:   p.Extra,
				})
			}
			if len(list) > 0 {
				if err := q.msgRepo.BatchCreate(list); err != nil {
					applog.Errorf("batch create messages error: %v", err)
				} else {
					applog.Infof("batch created %d messages", len(list))
				}
			}
		}
	}
}

func (q *Queue) ConsumeCertificates(ctx context.Context, genFn func(id uint) error) {
	sub := database.RDB.Subscribe(ctx, ChannelCertificates)
	defer sub.Close()
	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var p CertPayload
			if err := json.Unmarshal([]byte(msg.Payload), &p); err != nil {
				applog.Errorf("consume cert parse error: %v", err)
				continue
			}
			if err := genFn(p.CertificateID); err != nil {
				applog.Errorf("generate cert %d error: %v", p.CertificateID, err)
				cert, _ := q.certRepo.GetByID(p.CertificateID)
				if cert != nil {
					cert.RetryCount++
					if cert.RetryCount < 3 {
						cert.Status = "failed"
						_ = q.certRepo.Update(cert)
						go func() {
							time.Sleep(time.Duration(5*cert.RetryCount) * time.Second)
							_ = q.PublishCertificate(context.Background(), CertPayload{CertificateID: cert.ID})
						}()
					} else {
						cert.Status = "failed"
						_ = q.certRepo.Update(cert)
					}
				}
			}
		}
	}
}
