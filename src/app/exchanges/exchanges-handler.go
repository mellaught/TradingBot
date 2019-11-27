package exchanges

import (
	"TradingBot/src/app/tickers"
)

// Exchange -- interface for all exchanges
type Exchange interface {
	GetTickersFromChan() *tickers.Tickers
	GetName() string
}

// Return tickers format models.
func GetTickets(e Exchange) *tickers.Tickers {
	return e.GetTickersFromChan()
}
