package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var currencyAliases = map[string]string{
	"ДОЛЛАР": "USD", "ДОЛЛАРЫ": "USD", "ДОЛ": "USD",
	"ДОЛАР": "USD", "ДОЛЛАРОВ": "USD", "ДОЛАРОВ": "USD",
	"ДОЛЛ": "USD", "DOLLAR": "USD", "DOL": "USD",
	"$": "USD", "БАКС": "USD", "BUCKS": "USD", "БАКСОВ": "USD",
	"LJK": "USD", "LJKK": "USD", "ВЩД": "USD", "ГЫВ": "USD",

	"ЕВРО": "EUR", "ЕВР": "EUR", "€": "EUR",
	"УГКЩ": "EUR", "TDHJ": "EUR", "УГК": "EUR", "TDH": "EUR",

	"РУБЛЬ": "RUB", "РУБЛИ": "RUB", "РУБ": "RUB", "₽": "RUB",
	"РУБЛЕЙ": "RUB", "РУБЛЯ": "RUB", "RUBLES": "RUB", "КГИДУЫ": "RUB",
	"КГИ": "RUB", "КУИ": "RUB",

	"ЙЕНА": "JPY", "ЙЕНЫ": "JPY", "¥": "JPY", "ОЗН": "JPY",

	"ВОНА": "KRW", "ВОНЫ": "KRW", "₩": "KRW", "ЛКЦ": "KRW",
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

		resp, err := http.Post("http://localhost:8080/api/convert", "application/json", strings.NewReader(string(jsonReq)))
		if err != nil {
			fmt.Println("Ошибка соединения: сервер не запущен")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
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
