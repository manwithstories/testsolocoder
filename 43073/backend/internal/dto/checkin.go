package dto

type CheckInRequest struct {
	QrCode string `json:"qrCode" validate:"required"`
}

type CheckInListRequest struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"pageSize,default=10"`
	ActivityID uint   `form:"activityId"`
	CheckedIn *bool  `form:"checkedIn"`
	Keyword  string `form:"keyword"`
}
