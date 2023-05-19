package net

import (
	"encoding/json"
	"fmt"
	"github.com/authentication_app/global"
	"github.com/authentication_app/internal/domain/request"
	"github.com/authentication_app/internal/ui/constants"
	"github.com/authentication_app/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/url"
	"sync/atomic"
	"time"
)

const (
	heartbeatTime = 5 * time.Second
	connectTime   = 15 * time.Second
)

type WsClient struct {
	Address           string
	c                 *websocket.Conn
	Receive           chan []byte
	stopReconnectChan chan struct{}
	isConnect         atomic.Bool
	mac               string
}

func NewWsClient() *WsClient {
	resp := &WsClient{}
	resp.Receive = make(chan []byte, 1024)
	resp.stopReconnectChan = make(chan struct{})
	resp.isConnect = atomic.Bool{}
	resp.mac = utils.GetMac()
	resp.connect()
	go resp.reconnect()
	return resp
}

func (w *WsClient) readMessage() {
	done := make(chan struct{})
	go func() {
		defer func() {
			close(done)
			_ = w.c.Close()
			w.isConnect.Swap(false)
		}()
		for {
			_, resp, err := w.c.ReadMessage()
			if err != nil {
				logrus.Error("读取数据失败:", err)
				return
			}
			logrus.Infof("收到服务器发送消息: %s", resp)
			w.Receive <- resp
		}
	}()

}
func (w *WsClient) Close() {
	_ = w.c.Close()
	w.isConnect.Swap(false)
	close(w.Receive)
	close(w.stopReconnectChan)
}

func (w *WsClient) Heartbeat() {
	req := request.HeartbeatReq{
		Header: request.Header{
			Version:     constants.KVersion,
			MessageType: request.MessageHeartbeat,
			Mac:         w.mac,
		},
		Admin: true,
	}
	msg, _ := json.Marshal(&req)
	_ = w.SendMessage(msg)
}
func (w *WsClient) connect() {
	logrus.Info("开始连接服务器")
	if !w.isConnect.Load() {
		address := fmt.Sprintf("%s:%d", global.Cfg.ServerCfg.Host, global.Cfg.ServerCfg.Port)
		u := url.URL{Scheme: "ws", Host: address, Path: "/authentication_hub/connect"}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			logrus.Errorf("连接服务器失败:%s", err.Error())
			return
		} else {
			logrus.Info("连接成功")
			w.isConnect.Swap(true)
			w.c = c
			w.Address = address
			w.readMessage()
		}
	}
}
func (w *WsClient) reconnect() {
	ticker := time.NewTicker(connectTime)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-w.stopReconnectChan:
			return
		case <-ticker.C:
			w.connect()
		}
	}
}
func (w *WsClient) SendMessage(message []byte) (err error) {
	err = w.c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return
	}
	return
}
