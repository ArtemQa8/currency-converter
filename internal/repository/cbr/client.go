package cbr

//nolint:depguard
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
	// ЦБ отдает с большой буквы.
	Date   time.Time `json:"Date"` //nolint:tagliatelle
	Valute map[string]struct {
		CharCode string  `json:"CharCode"` //nolint:tagliatelle
		Nominal  int     `json:"Nominal"`  //nolint:tagliatelle
		Value    float64 `json:"Value"`    //nolint:tagliatelle
	} `json:"Valute"` //nolint:tagliatelle
}

// fetchAll скачивает JSON и возвращает мапу со всеми курсами к рублю.
func (c *Client) fetchAll(ctx context.Context) (map[string]model.Rate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API ЦБ вернул статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data cbrResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	rates := make(map[string]model.Rate)

	for code, val := range data.Valute {
		rateDec := decimal.NewFromFloat(val.Value)
		nominalDec := decimal.NewFromInt(int64(val.Nominal))
		normalizedValue := rateDec.Div(nominalDec)

		rates[code] = model.Rate{
			BaseCurrency:  code,
			QuoteCurrency: model.RUB,
			Value:         normalizedValue,
			Date:          data.Date,
			Source:        "cbr",
		}
	}

	return rates, nil
}

// GetAllRates возвращает курсы всех валют из ЦБ
func (c *Client) GetAllRates(ctx context.Context) (map[string]model.Rate, error) {
	return c.fetchAll(ctx)
}

// GetRate получает курс указанной валюты к рублю через API ЦБ
// Поддерживает только рубль в качестве второй валюты (quoteCurrency)
func (c *Client) GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error) {

	if quoteCurrency != model.RUB {
		return model.Rate{},
			fmt.Errorf("клиент ЦБ поддерживает только RUB "+
				"в качестве котируемой валюты, получено: %s",
				quoteCurrency)
	}

	// Курс RUB к RUB всегда =1, не делаем лишний запрос в API
	if baseCurrency == model.RUB {
		return model.Rate{Value: decimal.NewFromInt(1), Source: "cbr"}, nil
	}

	allRates, err := c.fetchAll(ctx)
	if err != nil {
		return model.Rate{}, err
	}

	rate, ok := allRates[baseCurrency]
	if !ok {
		return model.Rate{}, fmt.Errorf("валюта %s не найдена в ответе ЦБ", baseCurrency)
	}

	return rate, nil
}
