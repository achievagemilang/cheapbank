package util

import (
	"errors"
	"fmt"
)

var rates = map[string]float64{
	"USD": 1.0,
	"EUR": 0.85,
	"AUD": 1.35,
	"CAD": 1.25,
}

func ConvertCurrency(amount float64, from string, to string) (float64, error) {
	fromRate, fromExists := rates[from]
	toRate, toExists := rates[to]

	fmt.Println("debug", from, to, fromRate, toRate)

	if !fromExists || !toExists {
		return 0, errors.New("invalid currency code")
	}

	convertedAmount := (amount / fromRate) * toRate
	return convertedAmount, nil
}

