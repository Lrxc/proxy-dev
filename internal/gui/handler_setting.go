package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"proxy-dev/assets"
	"proxy-dev/internal/config"
	"proxy-dev/internal/server"
	"proxy-dev/internal/system"
)

func settingProxyOnClick(myWindow fyne.Window, itme *fyne.MenuItem) func() {
	return func() {
		itme.Checked = !itme.Checked

		var err error
		if itme.Checked {
			if AppRunning {
				err = system.SysProxyOn()
				Proxy_Status.Set(fmt.Sprintf("%s(:%d)", PROXY_STATUS_RUNNING, config.Conf.Proxy.Port))
			}
		} else {
			err = system.SysProxyOff()
			Proxy_Status.Set(PROXY_STATUS_OFF)
		}

		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		config.Conf.System.AutoProxy = itme.Checked
		config.WriteConf(config.Conf)
	}
}

func settingHttpsOnClick(myWindow fyne.Window, itme *fyne.MenuItem) func() {
	return func() {
		itme.Checked = !itme.Checked

		err := server.ReStartServer(itme.Checked)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		if itme.Checked {
			Https_Status.Set(PROXY_STATUS_RUNNING)
		} else {
			Https_Status.Set(PROXY_STATUS_OFF)
		}

		config.Conf.System.Https = itme.Checked
		config.WriteConf(config.Conf)
	}
}

func settingExitOnClick(myWindow fyne.Window, itme *fyne.MenuItem) func() {
	return func() {
		itme.Checked = !itme.Checked

		config.Conf.System.MinExit = itme.Checked
		config.WriteConf(config.Conf)
	}
}

func settingInstallCa(myWindow fyne.Window) func() {
	return func() {
		file, err := assets.ReadFile("server.crt")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		err = system.InstallCert(file)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
	}
}
