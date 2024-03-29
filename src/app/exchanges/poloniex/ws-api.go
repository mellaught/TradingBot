package poloniex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/gorilla/websocket"
	polo "github.com/iowar/poloniex"
)

const (
	TICKER     = 1002 // Ticker Channel Id
	SUBSBUFFER = 24   // Subscriptions Buffer
)

var (
	channelsByName = make(map[string]int) // channels map by name
	channelsByID   = make(map[int]string) // channels map by id
	marketChannels []int                  // channels list
)

type WSClient struct {
	Subs       map[string]chan interface{} // subscriptions map
	wsConn     *websocket.Conn             // websocket connection
	wsMutex    *sync.Mutex                 // prevent race condition for websocket RW
	sync.Mutex                             // embedded mutex
}

// subscription and unsubscription
type subscription struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
}

// Web socket reader.
func (ws *WSClient) readMessage() ([]byte, error) {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()
	_, rmsg, err := ws.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return rmsg, nil
}

// Web socket writer.
func (ws *WSClient) writeMessage(msg []byte) error {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()
	return ws.wsConn.WriteMessage(1, msg)
}

// Set channels.
func setChannelsId() (err error) {
	publicApi, err := polo.NewClient("", "")
	if err != nil {
		return err
	}

	tickers, err := publicApi.GetTickers()
	if err != nil {
		return err
	}

	for k, v := range tickers {
		channelsByName[k] = v.ID
		channelsByID[v.ID] = k
		marketChannels = append(marketChannels, v.ID)
	}

	channelsByName["TICKER"] = TICKER
	channelsByID[TICKER] = "TICKER"
	return
}

// Create new web socket client.
func NewWSClient() (wsClient *WSClient, err error) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Minute,
	}

	ws, _, err := dialer.Dial(pushAPIUrl, nil)
	if err != nil {
		return
	}

	wsClient = &WSClient{
		wsConn:  ws,
		Subs:    make(map[string]chan interface{}),
		wsMutex: &sync.Mutex{},
	}

	if err = setChannelsId(); err != nil {
		return
	}

	go wsClient.wsHandler()

	return
}

// Create handler.
// If the message comes from the channels that are subscribed,
// it is sent to the chans.
func (ws *WSClient) wsHandler() {
	for {
		msg, err := ws.readMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		var imsg []interface{}
		err = json.Unmarshal(msg, &imsg)
		if err != nil || len(imsg) < 3 {
			continue
		}

		arg, ok := imsg[0].(float64)
		if !ok {
			continue
		}

		chid := int(arg)
		args, ok := imsg[2].([]interface{})
		if !ok {
			continue
		}

		var wsupdate interface{}
		if chid == TICKER {
			wsupdate, err = convertArgsToTicker(args)
			if err != nil {
				continue
			}
		} else if intInSlice(chid, marketChannels) {
			wsupdate, err = convertArgsToMarketUpdate(args)
			if err != nil {
				continue
			}
		} else {
			continue
		}

		chname := channelsByID[chid]
		if ws.Subs[chname] != nil {
			if chid == TICKER {
				ws.Subs[chname] <- wsupdate
			} else {
				ws.Subs[chname] <- wsupdate
			}
		}
	}
}

// Convert ticker update arguments and fill wsticker.
func convertArgsToTicker(args []interface{}) (wsticker WSTicker, err error) {
	wsticker.Symbol = channelsByID[int(args[0].(float64))]
	wsticker.Last, err = decimal.NewFromString(args[1].(string))
	if err != nil {
		err = Error(WSTickerError, "Last")
		return
	}
	wsticker.LowestAsk, err = decimal.NewFromString(args[2].(string))
	if err != nil {
		err = Error(WSTickerError, "LowestAsk")
		return
	}
	wsticker.HighestBid, err = decimal.NewFromString(args[3].(string))
	if err != nil {
		err = Error(WSTickerError, "HighestBid")
		return
	}

	wsticker.PercentChange, err = decimal.NewFromString(args[4].(string))
	if err != nil {
		err = Error(WSTickerError, "PercentChange")
		return
	}

	wsticker.BaseVolume, err = decimal.NewFromString(args[5].(string))
	if err != nil {
		err = Error(WSTickerError, "BaseVolume")
		return
	}

	wsticker.QuoteVolume, err = decimal.NewFromString(args[6].(string))
	if err != nil {
		err = Error(WSTickerError, "QuoteVolume")
		return
	}

	if v, ok := args[7].(float64); ok {
		if v == 0 {
			wsticker.IsFrozen = false
		} else {
			wsticker.IsFrozen = true
		}
	} else {
		err = Error(WSTickerError, "IsFrozen")
		return
	}

	wsticker.High24hr, err = decimal.NewFromString(args[8].(string))
	if err != nil {
		err = Error(WSTickerError, "High24hr")
		return
	}

	wsticker.Low24hr, err = decimal.NewFromString(args[9].(string))
	if err != nil {
		err = Error(WSTickerError, "Low24hr")
		return
	}

	return
}

