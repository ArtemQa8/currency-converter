package service

//nolint:depguard
import (
	"context"

	"github.com/ArtemQa8/currency-converter/internal/model"
)

// RateProvider описывает контракт для получения курса валют.
type RateProvider interface {
	GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error)
}

// RateSaver описывает контракт для сохранения курса валют.
type RateSaver interface {
	SaveRate(ctx context.Context, rate model.Rate) error
}

// AllRatesProvider описывает контракт для получения всех курсов валют.
type AllRatesProvider interface {
	GetAllRates(ctx context.Context) (map[string]model.Rate, error)
}
