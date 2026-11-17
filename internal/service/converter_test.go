package service_test

//nolint:depguard
import (
	"context"
	"errors"
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	rates map[string]decimal.Decimal
	err   error
}

func (m *mockProvider) GetRate(_ context.Context, base, _ string) (model.Rate, error) {
	if m.err != nil {
		return model.Rate{}, m.err
	}
	if val, ok := m.rates[base]; ok {
		return model.Rate{BaseCurrency: base, QuoteCurrency: model.RUB, Value: val, Source: "cbr"}, nil
	}
	return model.Rate{}, errors.New("валюта не найдена")
}

func (m *mockProvider) SaveRate(_ context.Context, _ model.Rate) error { return nil }

func (m *mockProvider) GetAllRates(_ context.Context) (map[string]model.Rate, error) {
	if m.err != nil {
		return nil, m.err
	}
	res := make(map[string]model.Rate)
	for k, v := range m.rates {
		res[k] = model.Rate{BaseCurrency: k, QuoteCurrency: model.RUB, Value: v, Source: "cbr"}
	}
	return res, nil
}

func (m *mockProvider) GetHistory(_ context.Context, _ string, _ string, _ int) ([]model.Rate, error) {
	return nil, nil
}

type backupMock struct {
	rate model.Rate
}

func (b *backupMock) GetRate(_ context.Context, _, _ string) (model.Rate, error) {
	return b.rate, nil
}

func (b *backupMock) SaveRate(_ context.Context, _ model.Rate) error { return nil }

func (b *backupMock) GetAllRates(_ context.Context) (map[string]model.Rate, error) {
	return nil, nil
}

func (b *backupMock) GetHistory(_ context.Context, _, _ string, _ int) ([]model.Rate, error) {
	return nil, nil
}

func TestConverterService_Convert(t *testing.T) {
	mock := &mockProvider{
		rates: map[string]decimal.Decimal{
			model.USD: decimal.NewFromInt(90),
			model.EUR: decimal.NewFromInt(100),
		},
	}

	svc := service.NewConverterService(mock, mock, mock)

	t.Run("Одинаковые валюты", func(t *testing.T) {
		res, err := svc.Convert(context.Background(), decimal.NewFromInt(100), model.USD, model.USD)
		assert.NoError(t, err)
		assert.Equal(t, "100", res.String())
	})

	t.Run("Перевод из USD в RUB", func(t *testing.T) {
		res, err := svc.Convert(context.Background(), decimal.NewFromInt(100), model.USD, model.RUB)
		assert.NoError(t, err)
		assert.Equal(t, "9000", res.String())
	})

	t.Run("Перевод из RUB в USD", func(t *testing.T) {
		res, err := svc.Convert(context.Background(), decimal.NewFromInt(9000), model.RUB, model.USD)
		assert.NoError(t, err)
		assert.Equal(t, "100", res.String())
	})

	t.Run("Кросс-курс EUR в USD", func(t *testing.T) {
		res, err := svc.Convert(context.Background(), decimal.NewFromInt(100), model.EUR, model.USD)
		assert.NoError(t, err)
		assert.Equal(t, "111.11", res.Round(2).String())
	})

	t.Run("Ошибка при конвертации в RUB", func(t *testing.T) {
		errorMock := &mockProvider{err: errors.New("API недоступен")}
		errorSvc := service.NewConverterService(errorMock, errorMock, errorMock)

		_, err := errorSvc.Convert(context.Background(), decimal.NewFromInt(100), model.USD, model.RUB)
		assert.Error(t, err)
	})

	t.Run("Ошибка при конвертации из RUB", func(t *testing.T) {
		errorMock := &mockProvider{err: errors.New("API недоступен")}
		errorSvc := service.NewConverterService(errorMock, errorMock, errorMock)

		_, err := errorSvc.Convert(context.Background(), decimal.NewFromInt(100), model.RUB, model.USD)
		assert.Error(t, err)
	})

	t.Run("Ошибка при кросс-курсе (первая валюта)", func(t *testing.T) {
		errorMock := &mockProvider{err: errors.New("API недоступен")}
		errorSvc := service.NewConverterService(errorMock, errorMock, errorMock)

		_, err := errorSvc.Convert(context.Background(), decimal.NewFromInt(100), model.EUR, model.USD)
		assert.Error(t, err)
	})

	t.Run("Ошибка при кросс-курсе (вторая валюта)", func(t *testing.T) {
		_, err := svc.Convert(context.Background(), decimal.NewFromInt(100), model.EUR, "GBP")
		assert.Error(t, err)
	})
}

func TestConverterService_BackupFallback(t *testing.T) {
	brokenMock := &mockProvider{err: errors.New("API недоступен")}

	backup := &backupMock{
		rate: model.Rate{BaseCurrency: model.USD, QuoteCurrency: model.RUB, Value: decimal.NewFromInt(95)},
	}

	svc := service.NewConverterService(brokenMock, brokenMock, backup)

	t.Run("Спасает бэкап, когда ЦБ и кэш упали", func(t *testing.T) {
		res, err := svc.Convert(context.Background(), decimal.NewFromInt(100), model.USD, model.RUB)
		assert.NoError(t, err)
		assert.Equal(t, "9500", res.String())
	})
}
