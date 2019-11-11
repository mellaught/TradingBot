package binance

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
