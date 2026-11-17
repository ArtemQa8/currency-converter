package model

const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
)

// Currency описывает модель валюты.
type Currency struct {
	Code      string
	Name      string
	Symbol    string
	Precision int
}

var currencyPrecisions = map[string]int32{
	"JPY": 0, // Японская Йена
	"KRW": 0, // Корейская Вона
	"KWD": 3, // Кувейтский Динар
	"BHD": 3, // Бахрейнский Динар
	"BTC": 8, // Биткоин
}

// GetPrecision возвращает количество знаков после запятой.
func GetPrecision(currencyCode string) int32 {
	if p, ok := currencyPrecisions[currencyCode]; ok {
		return p
	}

	return 2
}
