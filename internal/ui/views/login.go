package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/authentication_app/internal/domain/request"
	"github.com/authentication_app/internal/domain/ui_message"
	"github.com/authentication_app/internal/net"
	"github.com/authentication_app/internal/ui/constants"
	"github.com/authentication_app/internal/ui/theme"
)

func NewLoginView() {
	go messageListener()
	w := a.NewWindow(constants.KLoginWindowText)
	windows[request.MessageAdminLogin] = w
	a.Settings().SetTheme(&theme.CustomerTheme{})
	w.Resize(fyne.Size{
		Width:  constants.KLoginWindowWidth,
		Height: constants.KLoginWindowHeight,
	})
	userNameEntry := widget.NewEntry()
	userNameEntry.SetPlaceHolder(constants.KUserNameText)
	userNameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	userPasswordEntry := widget.NewPasswordEntry()
	userPasswordEntry.SetPlaceHolder(constants.KUserPasswordText)
	userPasswordEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	form := widget.NewForm(
		widget.NewFormItem(constants.KUserNameText, userNameEntry),
		widget.NewFormItem(constants.KUserPasswordText, userPasswordEntry),
	)
	form.OnCancel = func() {
		w.Close()
	}
	form.CancelText = constants.KCancelText
	form.OnSubmit = func() {
		err := userNameEntry.Validate()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		err = net.Login(ui_message.Login{
			MessageType: ui_message.MessageLogin,
			UserName:    userNameEntry.Text,
			Password:    userPasswordEntry.Text,
		})
		if err != nil {
			dialog.ShowError(err, w)
		}
	}
	form.SubmitText = constants.KLoginButtonText
	w.SetFixedSize(true)
	w.SetContent(container.NewVBox(
		form,
	))
	w.CenterOnScreen()
	w.ShowAndRun()
}
