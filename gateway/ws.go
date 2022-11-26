package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"leopard-quant/util/recws"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Cmd struct {
	Op   string    `json:"op"`
	Args []ArgItem `json:"args"`
}

type ArgItem struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type Configuration struct {
	Addr          string `json:"addr"`
	Proxy         string `json:"proxy"`
	ApiKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	AutoReconnect bool   `json:"auto_reconnect"`
	DebugMode     bool   `json:"debug_mode"`
}

type WS struct {
	cfg           *Configuration
	ctx           context.Context
	cancel        context.CancelFunc
	conn          *recws.RecConn
	mu            sync.RWMutex
	Ended         bool
	subscribeCmds []Cmd
	callbacks     ComposeCallback
}

func New(config *Configuration, callbacks ComposeCallback) *WS {
	b := &WS{
		cfg:       config,
		callbacks: callbacks,
	}
	b.ctx, b.cancel = context.WithCancel(context.Background())

	b.conn = &recws.RecConn{
		KeepAliveTimeout: 60 * time.Second,
		NonVerbose:       true,
	}
	if config.Proxy != "" {
		proxy, err := url.Parse(config.Proxy)
		if err != nil {
			return nil
		}
		b.conn.Proxy = http.ProxyURL(proxy)
	}
	b.conn.SubscribeHandler = b.subscribeHandler
	return b
}

func (b *WS) subscribeHandler() error {
	if b.cfg.DebugMode {
		color.Greenln("Ws subscribeHandler")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, cmd := range b.subscribeCmds {
		err := b.SendCmd(cmd)
		if err != nil {
			color.Redln(fmt.Sprintf("Ws SendCmd return error: %v", err))
		}
	}

	return nil
}

func (b *WS) SendCmd(cmd Cmd) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return b.Send(string(data))
}

func (b *WS) Send(msg string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Ws send error: %v", r))
		}
	}()

	err = b.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return
}

func (b *WS) Start() error {
	b.connect()

	cancel := make(chan struct{})

	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				b.ping()
			case <-cancel:
				return
			}
		}
	}()

	go func() {
		defer close(cancel)

		for {
			messageType, data, err := b.conn.ReadMessage()
			if err != nil {
				color.Redln(fmt.Sprintf("Ws Read error, closing connection: %v", err))
				b.conn.Close()
				b.Ended = true
				return
			}

			b.processMessage(messageType, data)
		}
	}()

	return nil
}

func (b *WS) connect() {
	b.conn.Dial(b.cfg.Addr, nil)
}

func (b *WS) Subscribe(item ArgItem) {
	cmd := Cmd{
		Op:   "subscribe",
		Args: []ArgItem{item},
	}
	b.subscribeCmds = append(b.subscribeCmds, cmd)
	b.SendCmd(cmd)
}

func (b *WS) ping() {
	defer func() {
		if r := recover(); r != nil {
			color.Redln(fmt.Sprintf("Ws ping error: %v", r))
		}
	}()

	if !b.IsConnected() {
		return
	}
	err := b.conn.WriteMessage(websocket.TextMessage, []byte(`ping`))
	if err != nil {
		color.Redln(fmt.Sprintf("Ws ping error: %v", err))
	}
}

// IsConnected returns the WebSocket connection state
func (b *WS) IsConnected() bool {
	return b.conn.IsConnected()
}

func processMessagePipeline(data []byte, callback ComposeCallback) {

}

func (b *WS) processMessage(messageType int, data []byte) {

	s := string(data)
	if b.cfg.DebugMode {
		color.Greenln(fmt.Sprintf("Ws %v", s))
	}

	if s == "pong" {
		if b.cfg.DebugMode {
			color.Greenln("应答pong，忽略。")
		}
		return
	}

	ret := gjson.ParseBytes(data)

	if eventValue := ret.Get("event"); eventValue.Exists() {
		if b.cfg.DebugMode {
			color.Greenln("订阅应答，忽略。")
		}
		return
	}

	if channelType := ret.Get("arg.channel"); channelType.Exists() && channelType.Str == "tickers" {
		raw := ret.Get("data").Array()[0].Raw
		if b.cfg.DebugMode {
			color.Greenln(fmt.Sprintf("ticker data -> %s", raw))
		}
		if ticker, err := b.callbacks.TickerCallbackConverter(data, ret); err == nil {
			b.callbacks.TickerCallback(ticker)
		}
		return
	}

}
