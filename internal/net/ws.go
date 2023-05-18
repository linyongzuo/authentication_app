package net

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/url"
)

type WsClient struct {
	Address string
	c       *websocket.Conn
	Receive chan []byte
}

func NewWsClient() *WsClient {
	resp := &WsClient{}
	resp.Receive = make(chan []byte, 1024)
	address := "localhost:8090"
	u := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/authentication_hub/connect"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Errorf("连接服务器失败:%s", err.Error())
		go resp.connect()
	} else {
		resp.c = c
		resp.Address = address
		resp.readMessage()
	}
	return resp
}

func (w *WsClient) readMessage() {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, resp, err := w.c.ReadMessage()
			if err != nil {
				logrus.Error("read:", err)
				return
			}
			logrus.Infof("recv: %s", resp)
			w.Receive <- resp
		}
	}()
}
func (w *WsClient) Close() {

}

func (w *WsClient) Heartbeat() {

}
func (w *WsClient) connect() {

}

func (w *WsClient) SendMessage(message []byte) (err error) {
	err = w.c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return
	}
	return
}
