package exchange

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jeongpope/go-crix/logger"
)

var (
	reconnectInterval = 10
)

var websocketServe = func(c *websocket.Conn, m *sync.Mutex,
	endpoint string, msg []*[]byte, handler Handler, errHandler ErrHandler) (err error) {
	c = connect(endpoint)

	go func() {
		log.Println("[websocket.go] websocketServe go-routine")
		go func() {
			if msg != nil {
				tTicker := time.NewTicker(time.Millisecond * 250)

				for _, value := range msg {
					c.WriteMessage(websocket.TextMessage, *value)
					<-tTicker.C
				}
			}
		}()

		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				logger.Log.Errorf("[%s] %s", endpoint, err.Error())

				c = reconnect(c, m, endpoint)
				if msg != nil {
					tTicker := time.NewTicker(time.Millisecond * 250)

					for _, value := range msg {
						c.WriteMessage(websocket.TextMessage, *value)
						<-tTicker.C
					}
				}

				continue
			}

			handler(message)
		}
	}()
	return
}

func connect(endpoint string) (conn *websocket.Conn) {
	log.Println("[websocket.go] connect")

	var err error
	for {
		conn, _, err = websocket.DefaultDialer.Dial(endpoint, nil)

		if err != nil {
			logger.Log.Error("Dial failed, retry 10 second")

			time.Sleep(time.Second * time.Duration(reconnectInterval))
		} else {
			break
		}
	}

	log.Println("[websocket.go] connect close")

	return conn
}

func reconnect(c *websocket.Conn, m *sync.Mutex, endpoint string) *websocket.Conn {
	log.Println("[websocket.go] reconnect")

	m.Lock()
	defer m.Unlock()

	c.Close()
	return connect(endpoint)
}
