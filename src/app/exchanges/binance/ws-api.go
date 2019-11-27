package binance

import (
	"TradingBot/src/app/tickers"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance"
	"github.com/pkg/errors"
)

// Get Tickers from chan(for Exchange interface)
func (b *BinanceWorker) GetTickersFromChan() *tickers.Tickers {
	WsTickers := <-b.AllMarketTickersC
	tiks := tickers.NewTickers()
	for _, ticker := range WsTickers {
		si, ok := b.SymbolInfo[ticker.Symbol]
		if !ok {
			log.Printf("warn, Binance - symbol %s missed in exchange info, skipped", ticker.Symbol)
			continue
		}
		base, quote := si.BaseAsset, si.QuoteAsset
		fmt.Println(base, quote)
		buyPrice, err := strconv.ParseFloat(ticker.BidPrice, 64)
		if err != nil {
			fmt.Printf("Error: %s in GetTickersFromChan in Binance Service", err.Error())
			return nil
		}
		sellPrice, err := strconv.ParseFloat(ticker.AskPrice, 64)
		if err != nil {
			fmt.Printf("Error: %s in GetTickersFromChan in Binance Service", err.Error())
			return nil
		}

		if buyPrice < 1e-6 || sellPrice < 1e-6 {
			continue
		}

		//tiks.Add(base, quote, tickers.Ticker{BuyPrice: buyPrice, SellPrice: sellPrice})
	}

	return tiks
}

// -----------------------------------   SUBSCRIBE TICKERS   -----------------------------------
//
func (b *BinanceWorker) AllMarketTickers() error {
	wsAllMarketTickersHandler := func(event binance.WsAllMarketsStatEvent) {
		b.AllMarketTickersC <- event
	}
	doneC, stopC, err := binance.WsAllMarketsStatServe(wsAllMarketTickersHandler, b.makeErrorHandler())
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

// -----------------------------------   SUBSCRIBE TRADES   -----------------------------------
//
func (b *BinanceWorker) AggTrades(symbol string) error {
	wsAggTradesHandler := func(event *binance.WsAggTradeEvent) {
		b.AggTradesC <- event
	}

	doneC, stopC, err := binance.WsAggTradeServe(symbol, wsAggTradesHandler, b.makeErrorHandler())
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

// -----------------------------------   SUBSCRIBE KLINES   -----------------------------------
//
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

// -----------------------------------   SUBSCRIBE CANDLESTICK   -----------------------------------
//
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

// Subsctibe candlestick use time interval
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

// -----------------------------------   SUBSCRIBE ORDER BOOK   -----------------------------------
//
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
