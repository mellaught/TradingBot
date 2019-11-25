package poloniex

import (
	"time"

	polo "github.com/iowar/poloniex"
)

//
func (p *PoloniexWorker) GetPrice(symbols ...string) ([]polo.Ticker, error) {

	resp, err := p.PubCli.GetTickers()
	if err != nil {
		return nil, err
	}
	var t []polo.Ticker
	for _, s := range symbols {
		t = append(t, resp[s])
	}

	return t, nil
}

// Returns the past 200 trades for a given market, or up to 1,000 trades between a range specified in UNIX timestamps by the "start" and "end" GET parameters.
// Публичная история сделок.
func (p *PoloniexWorker) GetTradeHistory(symbol string) ([]polo.PublicTrade, error) {
	resp, err := p.PubCli.GetPublicTradeHistory(symbol, time.Now().AddDate(0, 0, -1), time.Now())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns candlestick chart data. Required GET parameters are "currencyPair", "period" (candlestick period in seconds; valid values are 300, 900, 1800, 7200, 14400, and 86400),
// "start", and "end". "Start" and "end" are given in UNIX timestamp format and used to specify the date range for the data returned.
// Свечи
func (p *PoloniexWorker) GetCandles(symbol string) ([]polo.CandleStick, error) {
	resp, err := p.PubCli.GetChartData(symbol, time.Now().AddDate(0, 0, -1), time.Now(), "1d")
	if err != nil {
		return nil, err
	}

	return resp, nil
}
