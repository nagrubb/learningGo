package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const hostName = "www.alphavantage.co"

type metaData struct {
	Symbol string `json:"2. Symbol"`
}

type intradayQuote struct {
	Open   json.Number `json:"1. open"`
	High   json.Number `json:"2. high"`
	Low    json.Number `json:"3. low"`
	Close  json.Number `json:"4. close"`
	Volume json.Number `json:"5. volume"`
}

type stockQuoteResponseByMinute struct {
	MetaData   metaData                 `json:"Meta Data"`
	TimeSeries map[string]intradayQuote `json:"Time Series (1min)"`
}

type AlphaAvantage struct {
	ApiKey string
}

func (a AlphaAvantage) GetLatestQuote(symbol string) (*Quote, error) {
	requestUrl := url.URL{Scheme: "https", Host: hostName, Path: "/query"}
	queryParams := requestUrl.Query()
	queryParams.Set("function", "TIME_SERIES_INTRADAY")
	queryParams.Add("interval", "1min")
	queryParams.Add("symbol", symbol)
	queryParams.Add("apikey", a.ApiKey)
	requestUrl.RawQuery = queryParams.Encode()

	rsp, err := http.Get(requestUrl.String())

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	var quotes stockQuoteResponseByMinute
	if err := json.Unmarshal(body, &quotes); err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("America/New_York")
	var latestTime *time.Time
	var latestQuote intradayQuote

	for k, v := range quotes.TimeSeries {
		layout := "2006-01-02 15:04:00"
		quoteTime, err := time.ParseInLocation(layout, k, loc)

		if err != nil {
			return nil, err
		}

		if latestTime == nil {
			latestTime = &quoteTime
			latestQuote = v
		} else if latestTime.Before(quoteTime) {
			latestTime = &quoteTime
			latestQuote = v
		}
	}

	var openPrice, closePrice, highPrice, lowPrice float64
	var volume int64

	openPrice, err = latestQuote.Open.Float64()

	if err == nil {
		closePrice, err = latestQuote.Close.Float64()
	} else {
		closePrice, _ = latestQuote.Close.Float64()
	}

	if err == nil {
		highPrice, err = latestQuote.High.Float64()
	} else {
		highPrice, _ = latestQuote.High.Float64()
	}

	if err == nil {
		lowPrice, err = latestQuote.Low.Float64()
	} else {
		lowPrice, _ = latestQuote.Low.Float64()
	}

	if err == nil {
		volume, err = latestQuote.Volume.Int64()
	} else {
		volume, _ = latestQuote.Volume.Int64()
	}

	return &Quote{Open: openPrice, Close: closePrice, High: highPrice, Low: lowPrice, Volume: volume}, err
}
