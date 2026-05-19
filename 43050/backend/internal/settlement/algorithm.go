package settlement

import (
	"sort"
)

type Transfer struct {
	From   uint
	To     uint
	Amount float64
}

type Balance struct {
	UserID  uint
	Balance float64
}

func CalculateOptimalTransfers(balances []Balance) []Transfer {
	var debtors []Balance
	var creditors []Balance

	for _, b := range balances {
		if b.Balance < -0.01 {
			debtors = append(debtors, Balance{UserID: b.UserID, Balance: -b.Balance})
		} else if b.Balance > 0.01 {
			creditors = append(creditors, b)
		}
	}

	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].Balance > debtors[j].Balance
	})
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].Balance > creditors[j].Balance
	})

	var transfers []Transfer
	i, j := 0, 0

	for i < len(debtors) && j < len(creditors) {
		amount := min(debtors[i].Balance, creditors[j].Balance)
		if amount > 0.01 {
			transfers = append(transfers, Transfer{
				From:   debtors[i].UserID,
				To:     creditors[j].UserID,
				Amount: roundFloat(amount, 2),
			})
		}

		debtors[i].Balance -= amount
		creditors[j].Balance -= amount

		if debtors[i].Balance < 0.01 {
			i++
		}
		if creditors[j].Balance < 0.01 {
			j++
		}
	}

	return transfers
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func roundFloat(val float64, precision int) float64 {
	multiplier := 1.0
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}
	return float64(int(val*multiplier+0.5)) / multiplier
}
