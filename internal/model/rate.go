package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Rate struct {
	BaseCurrency  string
	QuoteCurrency string
	Value         decimal.Decimal
	Date          time.Time
	Source        string
}
