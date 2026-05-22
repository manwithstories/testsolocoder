package utils

import "math"

func RoundToTwoDecimals(price float64) float64 {
	return math.Round(price*100) / 100
}

func RoundPrice(price float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(price*pow) / pow
}

func CalculateTotal(basePrice float64, durationMinutes int) float64 {
	hours := float64(durationMinutes) / 60.0
	return RoundToTwoDecimals(basePrice * hours)
}

func CalculatePlatformFee(totalAmount float64, commissionRate float64) float64 {
	return RoundToTwoDecimals(totalAmount * commissionRate)
}

func CalculateProviderIncome(totalAmount float64, platformFee float64) float64 {
	return RoundToTwoDecimals(totalAmount - platformFee)
}

func CalculatePenalty(totalAmount float64, penaltyRate float64) float64 {
	return RoundToTwoDecimals(totalAmount * penaltyRate)
}
