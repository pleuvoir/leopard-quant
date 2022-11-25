package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"leopard-quant/util"
	"net/url"
)

type Websocket struct {
	Scheme string
	Host   string
	Path   string
	done   chan struct{}
	c      *websocket.Conn
}

func NewWebsocket(scheme string, host string, path string) *Websocket {
	w := Websocket{}
	w.Scheme = scheme
	w.Host = host
	w.Path = path
	w.done = make(chan struct{})
	return &w
}

// Connect 连接WebSocket
func (w *Websocket) Connect() error {
	u := url.URL{Scheme: w.Scheme, Host: w.Host, Path: w.Path}
	color.Greenln(fmt.Sprintf("connecting to %s", u.String()))
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		color.Redf("dial:", err)
		return err
	}
	w.c = c
	return nil
}

// SendRawTextMessage 发送完整文本
func (w *Websocket) SendRawTextMessage(m string) error {
	if err := w.c.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
		color.Redf("SendTextMessage:", err)
		return err
	}
	return nil
}

// SendJSONTextMessage 转换为JSON类型进行发送
func (w *Websocket) SendJSONTextMessage(m any) error {
	conn := w.c
	if conn == nil {
		return errors.Errorf("请先连接，再发送消息。")
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		color.Redf("SendTextMessage，转换JSON失败", err)
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
		color.Redf("SendTextMessage，通知对端失败", err)
		return err
	}
	return nil
}

// ReadMessage 读取消息
func (w *Websocket) ReadMessage() (p []byte, e error) {
	_, message, err := w.c.ReadMessage()
	return message, err
}

// Close 关闭连接
func (w *Websocket) Close() {
	util.CloseQuietly(w.c)
}
