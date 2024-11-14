// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

func (e *Currency) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Currency(s)
	case string:
		*e = Currency(s)
	default:
		return fmt.Errorf("unsupported scan type for Currency: %T", src)
	}
	return nil
}

type NullCurrency struct {
	Currency Currency `json:"currency"`
	Valid    bool     `json:"valid"` // Valid is true if Currency is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCurrency) Scan(value interface{}) error {
	if value == nil {
		ns.Currency, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Currency.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCurrency) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Currency), nil
}

type Account struct {
	ID        int64        `json:"id"`
	Owner     string       `json:"owner"`
	Balance   int64        `json:"balance"`
	Currency  Currency     `json:"currency"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type Entry struct {
	ID int64 `json:"id"`
	// can be positive/negative
	Amount    int64     `json:"amount"`
	AccountID int64     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID int64 `json:"id"`
	// must be positive
	Amount    int64        `json:"amount"`
	Currency  Currency     `json:"currency"`
	FromAccID int64        `json:"from_acc_id"`
	ToAccID   int64        `json:"to_acc_id"`
	CreatedAt sql.NullTime `json:"created_at"`
}