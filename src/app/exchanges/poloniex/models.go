package poloniex

import (
	"os"
	"time"

	"github.com/shopspring/decimal"

	polo "github.com/iowar/poloniex"
	"github.com/sirupsen/logrus"
)

var (
	pushAPIUrl = "wss://api2.poloniex.com/realm1"

	PoloniexCandlestickIntervalList = []int{
		300, 900, 1800, 7200, 14400, 86400,
	}

	ConnectError     = "[ERROR] Connection could not be established!"
	RequestError     = "[ERROR] NewRequest Error!"
	SetApiError      = "[ERROR] Set the API KEY and API SECRET!"
	PeriodError      = "[ERROR] Invalid Period!"
	TimePeriodError  = "[ERROR] Time Period incompatibility!"
	TimeError        = "[ERROR] Invalid Time!"
	StartTimeError   = "[ERROR] Start Time Format Error!"
	EndTimeError     = "[ERROR] End Time Format Error!"
	LimitError       = "[ERROR] Limit Format Error!"
	ChannelError     = "[ERROR] Unknown Channel Name: %s"
	SubscribeError   = "[ERROR] Already Subscribed!"
	WSTickerError    = "[ERROR] WSTicker Parsing %s"
	WSOrderBookError = "[ERROR] MarketUpdate OrderBook Parsing %s"
	NewTradeError    = "[ERROR] MarketUpdate NewTrade Parsing %s"
	ServerError      = "[SERVER ERROR] Response: %s"
)

type Config struct {
	RequestInterval string `json:"request_interval"`
}

type PoloniexWorker struct {
	log             *logrus.Logger
	requestInterval time.Duration
	Tickers         map[string]*polo.Ticker
	CandleStick     []polo.CandleStick
	symbols         []string
	PubCli          *polo.Poloniex
	WsTickers       *WSClient
	WsMarkets       *WSClient
	quit            chan os.Signal
}

// for ticker update.
type WSTicker struct {
	Symbol        string          `json:"symbol"`
	Last          decimal.Decimal `json:"last,float"`
	LowestAsk     decimal.Decimal `json:"lowestAsk,float"`
	HighestBid    decimal.Decimal `json:"hihgestBid,float"`
	PercentChange decimal.Decimal `json:"percentChange,float"`
	BaseVolume    decimal.Decimal `json:"baseVolume,float"`
	QuoteVolume   decimal.Decimal `json:"quoteVolume,float"`
	IsFrozen      bool            `json:"isFrozen"`
	High24hr      decimal.Decimal `json:"high24hr,float"`
	Low24hr       decimal.Decimal `json:"low24hr,float"`
}

// for market update.
type MarketUpdate struct {
	Data       interface{}
	TypeUpdate string `json:"type"`
}

// "i" messages.
type OrderDepth struct {
	Symbol    string `json:"symbol"`
	OrderBook struct {
		Asks []Book `json:"asks"`
		Bids []Book `json:"bids"`
	} `json:"orderBook"`
}

type Book struct {
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
}

// "o" messages
type WSOrderBook struct {
	Rate      decimal.Decimal `json:"rate,string"`
	TypeOrder string          `json:"type"`
	Amount    decimal.Decimal `json:"amount,string"`
}

// "o" messages.
type WSOrderBookModify WSOrderBook

// "o" messages.
type WSOrderBookRemove struct {
	Rate      decimal.Decimal `json:"rate,string"`
	TypeOrder string          `json:"type"`
}

// "t" messages.
type NewTrade struct {
	TradeId   int64           `json:"tradeID,string"`
	Rate      decimal.Decimal `json:"rate,string"`
	Amount    decimal.Decimal `json:"amount,string"`
	Total     decimal.Decimal `json:"total,string"`
	TypeOrder string          `json:"type"`
}
