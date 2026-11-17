package api

import "time"

// ConvertRequest описывает запрос на конвертацию.
type ConvertRequest struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

// ConverterResponse описывает ответ конвертации.
type ConverterResponse struct {
	Result float64 `json:"result"`
}

// HistoryRequest описывает запрос истории.
type HistoryRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
	Days int    `json:"days"`
}

// HistoryObject описывает один объект истории.
type HistoryObject struct {
	Date  time.Time `json:"date"`
	From  string    `json:"from"`
	To    string    `json:"to"`
	Value float64   `json:"value"`
}

// HistoryResponse описывает ответ истории.
type HistoryResponse struct {
	Rates []HistoryObject `json:"rates"`
}
