package cbr_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/repository/cbr"
	"github.com/stretchr/testify/assert"
)

func newMockServer(status int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(responseBody))
	}))
}

func TestClient_GetRate(t *testing.T) {
	t.Run("Успешный курс USD и JPY", func(t *testing.T) {
		server := newMockServer(http.StatusOK, `{
			"Date": "2026-02-22T11:30:00+03:00",
			"Valute": {
				"USD": {"CharCode": "USD", "Nominal": 1, "Value": 85.7},
				"JPY": {"CharCode": "JPY", "Nominal": 100, "Value": 55}
			}
		}`)
		defer server.Close()

		client := cbr.NewClient(server.URL)

		rateUSD, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.NoError(t, err)
		assert.Equal(t, "85.7", rateUSD.Value.String())

		rateJPY, err := client.GetRate(context.Background(), "JPY", "RUB")
		assert.NoError(t, err)
		assert.Equal(t, "0.55", rateJPY.Value.String())
	})

	t.Run("Ошибка при неверной котируемой валюте", func(t *testing.T) {
		client := cbr.NewClient("http://localhost")
		_, err := client.GetRate(context.Background(), "USD", "EUR")
		assert.Error(t, err)
	})

	t.Run("Успешный курс RUB к RUB", func(t *testing.T) {
		client := cbr.NewClient("http://localhost")
		rate, err := client.GetRate(context.Background(), "RUB", "RUB")
		assert.NoError(t, err)
		assert.Equal(t, "1", rate.Value.String())
	})

	t.Run("Ошибка сервера 500", func(t *testing.T) {
		server := newMockServer(http.StatusInternalServerError, "")
		defer server.Close()
		client := cbr.NewClient(server.URL)

		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Ошибка парсинга JSON", func(t *testing.T) {
		server := newMockServer(http.StatusOK, `это не джисончик`)
		defer server.Close()
		client := cbr.NewClient(server.URL)

		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Валюта не найдена", func(t *testing.T) {
		server := newMockServer(http.StatusOK, `{"Date": "2025-11-12T11:30:00+03:00", "Valute": {}}`)
		defer server.Close()
		client := cbr.NewClient(server.URL)

		_, err := client.GetRate(context.Background(), "GBP", "RUB")
		assert.Error(t, err)
	})

	t.Run("Ошибка сети", func(t *testing.T) {
		client := cbr.NewClient("http://127.0.0.1:0")
		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Ошибка создания запроса (неверный URL)", func(t *testing.T) {
		client := cbr.NewClient("http://invalid url.com")
		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Ошибка чтения тела ответа (обрыв соединения)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"a":`))
		}))
		defer server.Close()

		client := cbr.NewClient(server.URL)
		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})
}
