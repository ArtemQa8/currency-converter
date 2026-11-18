package main

import "fmt"

// printHeader рисует шапку с курсами и цветной динамикой.
func printHeader() {
	fmt.Println("=========================================")
	fmt.Println("           КОНВЕРТЕР ВАЛЮТ")
	fmt.Println("=========================================")
	fmt.Println("Актуальные курсы ЦБ РФ:")

	currencies := map[string]string{
		"USD": "Доллар",
		"EUR": "Евро",
	}

	for code, name := range currencies {
		last, prev, err := fetchRate(code)
		if err != nil {
			fmt.Printf("%s: Данные недоступны\n", name)
			continue
		}

		if prev == 0 {
			fmt.Printf("%s: %.2f руб.\n", name, last)
			continue
		}

		diff := last - prev
		percent := (diff / prev) * 100

		switch {
		case diff > 0:
			fmt.Printf("%s: %.2f руб. \033[32m▲ +%.2f%%\033[0m\n", name, last, percent)
		case diff < 0:
			fmt.Printf("%s: %.2f руб. \033[31m▼ %.2f%%\033[0m\n", name, last, percent)
		default:
			fmt.Printf("%s: %.2f руб.\n", name, last)
		}
	}
	fmt.Println("-----------------------------------------")
}
