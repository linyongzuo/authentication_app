package net

import (
	"encoding/json"
	"github.com/authentication_app/internal/domain/request"
	"github.com/authentication_app/internal/domain/ui_message"
	"github.com/sirupsen/logrus"
)

var Client = NewWsClient()

func Login(login ui_message.Login) (err error) {
	msg, _ := json.Marshal(request.AdminLoginReq{
		Header: request.Header{
			Version:     "3",
			MessageType: 2,
			Mac:         "",
			Ip:          "",
		},
		UserName: login.UserName,
		Password: login.Password,
	})
	err = Client.SendMessage(msg)
	if err != nil {
		logrus.Errorf("发送消息失败:%s", err.Error())
		return
	}
	return
}

func UserInfo() (err error) {
	msg, _ := json.Marshal(request.UserInfoReq{
		Header: request.Header{
			Version:     "3",
			MessageType: request.MessageUserInfo,
		},
	})
	err = Client.SendMessage(msg)
	if err != nil {
		logrus.Errorf("发送消息失败:%s", err.Error())
		return
	}
	return
}
func GenerateCode(count int) (err error) {
	msg, _ := json.Marshal(request.GenerateCodeReq{
		Header: request.Header{
			Version:     "3",
			MessageType: request.MessageAdminGenerateCode,
		},
		Count: count,
	})
	err = Client.SendMessage(msg)
	if err != nil {
		logrus.Errorf("发送消息失败:%s", err.Error())
		return
	}
	return
}
