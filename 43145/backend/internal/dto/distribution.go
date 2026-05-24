package dto

import "time"

type CreateDistributionLinkRequest struct {
	Channel   string     `json:"channel" binding:"required,oneof=email qrcode wechat dm api general"`
	MaxUses   int        `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type SendInvitationsRequest struct {
	Emails    []string   `json:"emails" binding:"required,min=1"`
	ExpiresAt *time.Time `json:"expires_at"`
	CustomMessage string `json:"custom_message"`
}

type ExportReportRequest struct {
	SurveyID   uint       `json:"survey_id" binding:"required"`
	Format     string     `json:"format" binding:"required,oneof=pdf excel image"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	IncludeRaw bool       `json:"include_raw"`
}

type DistributionLinkResponse struct {
	ID        uint       `json:"id"`
	LinkToken string     `json:"link_token"`
	Channel   string     `json:"channel"`
	MaxUses   int        `json:"max_uses"`
	UseCount  int        `json:"use_count"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsActive  bool       `json:"is_active"`
	FullURL   string     `json:"full_url"`
	QRCodeURL string     `json:"qrcode_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
