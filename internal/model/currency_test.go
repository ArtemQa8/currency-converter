package model_test

import (
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetPrecision(t *testing.T) {
	t.Run("Стандартная точность для доллара", func(t *testing.T) {
		assert.Equal(t, int32(2), model.GetPrecision("USD"))
	})

	t.Run("Нулевая точность для йены", func(t *testing.T) {
		assert.Equal(t, int32(0), model.GetPrecision("JPY"))
	})

	t.Run("Неизвестная валюта возвращает 2", func(t *testing.T) {
		assert.Equal(t, int32(2), model.GetPrecision("UNKNOWN"))
	})
}
