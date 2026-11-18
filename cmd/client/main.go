// Package main реализует CLI-клиент.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printHeader()
		fmt.Println("1. Конвертация")
		fmt.Println("2. История курсов")
		fmt.Println("0. Выход")
		fmt.Print("Ваш выбор: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())
		switch choice {
		case "1":
			conversionLoop(scanner)
		case "2":
			historyLoop(scanner)
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			if choice != "" {
				fmt.Println("Неверный выбор. Попробуйте снова.")
			}
		}
	}
}
