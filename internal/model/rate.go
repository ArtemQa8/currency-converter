package model

//nolint:depguard
import (
	"time"

	"github.com/shopspring/decimal"
)

// Rate описывает модель курса валют.
type Rate struct {
	BaseCurrency  string
	QuoteCurrency string
	Value         decimal.Decimal
	Date          time.Time
	Source        string
}
