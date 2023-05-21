package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/authentication_app/internal/controller"
	"github.com/authentication_app/internal/domain/request"
	"github.com/authentication_app/internal/net"
	"github.com/authentication_app/internal/ui/constants"
	"strconv"
)

var userData = [][]string{[]string{"注册人数", "在线人数", "离线人数"}, []string{"0", "0", "0"}}
var codeData = [][]string{[]string{"编码"}, []string{""}}

func NewUserView() {
	w := a.NewWindow(constants.KMainText)
	windows[request.MessageUserInfo] = w
	w.Resize(fyne.Size{
		Width:  constants.KMainWindowWidth,
		Height: constants.KMainWindowHeight,
	})
	// 用户信息展示
	_ = net.UserInfo()
	userTable := widget.NewTable(func() (int, int) {
		return len(userData), len(userData[0])
	},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(userData[i.Row][i.Col])
		})
	tables[constants.KUserInfoTable] = userTable
	refreshButton := widget.NewButton(constants.KRefresh, func() {
		err := net.UserInfo()
		if err != nil {
			dialog.ShowError(err, w)
		}
	})
	bottomBox := container.NewVBox(refreshButton)
	userHSplit :=
		container.NewVSplit(userTable, bottomBox)
	user := container.NewTabItem(constants.KUserInfoText, userHSplit)

	codeTable := widget.NewTable(func() (int, int) {
		return len(codeData), 1
	},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(codeData[i.Row][i.Col])
		})
	codeCountEntry := widget.NewEntry()
	codeCountEntry.SetPlaceHolder(constants.KCodeCount)
	form := widget.NewForm(
		widget.NewFormItem(constants.KCodeCount, codeCountEntry),
	)
	form.OnSubmit = func() {
		err := codeCountEntry.Validate()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		count, _ := strconv.ParseInt(codeCountEntry.Text, 10, 64)
		err = net.GenerateCode((int)(count))
		if err != nil {
			dialog.ShowError(err, w)
		}
	}
	form.SubmitText = constants.KGenText
	saveButton := widget.NewButton(constants.KSave, func() {
		controller.SaveFile(codeData)
	})
	hBox := container.NewVBox(form, saveButton)
	box := container.NewHSplit(
		hBox,
		codeTable,
	)
	code := container.NewTabItem(constants.KGenTabText, box)
	tables[constants.KCodeTable] = codeTable

	appTabs := container.NewAppTabs(user, code)
	w.SetContent(appTabs)
	w.CenterOnScreen()
	w.Show()
}
