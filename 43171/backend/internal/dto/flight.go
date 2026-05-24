package dto

type CreateFlightRecordReq struct {
	OrderID     *uint   `json:"order_id"`
	ServiceID   *uint   `json:"service_id"`
	DroneID     uint    `json:"drone_id" binding:"required"`
	StartPoint  string  `json:"start_point"`
	EndPoint    string  `json:"end_point"`
	Route       string  `json:"route"`
	AltitudeMax float64 `json:"altitude_max"`
	AltitudeAvg float64 `json:"altitude_avg"`
	Duration    int     `json:"duration"`
	Distance    float64 `json:"distance"`
	FlightDate  string  `json:"flight_date"`
	FlightLog   string  `json:"flight_log"`
	Images      string  `json:"images"`
	Remark      string  `json:"remark"`
}

type FlightQuery struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	DroneID   uint   `form:"drone_id"`
	PilotID   uint   `form:"pilot_id"`
	OrderID   uint   `form:"order_id"`
	ServiceID uint   `form:"service_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
