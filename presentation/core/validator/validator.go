package validator

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

func ValidateId(ctx context.Context, id int64) error {
	if id < 0 {
		errStr := "Id must be positive integer!"
		err := errors.New(errStr)
		return err
	}
	return nil
}

func ValidateIds(ctx context.Context, FromId, ToId int64) error {
	if FromId < 0 || ToId < 0 {
		errStr := "Id must be positive integer!"
		err := errors.New(errStr)
		return err
	} else if FromId == ToId {
		errStr := "Ids can’t be equal!"
		err := errors.New(errStr)
		return err
	}
	return nil
}

func ValidateBalanceForTransaction(ctx context.Context, balance *decimal.Decimal, value string) error {
	decimalValue := decimal.RequireFromString(value)
	if !balance.GreaterThanOrEqual(decimalValue) {
		errStr := "Not enough money for transaction!"
		err := errors.New(errStr)
		return err
	}
	return nil
}

func ValidateQueryKey(s, value string) string {
	if s == "" {
		s = value
	}
	return s
}
