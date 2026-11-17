// Package service содержит бизнес-логику приложения.
package service

//nolint:depguard
import (
	"context"
	"fmt"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/singleflight"
)

// ConverterService описывает сервис конвертации.
type ConverterService struct {
	cache   RateProvider
	primary AllRatesProvider
	backups []RateProvider
	sf      singleflight.Group
}

// NewConverterService конструктор - создаёт новый сервис.
func NewConverterService(cache RateProvider, primary AllRatesProvider, backups ...RateProvider) *ConverterService {
	return &ConverterService{
		cache:   cache,
		primary: primary,
		backups: backups,
	}
}

// Convert конвертирует сумму из одной валюты в другую.
func (s *ConverterService) Convert(
	ctx context.Context, amount decimal.Decimal, from, to string,
) (decimal.Decimal, error) {

	if from == to {
		return amount, nil
	}

	// В Рубли
	if to == model.RUB {
		rate, err := s.getRateToRub(ctx, from)
		if err != nil {
			return decimal.Zero, err
		}
		return amount.Mul(rate.Value), nil
	}

	// В другую валюту из Рублей
	if from == model.RUB {
		rate, err := s.getRateToRub(ctx, to)
		if err != nil {
			return decimal.Zero, err
		}
		return amount.Div(rate.Value), nil
	}

	// Межвалютные конвертации без Рублей
	fromRate, err := s.getRateToRub(ctx, from)
	if err != nil {
		return decimal.Zero, err
	}

	toRate, err := s.getRateToRub(ctx, to)
	if err != nil {
		return decimal.Zero, err
	}

	return amount.Mul(fromRate.Value).Div(toRate.Value), nil

}

func (s *ConverterService) getRateToRub(ctx context.Context, currency string) (model.Rate, error) {
	rate, err := s.cache.GetRate(ctx, currency, model.RUB)
	if err == nil {
		return rate, nil
	}

	val, err, _ := s.sf.Do("fetch_all_cbr", func() (interface{}, error) {
		allRates, err := s.primary.GetAllRates(ctx)
		if err != nil {
			return nil, err
		}

		for _, r := range allRates {
			if saver, ok := s.cache.(RateSaver); ok {
				_ = saver.SaveRate(ctx, r)
			}

			for _, backup := range s.backups {
				if saver, ok := backup.(RateSaver); ok {
					_ = saver.SaveRate(ctx, r)
				}
			}
		}

		return allRates, nil
	})

	if err == nil {
		allRates := val.(map[string]model.Rate)
		if r, ok := allRates[currency]; ok {
			return r, nil
		}
	}

	for _, backup := range s.backups {
		rate, err := backup.GetRate(ctx, currency, model.RUB)
		if err == nil {
			return rate, nil
		}
	}

	return model.Rate{}, fmt.Errorf("курс %s недоступен во всех источниках", currency)
}
