package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// doConvert отправляет POST-запрос на сервер для конвертации суммы и выводит результат.
func doConvert(amount float64, from, to string) {
	reqData := map[string]interface{}{"amount": amount, "from": from, "to": to}
	jsonReq, _ := json.Marshal(reqData)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"http://localhost:8080/api/convert",
		strings.NewReader(string(jsonReq)),
	)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка соединения: сервер не запущен")
		return
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Result float64 `json:"result"`
		Error  string  `json:"error"`
	}
	_ = json.Unmarshal(body, &result)

	if result.Error != "" {
		fmt.Printf("\nОшибка сервера: %s\n", result.Error)
	} else {
		fmt.Printf("\n\033[1mРезультат:\033[0m \033[36m%v %s = %.2f %s\033[0m\n", amount, from, result.Result, to)
	}
}

// showHistory запрашивает у сервера историю курса валюты за N дней и печатает её.
func showHistory(from string, days int) {
	url := fmt.Sprintf("http://localhost:8080/api/history?from=%s&to=RUB&days=%d", from, days)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка соединения: сервер не запущен.")
		return
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Rates []struct {
			Date  string  `json:"date"`
			Value float64 `json:"value"`
		} `json:"rates"`
	}
	_ = json.Unmarshal(body, &result)

	if len(result.Rates) == 0 {
		fmt.Println("\nИстория пуста.")
		return
	}

	fmt.Printf("\nИстория %s -> RUB:\n", from)
	for _, r := range result.Rates {
		dateStr := r.Date
		if len(dateStr) > 10 {
			dateStr = dateStr[:10]
		}
		fmt.Printf(" %s: %.2f\n", dateStr, r.Value)
	}
}

// fetchRate получает последнюю историю валюты к рублю и считает динамику.
func fetchRate(currency string) (float64, float64, error) {
	url := fmt.Sprintf("http://localhost:8080/api/history?from=%s&to=RUB&days=2", currency)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var result struct {
		Rates []struct {
			Value float64 `json:"value"`
		} `json:"rates"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, err
	}

	if len(result.Rates) == 0 {
		return 0, 0, fmt.Errorf("нет данных")
	}

	last := result.Rates[len(result.Rates)-1].Value

	if len(result.Rates) < 2 {
		return last, 0, nil
	}

	prev := result.Rates[0].Value
	return last, prev, nil
}
