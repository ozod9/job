package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"job/domain/models"
	"job/presentation/core/jsonint"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewPgDatabase(DB_USER, DB_PASSWORD, DB_HOST, DB_NAME string, DB_PORT int) (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

func GetBalancePg(ctx context.Context, db *sql.DB, id int64) (*models.Balance, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM balances WHERE id = $1", id)
	if err != nil {
		if err == ctx.Err() {
			return nil, errors.New("request cancel")
		}
		return nil, err
	}
	var balance = models.Balance{}

	for rows.Next() {
		var balanceDTO models.BalanceDTO
		if err := rows.Scan(&balanceDTO.ID, &balanceDTO.Amount); err != nil {
			return nil, err
		}
		balance = balanceDTO.GetEntity()
	}

	return &balance, nil
}

func GetHistoryPg(ctx context.Context, db *sql.DB, userId int64, order_by, limit, offset string) ([]models.Transaction, error) {
	queryString := fmt.Sprintf("SELECT * FROM transactions WHERE balance_id = $1 ORDER BY %s LIMIT %s OFFSET %s;", order_by, limit, offset)
	rows, err := db.QueryContext(ctx, queryString, userId)
	if err != nil {
		if err == ctx.Err() {
			return nil, errors.New("request cancel")
		}
		return nil, err
	}

	var transactions = make([]models.Transaction, 0)

	for rows.Next() {
		var transaction models.TransactionDTO
		if err := rows.Scan(&transaction.ID, &transaction.BalanceID, &transaction.FromID, &transaction.Amount, &transaction.Reason, &transaction.Type, &transaction.Date); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction.GetEntity())
	}

	return transactions, nil
}

func TransferTransactionPg(ctx context.Context, db *sql.DB, transaction jsonint.TransactionJSON) error {
	err := OutcomeTransactionPg(ctx, db, transaction)
	if err != nil {
		return err
	}

	err = IncomeTransactionPg(ctx, db, transaction)
	if err != nil {
		return err
	}

	return nil
}

func IncomeTransactionPg(ctx context.Context, db *sql.DB, transaction jsonint.TransactionJSON) error {
	queryString := `INSERT INTO balances (id, balance)
					VALUES ($1, $2)
					ON CONFLICT (id) DO UPDATE SET balance = balances.balance + EXCLUDED.balance;`

	res, err := db.Exec(queryString, transaction.ToId.Value, transaction.Amount.Value)
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	transaction.Type.Value = "income"

	err = AddTransactionInformationPg(ctx, db, transaction)
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	if r == 0 {
		err = errors.New("Have no balance with that id!")
		return err
	}

	return nil
}

func OutcomeTransactionPg(ctx context.Context, db *sql.DB, transaction jsonint.TransactionJSON) error {
	queryString := `UPDATE balances SET balance = balance - $1 WHERE id = $2;`
	res, err := db.Exec(queryString, transaction.Amount.Value, transaction.FromId.Value)
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	transaction.Type.Value = "outcome"
	err = AddTransactionInformationPg(ctx, db, transaction)
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	if r == 0 {
		err = errors.New("Have no user with that id!")
		return err
	}

	return nil
}

func AddTransactionInformationPg(ctx context.Context, db *sql.DB, transaction jsonint.TransactionJSON) error {
	queryString := `INSERT INTO transactions(balance_id, from_id, amount, reason, type, date) 
	VALUES ($1, $2, $3, $4, $5, $6);`

	var balance_id, from_id int64

	switch transaction.Type.Value {
	case "outcome":
		balance_id = transaction.FromId.Value
		from_id = transaction.ToId.Value
	case "income":
		balance_id = transaction.ToId.Value
		from_id = transaction.FromId.Value
	}

	res, err := db.Exec(queryString, balance_id, from_id, transaction.Amount.Value, transaction.Reason.Value, transaction.Type.Value, time.Now())
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		if err == ctx.Err() {
			return errors.New("request cancel")
		}
		return err
	}

	if r == 0 {
		err = errors.New("Have no balancr with that id!")
		return err
	}

	return nil
}
