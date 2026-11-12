package cbr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/shopspring/decimal"
)

type Client struct {
	httpClient *http.Client
	apiURL     string
}

func NewClient(apiURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiURL:     apiURL,
	}
}

type cbrResponse struct {
	Date   time.Time `json:"Date"`
	Valute map[string]struct {
		CharCode string  `json:"CharCode"`
		Nominal  int     `json:"Nominal"`
		Value    float64 `json:"Value"`
	} `json:"Valute"`
}

// GetRate получает курс указанной валюты к рублю через API ЦБ
// Поддерживает только рубль в качестве второй валюты (quoteCurrency)
func (c *Client) GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error) {

	if quoteCurrency != "RUB" {
		return model.Rate{}, fmt.Errorf("клиент ЦБ поддерживает только RUB в качестве котируемой валюты, получено: %s", quoteCurrency)
	}

	// Курс RUB к RUB всегда =1, не делаем лишний запрос в API
	if baseCurrency == "RUB" {
		return model.Rate{Value: decimal.NewFromInt(1), Source: "cbr"}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL, nil)
	if err != nil {
		return model.Rate{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.Rate{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Rate{}, fmt.Errorf("API ЦБ вернул статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Rate{}, err
	}

	var data cbrResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return model.Rate{}, err
	}

	val, ok := data.Valute[baseCurrency]
	if !ok {
		return model.Rate{}, fmt.Errorf("валюта %s не найдена в ответе ЦБ", baseCurrency)
	}

	// ЦБ отдает некоторые валюты за 100/1000 единиц
	// Поэтому приводим курс к стоимости за 1 единицу
	rateDec := decimal.NewFromFloat(val.Value)
	nominalDec := decimal.NewFromInt(int64(val.Nominal))
	normalizedValue := rateDec.Div(nominalDec)

	return model.Rate{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: "RUB",
		Value:         normalizedValue,
		Date:          data.Date,
		Source:        "cbr",
	}, nil
}
