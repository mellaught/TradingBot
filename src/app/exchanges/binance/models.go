package binance

import (
	"sync"
	"time"

	"github.com/adshao/go-binance"
	"github.com/sirupsen/logrus"
)

var (
	BinanceCandlestickIntervalList = []string{
		"1m",
		"3m",
		"5m",
		"15m",
		"30m",
		"1h",
		"2h",
		"4h",
		"6h",
		"8h",
		"12h",
		"1d",
		"3d",
		"1w",
		"1M",
	}

	BittrexCandlestickIntervalList = []string{
		"oneMin", "fiveMin", "thirtyMin", "hour", "day",
	}
)

type BinanceWorker struct {
	Cli               *binance.Client
	log               *logrus.Logger
	symbols           []string
	requestInterval   time.Duration
	AggTradesC        chan *binance.WsAggTradeEvent
	TradesC           chan *binance.WsTradeEvent
	AllMarketTickersC chan binance.WsAllMarketsStatEvent
	KlinesC           chan *binance.WsKlineEvent
	stops             []chan struct{}
	dones             []chan struct{}
	orderBookCacheMu  sync.Mutex
	orderBookCache    map[string]OrderBookInternal
}

type OrderBookInternal struct {
	LastUpdateID int64             `json:"-"`
	Bids         map[string]string `json:"bids"`
	Asks         map[string]string `json:"asks"`
}

type OrderBookResponse struct {
	LastUpdateID int64       `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"` // price, quantity
	Asks         [][2]string `json:"asks"` // price, quantity
}

func SerializeBinanceOrderBookREST(data OrderBookResponse) OrderBookInternal {
	asks := make(map[string]string)
	bids := make(map[string]string)

	for _, ask := range data.Asks {
		asks[ask[0]] = ask[1]
	}

	for _, bid := range data.Bids {
		bids[bid[0]] = bid[1]
	}

	return OrderBookInternal{
		LastUpdateID: data.LastUpdateID,
		Asks:         asks,
		Bids:         bids,
	}
}

func IsValidInterval(s string) bool {
	for _, v := range BinanceCandlestickIntervalList {
		if v == s {
			return true
		}
	}
	return false
}

type TradeHistories struct {
	Orders []TradeHistory `json:"orders"`
}

type TradeHistory struct {
	AggTradeID       int64  `json:"aggTradeID"`
	FirstTradeID     int64  `json:"firstTradeID"`
	IsBestPriceMatch bool   `json:"isBestPriceMatch"`
	IsBuyerMaker     bool   `json:"isBuyerMaker"`
	LastTradeID      int64  `json:"lastTradeID"`
	Price            string `json:"price"`
	Quantity         string `json:"quantity"`
	Timestamp        int64  `json:"timestamp"`
}