// Convert market update arguments and fill marketupdate.
func convertArgsToMarketUpdate(args []interface{}) (res []MarketUpdate, err error) {
	res = make([]MarketUpdate, len(args))
	for i, val := range args {
		vals := val.([]interface{})
		var marketupdate MarketUpdate

		switch vals[0].(string) {
		// Init.
		// case "i":
		// 	var orderdepth OrderDepth
		// 	val := vals[1].(map[string]interface{})
		// 	orderdepth.Symbol = val["currencyPair"].(string)

		// 	asks := val["orderBook"].([]interface{})[0].(map[string]interface{})
		// 	bids := val["orderBook"].([]interface{})[1].(map[string]interface{})

		// 	for k, v := range bids {
		// 		price, _ := strconv.ParseFloat(k)
		// 		quantity, _ := strconv.ParseFloat(v.(string))
		// 		book := Book{Price: price, Quantity: quantity}
		// 		orderdepth.OrderBook.Bids = append(orderdepth.OrderBook.Bids, book)
		// 	}

		// 	for k, v := range asks {
		// 		price, _ := strconv.ParseFloat(k)
		// 		quantity, _ := strconv.ParseFloat(v.(string))
		// 		book := Book{Price: price, Quantity: quantity}
		// 		orderdepth.OrderBook.Asks = append(orderdepth.OrderBook.Asks, book)
		// 	}

		// 	marketupdate.TypeUpdate = "OrderDepth"
		// 	fmt.Println(len(orderdepth.OrderBook.Asks), len(orderdepth.OrderBook.Bids))
		// 	marketupdate.Data = orderdepth
		// Update Order Book.
		case "o":
			var orderdatafield WSOrderBook

			marketupdate.TypeUpdate = "modify"
			if t := vals[3].(string); t == "0.00000000" {
				marketupdate.TypeUpdate = "remove"
			}
			orderdatafield.TypeOrder = "ask"
			if t := vals[1].(float64); t == 1 {
				orderdatafield.TypeOrder = "bid"
			}

			orderdatafield.Rate, err = decimal.NewFromString(vals[2].(string))
			if err != nil {
				err = Error(WSOrderBookError, "Rate")
				return
			}

			orderdatafield.Amount, err = decimal.NewFromString(vals[3].(string))
			fmt.Println(vals[3])
			if err != nil {
				err = Error(WSOrderBookError, "Amount")
				return
			}

			marketupdate.Data = orderdatafield
		// Update Order Trade.
		case "t":

			var tradedatafield NewTrade

			tradedatafield.TradeId, err = strconv.ParseInt(vals[1].(string), 10, 64)
			if err != nil {
				err = Error(NewTradeError, "TradeId")
				return
			}

			if vals[2].(float64) == 1 {
				tradedatafield.TypeOrder = "buy"
			} else {
				tradedatafield.TypeOrder = "sell"
			}

			tradedatafield.Rate, err = decimal.NewFromString(vals[3].(string))
			if err != nil {
				err = Error(NewTradeError, "Rate")
				return
			}

			tradedatafield.Amount, err = decimal.NewFromString(vals[4].(string))
			if err != nil {
				err = Error(NewTradeError, "Amount")
				return
			}

			tradedatafield.Total = decimal.NewFromFloat(vals[5].(float64))

			marketupdate.TypeUpdate = "NewTrade"
			marketupdate.Data = tradedatafield
		}

		res[i] = marketupdate
	}

	return res, nil
}

// sub-function for subscription.
func (ws *WSClient) subscribe(chid int, chname string) (err error) {
	ws.Lock()
	defer ws.Unlock()

	//	if ws.Subs[chname] != nil {
	//		err = Error(SubscribeError)
	//		return
	//	}
	//
	//	ws.Subs[chname] = make(chan interface{}, SUBSBUFFER)

	if ws.Subs[chname] == nil {
		ws.Subs[chname] = make(chan interface{}, SUBSBUFFER)
	}

	subsMsg, _ := subscription{
		Command: "subscribe",
		Channel: strconv.Itoa(chid),
	}.toJSON()
	err = ws.writeMessage(subsMsg)
	if err != nil {
		return
	}

	return
}

// sub-function for unsubscription.
// the chans are not closed once the subscription is made to protect chan address.
// To prevent chans taking a new address on the memory, thus chans can be used repeatedly.
func (ws *WSClient) unsubscribe(chname string) (err error) {
	ws.Lock()
	defer ws.Unlock()

	if ws.Subs[chname] == nil {
		return
	}

	unSubsMsg, _ := subscription{
		Command: "unsubscribe",
		Channel: chname,
	}.toJSON()

	err = ws.writeMessage(unSubsMsg)
	if err != nil {
		return err
	}

	// close(ws.Subs[chname])
	// delete(ws.Subs, chname)
	return
}

// Subscribe to ticker channel.
// It returns nil if successful.
func (ws *WSClient) SubscribeTicker() error {
	return (ws.subscribe(TICKER, "TICKER"))
}

// Subscribe to Price channel.
// It returns nil if successful.
func (ws *WSClient) SubscribePrice(symbol string) error {

	return (ws.subscribe(TICKER, "TICKER"))
}

// Unsubscribe from ticker channel.
// It returns nil if successful.
func (ws *WSClient) UnsubscribeTicker() error {
	return (ws.unsubscribe("TICKER"))
}

// Subscribe to market channel.
// It returns nil if successful.
func (ws *WSClient) SubscribeMarket(chname string) error {
	chname = strings.ToUpper(chname)
	chid, ok := channelsByName[chname]
	if !ok {
		return Error(ChannelError, chname)
	}
	return (ws.subscribe(chid, chname))
}

// Unsubscribe from market channel.
// It returns nil if successful.
func (ws *WSClient) UnsubscribeMarket(chname string) error {
	chname = strings.ToUpper(chname)
	_, ok := channelsByName[chname]
	if !ok {
		return Error(ChannelError, chname)
	}
	return (ws.unsubscribe(chname))
}
