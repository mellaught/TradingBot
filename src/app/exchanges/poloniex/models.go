package poloniex

import (
	"os"
	"time"

	polo "github.com/iowar/poloniex"
	"github.com/sirupsen/logrus"
)

var (
	PoloniexCandlestickIntervalList = []int{
		300, 900, 1800, 7200, 14400, 86400,
	}
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
	WsCli           *polo.WSClient
	quit            chan os.Signal
}
