package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	//model "tradingBot/src/exchanges/binance/model"
	models "TradingBot/src/models"

	"github.com/adshao/go-binance"
	"github.com/sirupsen/logrus"
)

const (
	priceURL          = "https://api.binance.com/api/v3/ticker/price"
	depthURL          = "https://api.binance.com/api/v1/depth"
	zero              = "0.00000000"
	orderBookMaxLimit = 1000
	candlestickLimit  = 1000
	apiInterval       = 1 * time.Second
)

type BinanceWorker struct {
	Cli              *binance.Client
	log              *logrus.Logger
	symbols          []string
	requestInterval  time.Duration
	AggTradesC       chan *binance.WsAggTradeEvent
	TradesC          chan *binance.WsTradeEvent
	KlinesC          chan *binance.WsKlineEvent
	stops            []chan struct{}
	dones            []chan struct{}
	orderBookCacheMu sync.Mutex
	orderBookCache   map[string]OrderBookInternal
}

// Create Binance worker: ApiKey, ApiSecret and other params from conf
func CreateWorker(conf *models.ExchangeConfig) *BinanceWorker {

	requestInterval, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	return &BinanceWorker{
		Cli:             binance.NewClient(conf.ApiKey, conf.ApiSecret),
		requestInterval: requestInterval,
		orderBookCache:  make(map[string]OrderBookInternal),
	}
}

// Start a new Binance worker.
func (b *BinanceWorker) Start() {

	symbol := "BTCUSDT"
	//for _, symbol := range b.symbols {
	go func(symbol string) {
		err := b.SubscribeOrderBook(symbol)
		if err != nil {
			b.log.Printf("Couldn't get diff depths on symbol %s: %v", symbol, err)
		}
	}(symbol)
	go b.SubscribeCandlestickAll(symbol)
	//}
}

func (b *BinanceWorker) SubscribeCandlestick(symbol, interval string) error {
	for ; ; <-time.Tick(b.requestInterval) {
		wsCandlestickHandler := func(event *binance.WsKlineEvent) {
			fmt.Println("Kline:", event.Kline.Open)
		}

		// Open a stream to wss://stream.binance.com:9443/ws/bnbbtc@depth
		doneC, _, err := binance.WsKlineServe(symbol, interval, wsCandlestickHandler, b.makeErrorHandler())
		if err != nil {
			return err
		}

		<-doneC
	}
}

func (b *BinanceWorker) SubscribeCandlestickAll(symbol string) {
	for _, v := range BinanceCandlestickIntervalList {
		go func(v string) {
			//b.initCandlesticks(symbol, s)
			if err := b.SubscribeCandlestick(symbol, v); err != nil {
				b.log.Errorf("Could not subscribe to candlestick interval %v symbol %v: %v", v, symbol, err)
			}
		}(v)
	}
}

func (b *BinanceWorker) initCandlesticks(symbol, interval string) {
	client := binance.NewClient("", "")
	candlesticks, err := client.NewKlinesService().Symbol(symbol).
		Interval(interval).Limit(candlestickLimit).Do(context.Background())
	if err != nil {
		b.log.Errorf("Could not load candlesticks from REST API with interval %v and symbol %v: %v",
			interval, symbol, err)

		return
	}

	for _, k := range candlesticks {
		fmt.Println(k)
		// if err := b.updateCandlestickAPI(symbol, interval, k); err != nil {
		// 	b.log.Errorf("Could not update candlesticks from REST API: %v", err)
		// }
	}
}

func (b *BinanceWorker) SubscribeOrderBook(symbol string) error {
	for ; ; <-time.Tick(b.requestInterval) {
		// Get a depth snapshot from https://www.binance.com/api/v1/depth?symbol=BNBBTC&limit=1000
		orderBook, err := b.getOrderBook(symbol, orderBookMaxLimit)

		// b.log.Debugf("Got order book for symbol %v: %+v", symbol, orderBook)

		if err != nil {
			return errors.Wrapf(err, "could not get order book")
		}

		b.orderBookCacheMu.Lock()
		b.orderBookCache[symbol] = orderBook
		b.orderBookCacheMu.Unlock()

		// Buffer the events you receive from the stream
		wsDepthHandler := func(event *binance.WsDepthEvent) {
			fmt.Println("Asks:", event.Asks)
			// if err = b.updateOrderBook(symbol, event); err != nil {
			// 	b.log.Errorf("Could not update order book: %v", err)
			// }
		}

		doneC, _, err := binance.WsDepthServe(symbol, wsDepthHandler, b.makeErrorHandler())
		if err != nil {
			return err
		}

		<-doneC
	}
}

