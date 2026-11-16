package service

import (
	"context"

	"github.com/ArtemQa8/currency-converter/internal/model"
)

type RateProvider interface {
	GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error)
}

type RateSaver interface {
	SaveRate(ctx context.Context, rate model.Rate) error
}

type AllRatesProvider interface {
	GetAllRates(ctx context.Context) (map[string]model.Rate, error)
}
