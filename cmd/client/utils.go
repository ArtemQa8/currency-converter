package main

import (
	"fmt"
	"strings"
)

// isExitCommand проверяет, является ли ввод пользователя командой выхода в меню.
func isExitCommand(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input))
	exits := []string{
		"выход", "назад", "exit", "quit", "back",
		"ds[jl", "yfpfl", "учше", "йгше", "ифсл",
	}
	for _, e := range exits {
		if input == e {
			return true
		}
	}
	return false
}

// normalizedCodes приводит введенный текст к стандартному коду валюты с помощью словаря синонимов.
func normalizedCodes(input string) string {
	input = strings.ToUpper(input)
	if code, ok := currencyAliases[input]; ok {
		return code
	}
	return input
}

// formatAmount форматирует число с учетом точности валюты.
func formatAmount(amount float64, currency string) string {
	// Словарь точности для валют без копеек или с 3+ знаками.
	precisions := map[string]int{
		"JPY": 0, // Японская Йена.
		"KRW": 0, // Корейская Вона.
		"KWD": 3, // Кувейтский Динар.
		"BHD": 3, // Бахрейнский Динар.
		"BTC": 8, // Биткоин.
	}

	precision := 2 // По умолчанию.
	if p, ok := precisions[currency]; ok {
		precision = p
	}

	return fmt.Sprintf("%.*f", precision, amount)
}
