package exchangerate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"job/presentation/core/jsonint"
	"net/http"

	"github.com/shopspring/decimal"
)

type SavedRates jsonint.AllRatesJSON

var rates SavedRates

func (rates *SavedRates) getCurrencyRates() (*map[string]float64, error) {
	url := "https://api.exchangeratesapi.io/latest?base=RUB"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, err
	}

	return &rates.Rates, nil
}

func ExchangeCurrency(amount *decimal.Decimal, currency string) (*decimal.Decimal, error) {
	if rates.Rates == nil {
		_, err := rates.getCurrencyRates()
		if err != nil {
			return nil, err
		}
	}

	current, ok := rates.Rates[currency]
	if !ok {
		err := errors.New("currency doesn't exist")
		return nil, err
	}
	decimalCurrent := decimal.NewFromFloatWithExponent(current, -6)
	amountInCurrency := amount.Mul(decimalCurrent)
	roundedAmount := amountInCurrency.RoundCash(10)

	return &roundedAmount, nil
}

func Get() *map[string]float64 {
	return &rates.Rates
}
