package routes

import (
	"job/presentation/controller"

	"job/domain/models"
	"job/presentation/core/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(env *controller.Environment, conf *models.Config) (*mux.Router, error) {
	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/balances/{id}", middleware.Requests(env.GetBalance)).Methods("GET")
	r.HandleFunc("/balances/history/{id}", middleware.Requests(env.GetHistory)).Methods("GET")
	r.HandleFunc("/balances/transfer", middleware.Requests(env.TransferTransaction)).Methods("POST")
	r.HandleFunc("/balances/income", middleware.Requests(env.IncomeTransaction)).Methods("POST")
	r.HandleFunc("/balances/outcome", middleware.Requests(env.OutcomeTransaction)).Methods("POST")

	return r, nil
}
