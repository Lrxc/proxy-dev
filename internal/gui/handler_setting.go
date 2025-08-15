package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"net/url"
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

		if itme.Checked {
			Https_Status.Set(PROXY_STATUS_RUNNING)
		} else {
			Https_Status.Set(PROXY_STATUS_OFF)
		}

		config.Conf.System.Https = itme.Checked
		config.WriteConf(config.Conf)

		err := server.ReDeploy(itme.Checked)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
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

func about(myWindow fyne.Window) func() {
	return func() {
		// 创建超链接组件（蓝色字体）
		url, _ := url.Parse("https://github.com/Lrxc/proxy-dev")
		link := widget.NewHyperlink("Github", url)
		link.Alignment = fyne.TextAlignCenter // 居中显示

		// 创建对话框内容
		content := widget.NewLabel(config.AppName)
		container := container.NewVBox(
			content,
			link,
		)

		// 显示对话框
		dialog.ShowCustom(
			"",
			"OK",
			container,
			myWindow,
		)
	}
}
