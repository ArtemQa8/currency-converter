package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	usd = "USD"
	eur = "EUR"
	rub = "RUB"
	jpy = "JPY"
	krw = "KRW"
	cny = "CNY"
)

var currencyAliases = map[string]string{
	"ДОЛЛАР": usd, "ДОЛЛАРЫ": usd, "ДОЛ": usd,
	"ДОЛАР": usd, "ДОЛЛАРОВ": usd, "ДОЛАРОВ": usd,
	"ДОЛЛ": usd, "БАКС": usd, "БАКСОВ": usd,
	"$": usd, "BUCKS": usd, "DOLLAR": usd, "DOL": usd, "DOLAR": usd,
	"LJK": usd, "LJKK": usd, "ВЩД": usd, "ГЫВ": usd,

	"ЕВРО": eur, "ЕВР": eur,
	"€": eur, "EURO": eur,
	"УГКЩ": eur, "TDHJ": eur, "УГК": eur, "TDH": eur,

	"РУБЛЬ": rub, "РУБЛИ": rub, "РУБ": rub, "РУБЛЕЙ": rub, "РУБЛЯ": rub,
	"₽": rub, "RUBLES": rub,
	"КГИДУЫ": rub, "КГИ": rub, "КУИ": rub, "HE<": rub,

	"ЙЕНА": jpy, "ЙЕНЫ": jpy, "ЕНА": jpy, "ЙЕН": jpy, "ЕН": jpy,
	"¥": jpy, "JENA": jpy, "YENA": jpy, "YEN": jpy, "JEN": jpy,
	"ОЗН": jpy, "QTYF": jpy, "НУТФ": jpy, "ОУТФ": jpy,

	"ВОНА": krw, "ВОНЫ": krw, "ВОН": krw,
	"₩": krw, "WON": krw,
	"ЛКЦ": krw, "DJYF": krw, "DJY": krw, "ЦЩТ": krw,

	"ЮАНЬ": cny, "ЮАНЕЙ": cny, "ЮАНЯ": cny, "ЮАНИ": cny,
	"Ұ": cny, "YUAN": cny, "RENMINBI": cny, "RMB": cny,
	"ЖЕНЬМИНЬБИ": cny, "РЕНМИНБИ": cny, "ЖМБ": cny, "РМБ": cny,
	">FYM": cny, "НГФТ": cny,
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Конвертер Валют ===")
	fmt.Println("Введите запрос в формате: <сумма> <откуда> <куда>")
	fmt.Println("Пример: 100 USD RUB")
	fmt.Println("Для выхода введите 'exit', 'quit' или 'выход'")
	fmt.Println("-----------------------")

	for {
		fmt.Print("Введите запрос: ")

		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			continue
		}

		lowerText := strings.ToLower(text)

		if lowerText == "exit" || lowerText == "quit" ||
			lowerText == "q" || lowerText == "ds[jl" ||
			lowerText == "выход" || lowerText == "й" ||
			lowerText == "учше" || lowerText == "йгше" {
			fmt.Println("До свидания!")
			break
		}

		parts := strings.Fields(text)
		if len(parts) != 3 {
			fmt.Println("Ошибка: Введите ровно 3 параметра. Пример: 100 USD RUB")
			continue
		}

		amount, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			fmt.Println("Ошибка: Сумма должна быть числом")
			continue
		}

		fromCur := normalizeCodes(parts[1])
		toCur := normalizeCodes(parts[2])

		reqData := map[string]interface{}{
			"amount": amount,
			"from":   fromCur,
			"to":     toCur,
		}
		jsonReq, _ := json.Marshal(reqData)

		req, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"http://localhost:8080/api/convert",
			strings.NewReader(string(jsonReq)))
		if err != nil {
			fmt.Println("Ошибка создания запроса:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Ошибка соединения: сервер не запущен")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			fmt.Println("Ошибка чтения ответа:", err)
			continue
		}

		var result struct {
			Result float64 `json:"result"`
			Error  string  `json:"error"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Ошибка парсинга ответа сервера")
			continue
		}

		if result.Error != "" {
			fmt.Printf("Ошибка сервера: %s\n", result.Error)
		} else {
			fmt.Printf("%v %s = %v %s\n", amount, fromCur, result.Result, toCur)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения:", err)
	}
}

func normalizeCodes(input string) string {
	input = strings.ToUpper(input)

	if code, ok := currencyAliases[input]; ok {
		return code
	}

	return input
}
