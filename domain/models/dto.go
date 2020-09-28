package models

import (
	"database/sql"
	"job/presentation/core/mytime"
	"time"

	"github.com/shopspring/decimal"
)

type BalanceDTO struct {
	ID     sql.NullInt64
	Amount sql.NullString
}

type TransactionDTO struct {
	ID        sql.NullInt64
	BalanceID sql.NullInt64
	FromID    sql.NullInt64
	Amount    sql.NullString
	Reason    sql.NullString
	Type      sql.NullString
	Date      time.Time
}

func (user BalanceDTO) GetEntity() Balance {
	return Balance{
		ID:     getInt64Pointer(user.ID),
		Amount: getDecimalPointer(user.Amount),
	}
}

func (transaction TransactionDTO) GetEntity() Transaction {
	return Transaction{
		ID:        getInt64Pointer(transaction.ID),
		BalanceID: getInt64Pointer(transaction.BalanceID),
		FromID:    getInt64Pointer(transaction.FromID),
		Amount:    getDecimalPointer(transaction.Amount),
		Reason:    getStringPointer(transaction.Reason),
		Type:      getStringPointer(transaction.Type),
		Date:      getTimePointer(transaction.Date),
	}
}

func getTimePointer(time time.Time) *mytime.MyTime {
	if time.IsZero() {
		return nil
	}
	return &mytime.MyTime{&time}
}

func getInt64Pointer(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}
	return nil
}

func getStringPointer(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

func getDecimalPointer(value sql.NullString) *decimal.Decimal {
	if value.Valid {
		n := decimal.RequireFromString(value.String)
		return &n
	}
	return nil
}
