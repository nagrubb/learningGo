package main

type Quote struct {
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume int64
}

type StockQuoteService interface {
	GetLatestQuote(symbol string) (*Quote, error)
}
