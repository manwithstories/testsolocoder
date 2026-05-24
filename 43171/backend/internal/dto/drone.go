package dto

import "drone-rental/internal/model"

type CreateDroneReq struct {
	Name          string  `json:"name" binding:"required"`
	Brand         string  `json:"brand"`
	Model         string  `json:"model"`
	SerialNo      string  `json:"serial_no" binding:"required"`
	Weight        float64 `json:"weight"`
	BatteryLife   int     `json:"battery_life"`
	GimbalSpec    string  `json:"gimbal_spec"`
	CameraSpec    string  `json:"camera_spec"`
	MaxSpeed      float64 `json:"max_speed"`
	MaxAltitude   float64 `json:"max_altitude"`
	Region        string  `json:"region" binding:"required"`
	Description   string  `json:"description"`
	Images        string  `json:"images"`
	PricePerDay   float64 `json:"price_per_day" binding:"required,gt=0"`
	Deposit       float64 `json:"deposit"`
	AvailableFrom string  `json:"available_from"`
	AvailableTo   string  `json:"available_to"`
}

type UpdateDroneReq struct {
	Name          string  `json:"name"`
	Brand         string  `json:"brand"`
	Model         string  `json:"model"`
	Weight        float64 `json:"weight"`
	BatteryLife   int     `json:"battery_life"`
	GimbalSpec    string  `json:"gimbal_spec"`
	CameraSpec    string  `json:"camera_spec"`
	MaxSpeed      float64 `json:"max_speed"`
	MaxAltitude   float64 `json:"max_altitude"`
	Region        string  `json:"region"`
	Description   string  `json:"description"`
	Images        string  `json:"images"`
	PricePerDay   float64 `json:"price_per_day"`
	Deposit       float64 `json:"deposit"`
	AvailableFrom string  `json:"available_from"`
	AvailableTo   string  `json:"available_to"`
}

type DroneQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	Keyword    string `form:"keyword"`
	Region     string `form:"region"`
	Brand      string `form:"brand"`
	Status     string `form:"status"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
}

type BatchImportReq struct {
	Drones []CreateDroneReq `json:"drones" binding:"required"`
}

type UpdateStatusReq struct {
	Status model.DroneStatus `json:"status" binding:"required,oneof=offline online rented maintenance"`
}
