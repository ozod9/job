package controller

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

type mockLogger struct{}

func (logger *mockLogger) Info(message string, source string) { log.Println(message) }

func (logger *mockLogger) Error(message string, source string) { log.Println(message) }

func (logger *mockLogger) Warning(message string, source string) { log.Println(message) }

func (logger *mockLogger) Fatal(message string, source string) { log.Println(message) }

func newLogger() *mockLogger {
	return new(mockLogger)
}

func Test_GetBalancePg_ShouldReturn_SuccessResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT (.+) FROM balances WHERE").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, "100"))

	env := &Environment{Balances: db, logger: newLogger()}

	vars := map[string]string{
		"id": "1",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/balances/{id}", nil)
	if err != nil {
		log.Println(err)
		return
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.GetBalance)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		log.Printf("Expected 200, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}
	expected := "application/json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expected {
		log.Printf("Expected application/json, but got %s\n", ctype)
		t.Fatal(ctype)
	}
}

func Test_GetBalancePg_ShouldReturn_ErrorResultBalance1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT (.+) FROM balances WHERE").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 100))

	env := &Environment{Balances: db, logger: newLogger()}

	vars := map[string]string{
		"id": "-w",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/balances/{id}", nil)
	if err != nil {
		log.Println(err)
		return
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.GetBalance)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		log.Printf("Expected 400, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}

	expected := "application/problem+json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expected {
		log.Printf("Expected application/problem+json, but got %s\n", ctype)
		t.Fatal(ctype)
	}

}

func Test_GetBalancePg_ShouldReturn_ErrorResultBalance2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM balances WHERE").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(10, 100))

	env := &Environment{Balances: db, logger: newLogger()}

	vars := map[string]string{
		"id": "-2",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/balances/{id}", nil)
	if err != nil {
		log.Println(err)
		return
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.GetBalance)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		log.Printf("Expected 400, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}

	expected := "application/problem+json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expected {
		log.Printf("Expected application/problem+json, but got %s\n", ctype)
		t.Fatal(ctype)
	}

}

func Test_GetHistoryPg_ShouldReturn_SuccessResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM transactions WHERE balance_id = (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+);").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "balance_id", "from_id", "amount", "reason", "type", "date"}).AddRow(1, "100", 0, "100", "Some", "income", time.Now()))

	env := &Environment{Balances: db, logger: newLogger()}

	vars := map[string]string{
		"id": "1",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/balances/history/{id}", nil)
	if err != nil {
		log.Println(err)
		return
	}

	req = mux.SetURLVars(req, vars)

	q := req.URL.Query()
	q.Add("order_by", "date")
	q.Add("limit", "5")
	q.Add("offset", "null")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.GetHistory)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		log.Printf("Expected 200, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}
	expected := "application/json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expected {
		log.Printf("Expected application/json, but got %s\n", ctype)
		t.Fatal(ctype)
	}
}

func Test_GetHistoryPg_ShouldReturn_ErrorResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM transactions WHERE balance_id = (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+);").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "balance_id", "from_id", "amount", "reason", "type", "date"}).AddRow(1, "100", 0, "100", "Some", "income", time.Now()))

	env := &Environment{Balances: db, logger: newLogger()}

	vars := map[string]string{
		"id": "-1",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/balances/history/{id}", nil)
	if err != nil {
		log.Println(err)
		return
	}

	req = mux.SetURLVars(req, vars)

	q := req.URL.Query()
	q.Add("order_by", "date")
	q.Add("limit", "5")
	q.Add("offset", "null")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.GetBalance)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		log.Printf("Expected 400, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}

	expected := "application/problem+json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expected {
		log.Printf("Expected application/problem+json, but got %s\n", ctype)
		t.Fatal(ctype)
	}
}

func Test_IncomeTransactionPg_ShouldReturn_SuccessResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO balances (.+) VALUES (.+) ON CONFLICT (id) DO UPDATE SET (.+);").WithArgs(2)

	env := &Environment{Balances: db, logger: newLogger()}

	var jsonStr = []byte(`{"toId": 1, "amount":"200", "reason":"Some"}`)
	
	req, err := http.NewRequest("POST", "http://localhost:8080/balances/income", bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.IncomeTransaction)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		log.Printf("Expected 200, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}
}


func Test_IncomeTransactionPg_ShouldReturn_ErrorResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO balances (.+) VALUES (.+) ON CONFLICT (id) DO UPDATE SET (.+);").WithArgs(2)

	env := &Environment{Balances: db, logger: newLogger()}

	var jsonStr = []byte(`{"toId": "1", "amount":"200", "reason":"Some"}`)
	
	req, err := http.NewRequest("POST", "http://localhost:8080/balances/income", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.IncomeTransaction)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		log.Printf("Expected 400, but got %d\n", rr.Code)
		t.Fatal(rr.Code)
	}
}
