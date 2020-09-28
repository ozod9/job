package models

import (
	"job/presentation/core/mytime"

	"github.com/shopspring/decimal"
)

type Balance struct {
	ID     *int64           `json:"id"`
	Amount *decimal.Decimal `json:"amount"`
}

type Transaction struct {
	ID        *int64           `json:"id"`
	BalanceID *int64           `json:"balance_id"`
	FromID    *int64           `json:"from_id"`
	Amount    *decimal.Decimal `json:"amount"`
	Reason    *string          `json:"reason"`
	Type      *string          `json:"type"`
	Date      *mytime.MyTime   `json:"date"`
}
