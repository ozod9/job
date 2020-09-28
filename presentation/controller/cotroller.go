package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"job/application/exchangerate"

	"job/domain/repository"

	"github.com/jimlawless/whereami"

	"job/presentation/core/jsonint"
	"job/presentation/core/rfc7807"

	"job/presentation/core/validator"

	"github.com/gorilla/mux"
)

func (env *Environment) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		errStr := "Id must be positive integer!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		env.logger.Info(errStr, whereami.WhereAmI())
		if err != nil {
			return
		}
		return
	}
	if err := validator.ValidateId(ctx, id); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	user, err := repository.GetBalancePg(ctx, env.Balances, int64(id))
	if user.ID == nil {
		err = errors.New("Have no balance with that id!")
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	currency := "RUB"
	value := user.Amount

	keys, ok := r.URL.Query()["currency"]
	if ok && len(keys[0]) > 1 {
		currency = keys[0]
		value, err = exchangerate.ExchangeCurrency(value, currency)
		if err != nil {
			errStr := "Url Param 'currency' is not allowable! Have to use existing currency parameter values!"
			problem := rfc7807.NewProblem().
				AppendError("Id", errStr).
				SetType("business").
				SetStatus(http.StatusBadRequest)
			env.logger.Info(err.Error(), whereami.WhereAmI())
			err = problem.Write(w)
			if err != nil {
				return
			}
			return
		}
	}

	result := fmt.Sprintf("%s %s", value, currency)
	body, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Println(err)
		return
	}
}

func (env *Environment) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		errStr := "Id must be positive integer!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		env.logger.Info(errStr, whereami.WhereAmI())
		if err != nil {
			return
		}
		return
	}
	if err := validator.ValidateId(ctx, id); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	keys := r.URL.Query()
	order_by := validator.ValidateQueryKey(keys.Get("order_by"), "id")
	limit := validator.ValidateQueryKey(keys.Get("limit"), "null")
	offset := validator.ValidateQueryKey(keys.Get("offset"), "null")

	transactions, err := repository.GetHistoryPg(ctx, env.Balances, int64(id), order_by, limit, offset)
	if len(transactions) < 1 {
		err = errors.New("Have no user or transactions with that id!")
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	body, err := json.Marshal(transactions)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Println(err)
		return
	}
}

func (env *Environment) TransferTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transaction := jsonint.TransactionJSON{}
	err := jsonint.BodyToJSON(r.Body, &transaction)
	if err != nil {
		errStr := "Id must be positive integer! Reason and amount must be decimal in string!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.FromId.Valid {
		errStr := "Id must be positive integer, not null!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.ToId.Valid {
		errStr := "Id must be positive integer, not null!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.Amount.Valid {
		errStr := "Amount must be positive integer in string!"
		problem := rfc7807.NewProblem().
			AppendError("Amount", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.Reason.Valid {
		errStr := "Reason must be string!"
		problem := rfc7807.NewProblem().
			AppendError("Amount", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err := validator.ValidateIds(ctx, transaction.FromId.Value, transaction.ToId.Value); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	user, err := repository.GetBalancePg(ctx, env.Balances, transaction.FromId.Value)
	if user.ID == nil {
		err = errors.New("Have no balance with that id!")
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	if err := validator.ValidateBalanceForTransaction(ctx, user.Amount, transaction.Amount.Value); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Balance", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	err = repository.TransferTransactionPg(ctx, env.Balances, transaction)
	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	body, err := json.Marshal("Done!")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Println(err)
		return
	}
}

func (env *Environment) IncomeTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transaction := jsonint.TransactionJSON{}
	err := jsonint.BodyToJSON(r.Body, &transaction)
	if err != nil {
		errStr := "Id must be positive integer! Amount must be decimal in string!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.ToId.Valid {
		errStr := "Id must be positive integer, not null!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.Amount.Valid {
		errStr := "Amount must be positive integer!"
		problem := rfc7807.NewProblem().
			AppendError("Amount", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err := validator.ValidateId(ctx, transaction.ToId.Value); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}
	transaction.Type.Value = "income"
	transaction.FromId.Value = 0

	err = repository.IncomeTransactionPg(ctx, env.Balances, transaction)
	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	body, err := json.Marshal("Done!")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Println(err)
		return
	}
}

func (env *Environment) OutcomeTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transaction := jsonint.TransactionJSON{}
	err := jsonint.BodyToJSON(r.Body, &transaction)
	if err != nil {
		errStr := "Id must be positive integer! Value must be decimal in string "
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.FromId.Valid {
		errStr := "Id must be positive integer, not null!"
		problem := rfc7807.NewProblem().
			AppendError("Id", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if !transaction.Amount.Valid {
		errStr := "Amount must be positive integer!"
		problem := rfc7807.NewProblem().
			AppendError("Amount", errStr).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err := validator.ValidateId(ctx, transaction.FromId.Value); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	user, err := repository.GetBalancePg(ctx, env.Balances, transaction.FromId.Value)
	if user.ID == nil {
		err = errors.New("Have no balance with that id!")
		problem := rfc7807.NewProblem().
			AppendError("Id", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	if err := validator.ValidateBalanceForTransaction(ctx, user.Amount, transaction.Amount.Value); err != nil {
		problem := rfc7807.NewProblem().
			AppendError("Balance", err.Error()).
			SetType("business").
			SetStatus(http.StatusBadRequest)
		env.logger.Info(err.Error(), whereami.WhereAmI())
		err = problem.Write(w)
		if err != nil {
			return
		}
		return
	}

	transaction.Type.Value = "outcome"
	transaction.ToId.Value = 0

	err = repository.OutcomeTransactionPg(ctx, env.Balances, transaction)
	if err != nil {
		log.Println(err)
		env.logger.Error(err.Error(), whereami.WhereAmI())
		return
	}

	body, err := json.Marshal("Done!")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Println(err)
		return
	}
}
