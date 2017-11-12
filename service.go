package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var provider StockQuoteService

type ServiceConfig struct {
	ApiKey string
}

type RestService struct {
	stockProvider StockQuoteService
}

func main() {
	var config ServiceConfig
	_, err := toml.DecodeFile("config.toml", &config)

	if err != nil {
		panic(err)
	}

	provider = &AlphaAvantage{ApiKey: config.ApiKey}
	service := &RestService{stockProvider: provider}
	service.startService(provider)
}

func (rs RestService) startService(provider StockQuoteService) {
	rs.stockProvider = provider
	router := mux.NewRouter()
	router.HandleFunc("/stock/quote/{symbol}", GetQuote).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetQuote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	quote, err := provider.GetLatestQuote(params["symbol"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(quote)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
