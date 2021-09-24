package exchange

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/jeongpope/go-crix/model"
)

const (
	signature_key = "SIGNATURE_JEONGPOPE"
)

// Handler handle raw websocket message
type Handler func(msg []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// Composition
type exchange struct {
	tickerEndpoint  string            // endpoint
	c               *websocket.Conn   // websocket connection
	reconnectLock   *sync.Mutex       // for single reconnect
	subsMessage     []*[]byte         // subscribe request message
	chanSendMessage chan model.Ticker // to send message channel

	supportAssets []string                // supported assets (uppercase)
	tickers       map[string]model.Ticker // availiable tickers
	updateLock    *sync.Mutex             // concurrent read/write
}

// Polymorphism
type IExchange interface {
	Initialize(currencies *[]string) []string // Initialize
	Execute() (err error)                     // Execute subscribe
	Release()                                 // Release memory

	initSnapshot([]string)
}

//
func (ex *exchange) AttatchChannel(ch chan model.Ticker) {
	ex.chanSendMessage = ch
}
