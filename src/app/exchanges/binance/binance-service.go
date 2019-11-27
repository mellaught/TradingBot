package binance

import (
	"log"
	"time"

	//model "tradingBot/src/exchanges/binance/model"
	models "TradingBot/src/models"

	"github.com/adshao/go-binance"
)

const (
	comission         = 0.001
	priceURL          = "https://api.binance.com/api/v3/ticker/price"
	depthURL          = "https://api.binance.com/api/v1/depth"
	zero              = "0.00000000"
	orderBookMaxLimit = 1000
	candlestickLimit  = 1000
	apiInterval       = 1 * time.Second
)

// Returns Binance operation comission
func (b *BinanceWorker) GetBinanceComission() float64 {
	return comission
}

// Returns Binance operation comission
func (b *BinanceWorker) GetName() string {
	return "Binance"
}

// Create Binance worker: ApiKey, ApiSecret and other params from conf
func CreateWorker(conf *models.ExchangeConfig) *BinanceWorker {

	requestInterval, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	return &BinanceWorker{
		Cli:               binance.NewClient(conf.ApiKey, conf.ApiSecret),
		requestInterval:   requestInterval,
		AggTradesC:        make(chan *binance.WsAggTradeEvent),
		SymbolInfo:        make(map[string]*binance.Symbol),
		AllMarketTickersC: make(chan binance.WsAllMarketsStatEvent),
		orderBookCache:    make(map[string]OrderBookInternal),
	}

}

// Start a new Binance worker.
func (b *BinanceWorker) Start() {

	// Start subscribe tickers
	go b.AllMarketTickers()

}

// func (b *BinanceWorker) initCandlesticks(symbol, interval string) {
// 	client := binance.NewClient("", "")
// 	candlesticks, err := client.NewKlinesService().Symbol(symbol).
// 		Interval(interval).Limit(candlestickLimit).Do(context.Background())
// 	if err != nil {
// 		b.log.Errorf("Could not load candlesticks from REST API with interval %v and symbol %v: %v",
// 			interval, symbol, err)

// 		return
// 	}

// 	for _, k := range candlesticks {
// 		fmt.Println(k)
// 		// if err := b.updateCandlestickAPI(symbol, interval, k); err != nil {
// 		// 	b.log.Errorf("Could not update candlesticks from REST API: %v", err)
// 		// }
// 	}
// }
