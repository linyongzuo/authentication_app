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
	"strconv"
)

var a = app.New()
var windows = map[request.MessageType]fyne.Window{}
var tables = map[int]*widget.Table{}

func messageListener() {
	for {
		select {
		case message, ok := <-net.Client.Receive:
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
			net.Client.Heartbeat()
			windows[messageType].Hide()
			NewUserView()
		}
	case request.MessageUserInfo:
		{
			resp := response.UserInfoResponse{}
			_ = json.Unmarshal(message, &resp)
			userData[1][0] = strconv.FormatInt(int64(resp.RegisterCount), 10)
			userData[1][1] = strconv.FormatInt(int64(resp.OnlineCount), 10)
			userData[1][2] = strconv.FormatInt(int64(resp.OfflineCount), 10)
			tables[constants.KUserInfoTable].Refresh()
		}
	case request.MessageAdminGenerateCode:
		{
			codeData = [][]string{[]string{"编码"}}
			resp := response.GenerateCodeResponse{}
			_ = json.Unmarshal(message, &resp)
			for _, v := range resp.Codes {
				codeData = append(codeData, []string{v})
			}
			tables[constants.KCodeTable].Refresh()
		}
	}
}
