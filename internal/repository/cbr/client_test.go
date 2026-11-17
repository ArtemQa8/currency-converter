package cbr_test

//nolint:depguard
import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/repository/cbr"
	"github.com/stretchr/testify/assert"
)

func setupTestClient(t *testing.T, status int, body string) (*cbr.Client, func()) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))

	client := cbr.NewClient(server.URL)

	return client, server.Close
}

func TestClient_GetRate(t *testing.T) {
	t.Run("Успешный курс USD и JPY", func(t *testing.T) {
		body := `{
			"Date": "2026-02-22T11:30:00+03:00",
			"Valute": {
				"USD": {"CharCode": "USD", "Nominal": 1, "Value": 85.7},
				"JPY": {"CharCode": "JPY", "Nominal": 100, "Value": 55}
			}
		}`
		client, cleanup := setupTestClient(t, http.StatusOK, body)
		defer cleanup()

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
		client, cleanup := setupTestClient(t, http.StatusInternalServerError, "")
		defer cleanup()

		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Ошибка парсинга JSON", func(t *testing.T) {
		client, cleanup := setupTestClient(t, http.StatusOK, `это не джисончик`)
		defer cleanup()

		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Валюта не найдена", func(t *testing.T) {
		client, cleanup := setupTestClient(t, http.StatusOK, `{"Date": "2025-11-12T11:30:00+03:00", "Valute": {}}`)
		defer cleanup()

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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"a":`))
		}))
		defer server.Close()

		client := cbr.NewClient(server.URL)
		_, err := client.GetRate(context.Background(), "USD", "RUB")
		assert.Error(t, err)
	})

	t.Run("Успешное получение всех курсов", func(t *testing.T) {
		body := `{
			"Date": "2026-02-22T11:30:00+03:00",
			"Valute": {
				"USD": {"CharCode": "USD", "Nominal": 1, "Value": 85.7},
				"JPY": {"CharCode": "JPY", "Nominal": 100, "Value": 55}
			}
		}`
		client, cleanup := setupTestClient(t, http.StatusOK, body)
		defer cleanup()

		rates, err := client.GetAllRates(context.Background())
		assert.NoError(t, err)
		assert.Len(t, rates, 2)
	})
}
