package api_test

//nolint:depguard
import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/api"
	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m *mockRepo) GetRate(_ context.Context, base, _ string) (model.Rate, error) {
	return model.Rate{BaseCurrency: base, QuoteCurrency: model.RUB, Value: decimal.NewFromInt(90)}, nil
}
func (m *mockRepo) GetAllRates(_ context.Context) (map[string]model.Rate, error) { return nil, nil }
func (m *mockRepo) SaveRate(_ context.Context, _ model.Rate) error               { return nil }
func (m *mockRepo) GetHistory(_ context.Context, from, to string, _ int) ([]model.Rate, error) {
	return []model.Rate{
		{BaseCurrency: from, QuoteCurrency: to, Value: decimal.NewFromFloat(90.5)},
	}, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	repo := &mockRepo{}
	svc := service.NewConverterService(repo, repo, repo)
	h := api.NewHandler(svc)
	return api.NewRouter(h)
}

func TestAPI_Convert(t *testing.T) {
	router := setupRouter()

	t.Run("Успешная конвертация (200)", func(t *testing.T) {
		body := `{"amount": 100, "from": "USD", "to": "RUB"}`
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/convert",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"result":9000`)
	})

	t.Run("Битый JSON (400)", func(t *testing.T) {
		body := `{"amount": "abc", "from": "USD"}`
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/convert",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAPI_History(t *testing.T) {
	router := setupRouter()

	t.Run("Успешная история (200)", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/api/history?from=USD&to=RUB&days=7",
			nil,
		)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"from":"USD"`)
		assert.Contains(t, w.Body.String(), `"value":90.5`)
	})

	t.Run("Не переданы обязательные параметры (400)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/history?days=7", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
