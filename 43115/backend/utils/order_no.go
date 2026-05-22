package utils

import (
	"fmt"
	"time"
)

func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("HK%04d%02d%02d%02d%02d%02d%06d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond()/1000%1000000,
	)
}

func GenerateTransactionNo() string {
	now := time.Now()
	return fmt.Sprintf("TX%04d%02d%02d%02d%02d%02d%06d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond()/1000%1000000,
	)
}
