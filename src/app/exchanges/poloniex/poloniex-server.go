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
	// Websocket clients
	ticker, err := NewWSClient()
	if err != nil {
		log.Fatal(err)
	}

	markt, err := NewWSClient()
	if err != nil {
		log.Fatal(err)
	}

	return &PoloniexWorker{
		requestInterval: interval,
		CandleStick:     []polo.CandleStick{},
		Tickers:         map[string]*polo.Ticker{},
		PubCli:          cli,
		WsTickers:       ticker,
		WsMarkets:       markt,
	}
}

// Start Worker
func (p *PoloniexWorker) Start() {
	// symbol := "USDC_BTC"
	// go p.SubscribeMarkets(symbol)
	go p.SubscribeTikers()
}

func (p *PoloniexWorker) SubscribeTikers() {
	err := p.WsTickers.SubscribeTicker()
	if err != nil {
		return
	}
	// for {
	// 	t := <-p.WsTickers.Subs["TICKER"]
	// 	fmt.Println("SymBol: ", t.(polo.WSTicker).Symbol)
	// }
}

func (p *PoloniexWorker) SubscribeMarkets(symbol string) {
	err := p.WsMarkets.SubscribeMarket(symbol)
	if err != nil {
		return
	}
	// for {
	// 	t := <-p.WsMarkets.Subs[symbol]
	// 	fmt.Println("Market: ", t.([]polo.MarketUpdate)[0].Data, "|          |", t.([]polo.MarketUpdate))
	// }
}
