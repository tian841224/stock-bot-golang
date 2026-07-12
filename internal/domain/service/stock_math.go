package logic

func CalculateStockPerformance(todayClosePrice, prevClosePrice float64) (float64, float64, string) {
	changeAmount := todayClosePrice - prevClosePrice
	changeRate := (changeAmount / prevClosePrice) * 100
	upDownSign := ""
	if changeAmount > 0 {
		upDownSign = "+"
	} else if changeAmount < 0 {
		upDownSign = "-"
	}

	return changeAmount, changeRate, upDownSign
}
