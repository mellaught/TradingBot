package tickers

// Tickers have specified format for strategies.
type Tickers struct {
	Currencies      []string
	CurrencyToIndex map[string]int

	Tickers map[TickerKey]Ticker
}

// tickerKey contains base and quote currency
type TickerKey struct {
	Base  string
	Quote string
}

// Ticker contains buy price and sell price in decimal
type Ticker struct {
	BuyPrice  float64
	SellPrice float64
}

// Create New Ticker struct for stragies
func NewTickers() *Tickers {
	return &Tickers{
		Tickers:         make(map[TickerKey]Ticker),
		CurrencyToIndex: make(map[string]int),
	}
}

// //
// func (t *Tickers) Add(base, quote string, ticker Ticker) {
// 	t.putCurrency(base)
// 	t.putCurrency(quote)
// 	tk := TickerKey{Base: base, Quote: quote}
// 	t.Tickers[tk] = ticker
// }

// func (t *Tickers) putCurrency(currency string) {
// 	_, ok := t.CurrencyToIndex[currency]
// 	if !ok {
// 		t.Currencies = append(t.Currencies, currency)
// 		t.CurrencyToIndex[currency] = len(t.Currencies) - 1
// 	}
// }