func (b *BinanceWorker) getOrderBook(symbol string, depth int) (response OrderBookInternal, err error) {
	orderBookURL, err := b.makeOrderBookURL(symbol, depth)
	if err != nil {
		return OrderBookInternal{}, errors.Wrapf(err, "could not make order book URL")
	}

	resp, err := http.Get(orderBookURL)
	if err != nil {
		return OrderBookInternal{}, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(apiInterval)
	} else if resp.StatusCode != http.StatusOK {
		return OrderBookInternal{}, fmt.Errorf("getOrderBook received bad status code: %v", resp.StatusCode)
	}

	var data OrderBookResponse

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return OrderBookInternal{}, err
	}

	return SerializeBinanceOrderBookREST(data), nil
}

func (b *BinanceWorker) updateOrderBook(symbol string, event *binance.WsDepthEvent) error {
	b.orderBookCacheMu.Lock()
	defer b.orderBookCacheMu.Unlock()

	// Drop any event where u is <= lastUpdateId in the snapshot
	if event.UpdateID <= b.orderBookCache[symbol].LastUpdateID {
		return nil
	}

	for _, bid := range event.Bids {
		if bid.Quantity == zero {
			// b.log.Debugf("deleting bid with price %v for symbol %v", bid.Price, symbol)
			delete(b.orderBookCache[symbol].Bids, bid.Price)
			continue
		}

		b.orderBookCache[symbol].Bids[bid.Price] = bid.Quantity
	}

	for _, ask := range event.Asks {
		if ask.Quantity == zero {
			// b.log.Debugf("deleting ask with price %v for symbol %v", ask.Price, symbol)
			delete(b.orderBookCache[symbol].Asks, ask.Price)
			continue
		}

		b.orderBookCache[symbol].Asks[ask.Price] = ask.Quantity
	}

	// if err := b.database.StoreOrderBookInternal(symbol, w.orderBookCache[symbol]); err != nil {
	// 	b.log.Errorf("Could not store order book to database: %v", err)
	// }

	return nil
}

func (b *BinanceWorker) Klines(symbol, interval string) error {
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		b.KlinesC <- event
	}
	doneC, stopC, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, b.makeErrorHandler())
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *BinanceWorker) GetHistoryTrades(symbol string, start, end int64, number int) {
	trades, err := b.Cli.NewAggTradesService().
		Symbol(symbol).StartTime(start).EndTime(end).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	var sum int = 0
	var qul float64 = 0.
	var allQul float64 = 0.
	for _, t := range trades {
		i, _ := strconv.ParseFloat(t.Quantity, 64)
		allQul += i
		if t.IsBuyerMaker {
			sum++
			qul += i
		}

	}

	fmt.Println(float64(float64(sum)/float64(len(trades)))*100, float64(float64(qul)/float64(allQul))*100)
	// err = json.Unmarshal(jsonBlob, &rankings)
	// if err != nil {
	// 	// nozzle.printError("opening config file", err.Error())
	// }
	orders, err := json.Marshal(trades)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(orders))
	// rankingsJson, _ := json.Marshal(rankings)
	err = ioutil.WriteFile(fmt.Sprintf("%d.json", number), orders, 0644)
	// fmt.Printf("%+v", rankings)
}

func (b *BinanceWorker) makeErrorHandler() binance.ErrHandler {
	return func(err error) {
		b.log.Printf("Error in WS connection with Binance: %v", err)
	}
}

func (b *BinanceWorker) makeOrderBookURL(symbol string, depth int) (string, error) {
	u, err := url.Parse(depthURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("limit", strconv.Itoa(depth))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// func (b *BinanceWorker) updateCandlestick(symbol, interval string, event *binance.WsKlineEvent) error {
// 	if err := b.database.StoreCandlestickBinance(symbol, interval, event); err != nil {
// 		b.log.Errorf("Could not store candlestick to database: %v", err)
// 	}

// 	return nil
// }

// func (b *BinanceWorker) updateCandlestickAPI(symbol, interval string, candlestick *binance.Kline) error {
// 	if err := b.database.StoreCandlestickBinanceAPI(symbol, interval, candlestick); err != nil {
// 		b.log.Errorf("Could not store candlestick from REST API to database: %v", err)
// 	}

// 	return nil
// }
