package poloniex

import (
	"TradingBot/src/models"
	"log"
	"time"

	polo "github.com/iowar/poloniex"
)

// Create Poloniex Worker
func CreateWorker(conf *models.ExchangeConfig) *PoloniexWorker {
	interval, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	// Public api client
	cli, err := polo.NewClient(conf.ApiKey, conf.ApiSecret)
	if err != nil {
		log.Fatal(err)
	}
	// Websocket client
	ws, err := polo.NewWSClient()
	if err != nil {
		log.Fatal(err)
	}

	return &PoloniexWorker{
		requestInterval: interval,
		CandleStick:     []polo.CandleStick{},
		Tickers:         map[string]*polo.Ticker{},
		PubCli:          cli,
		WsCli:           ws,
	}
}

// Start Worker
func (p *PoloniexWorker) Start() {

	//for _, symbol := range p.symbols {
	// go func(symbol string) {
	// 	err := w.SubscribeOrderBook(symbol)
	// 	if err != nil {
	// 		w.log.Printf("Couldn't get diff depths on symbol %s: %v", symbol, err)
	// 	}
	// }(symbol)
	symbol := "USDC_BTC"
	p.SubscribeMarkets(symbol)
	p.SubscribeTikers()

	//}
}

// func (p *PoloniexWorker) SubscribeCandlestickAll(symbol string) {
// 	for _, v := range PoloniexCandlestickIntervalList {
// 		go func(v int) {
// 			//p.initCandlesticks(symbol, s)

// 			if err := p.SubscribeCandlestick(symbol, v); err != nil {
// 				p.log.Errorf("Could not subscribe to candlestick interval %v symbol %v: %v", v, symbol, err)
// 			}
// 		}(v)
// 	}
// }

// func (p *PoloniexWorker) SubscribeCandlestick(symbol string, interval int) error {
// 	for ; ; <-time.Tick(p.requestInterval) {
// 		candles, err := p.poloniex.ChartData(symbol, interval, time.Now().Add(-3*p.requestInterval), time.Now().Add(3*p.requestInterval))

// 		if err != nil {
// 			p.log.Errorf("Could not get latest tick on poloniex: %v", err)
// 		}

// 		for _, candle := range candles {
// 			if err := p.updateCandlestickAPI(symbol, interval, candle); err != nil {
// 				p.log.Errorf("Could not update candlesticks from REST API: %v", err)
// 			}
// 		}
// 	}
// }

// func (p *PoloniexWorker) updateCandlestickAPI(symbol string, interval int, candlestick *poloniex.CandleStick) error {
// 	if err := p.database.StoreCandlestickPoloniexAPI(symbol, PoloniexIntervalToBinance(interval), candlestick); err != nil {
// 		p.log.Errorf("Could not store candlestick from REST API to database: %v", err)
// 	}

// 	return nil
// }
