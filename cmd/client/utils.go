package main

import "strings"

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
