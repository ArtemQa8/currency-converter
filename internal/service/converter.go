package service

import (
	"context"
	"fmt"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/shopspring/decimal"
)

type ConverterService struct {
	providers []RateProvider
}

func NewConverterService(providers ...RateProvider) *ConverterService {
	return &ConverterService{providers: providers}
}

func (s *ConverterService) Convert(ctx context.Context, amount decimal.Decimal, from, to string) (decimal.Decimal, error) {

	if from == to {
		return amount, nil
	}

	// В Рубли
	if to == "RUB" {
		rate, err := s.getRateToRub(ctx, from)
		if err != nil {
			return decimal.Zero, err
		}
		return amount.Mul(rate.Value), nil
	}

	// В другую валюту из Рублей
	if from == "RUB" {
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

	if len(s.providers) == 0 {
		return model.Rate{}, fmt.Errorf("нет доступных источников для получения курса")
	}

	var lastErr error

	for _, provider := range s.providers {
		rate, err := provider.GetRate(ctx, currency, "RUB")
		if err != nil {
			lastErr = err
			continue
		}
		return rate, nil
	}

	return model.Rate{}, lastErr
}
