package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// conversionLoop запускает цикл конвертации.
func conversionLoop(scanner *bufio.Scanner) {
	for {
		fmt.Print("\nВведите запрос (например: 100 USD RUB) или 'выход'/'exit':\n")
		if !scanner.Scan() {
			return
		}

		input := strings.TrimSpace(scanner.Text())
		if isExitCommand(input) {
			return
		}
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		if len(parts) != 3 {
			fmt.Println("Ошибка: введите ровно 3 параметра.")
			continue
		}

		amount, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			fmt.Println("Ошибка: сумма должна быть числом.")
			continue
		}

		fromCur := normalizedCodes(parts[1])
		toCur := normalizedCodes(parts[2])
		doConvert(amount, fromCur, toCur)
	}
}

// historyLoop запускает цикл просмотра истории.
func historyLoop(scanner *bufio.Scanner) {
	for {
		fmt.Print("\nВведите валюту (например: USD) или 'выход'/'exit':\n")
		if !scanner.Scan() {
			return
		}

		input := strings.TrimSpace(scanner.Text())
		if isExitCommand(input) {
			return
		}
		if input == "" {
			continue
		}

		fromCur := normalizedCodes(input)

		fmt.Print("Введите количество дней (например: 7): ")
		if !scanner.Scan() {
			return
		}
		daysStr := strings.TrimSpace(scanner.Text())
		if isExitCommand(daysStr) {
			return
		}

		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 {
			fmt.Println("Ошибка: Дни должны быть положительным числом.")
			continue
		}
		showHistory(fromCur, days)
	}
}
