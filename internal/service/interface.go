package service

import (
	"context"

	"github.com/ArtemQa8/currency-converter/internal/model"
)

type RateProvider interface {
	GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error)
}
