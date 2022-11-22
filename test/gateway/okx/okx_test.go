package okx

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

type ArgItem struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type Op struct {
	Op   string    `json:"op"`
	Args []ArgItem `json:"args"`
}

var addr = flag.String("addr", "ws.okx.com:8443", "http service address")

func TestNewOKExSpotV3Ws(t *testing.T) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/ws/v5/public"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	item := ArgItem{Channel: "tickers", InstId: "STARL-USDT"}

	items := []ArgItem{item}

	//发送订阅消息
	m := Op{Op: "subscribe", Args: items}
	bytes, _ := json.Marshal(m)

	err = c.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println("write:", err)
		return
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println(t)
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			//if err != nil {
			//	log.Println("write:", err)
			//	return
			//}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
