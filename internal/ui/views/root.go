package views

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/authentication_app/internal/domain/request"
	"github.com/authentication_app/internal/domain/response"
	"github.com/authentication_app/internal/net"
	"github.com/authentication_app/internal/ui/constants"
	"github.com/sirupsen/logrus"
	"strconv"
)

var a = app.New()
var windows = map[request.MessageType]fyne.Window{}
var tables = map[int]*widget.Table{}

func messageListener() {
	for {
		select {
		case message, ok := <-net.Client.Receive():
			if !ok {
				continue
			}
			resp := &response.BaseResp{}
			_ = json.Unmarshal(message, resp)
			if resp.Code != "200" {
				dialog.ShowError(fmt.Errorf("%s", resp.Message), windows[resp.MessageType])
			} else {
				processMessage(resp.MessageType, message)
			}
		}
	}
}

func processMessage(messageType request.MessageType, message []byte) {
	switch messageType {
	case request.MessageAdminLogin:
		{
			logrus.Info("用户登陆成功，展示主页面")
			go net.Client.Heartbeat()
			windows[messageType].Hide()
			NewUserView()
		}
	case request.MessageUserInfo:
		{
			logrus.Info("获取用户信息成功")
			resp := response.UserInfoResponse{}
			_ = json.Unmarshal(message, &resp)
			userData[1][0] = strconv.FormatInt(int64(resp.RegisterCount), 10)
			userData[1][1] = strconv.FormatInt(int64(resp.OnlineCount), 10)
			userData[1][2] = strconv.FormatInt(int64(resp.OfflineCount), 10)
			tables[constants.KUserInfoTable].Refresh()
		}
	case request.MessageAdminGenerateCode:
		{
			logrus.Info("服务器生成编码成功")
			codeData = [][]string{[]string{constants.KCode}}
			resp := response.GenerateCodeResponse{}
			_ = json.Unmarshal(message, &resp)
			for _, v := range resp.Codes {
				codeData = append(codeData, []string{v})
			}
			tables[constants.KCodeTable].Refresh()
		}
	}
}
