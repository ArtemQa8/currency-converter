package api

import "time"

type ConvertRequest struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

type ConverterResponse struct {
	Result float64 `json:"result"`
}

type HistoryRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
	Days int    `json:"days"`
}

type HistoryObject struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type HistoryResponse struct {
	Rates []HistoryObject `json:"rates"`
}
