package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/energye/systray"
	log "github.com/sirupsen/logrus"
	"image/color"
	"proxy-dev/assets"
	"proxy-dev/internal/config"
	"proxy-dev/internal/gui/cus"
	"proxy-dev/internal/system"
	"runtime"
	"syscall"
)

func Gui() {
	myApp := app.New()
	myApp.SetIcon(&fyne.StaticResource{StaticContent: assets.Read("logo.png")})
	myWindow := myApp.NewWindow(config.AppName)          //主窗口
	myWindow.Resize(fyne.NewSize(APP_WIDTH, APP_HEIGHT)) //窗口大小
	myWindow.CenterOnScreen()                            //窗口居中
	myWindow.SetFixedSize(true)                          //禁止最大化窗口

	content := initView(myApp, myWindow)
	initTray(myApp, myWindow)
	go initTask(myWindow)

	config.AppWindow = myWindow

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func initTray(myApp fyne.App, myWindow fyne.Window) {
	// 创建自定义托盘图标
	//icon := &fyne.StaticResource{
	//	StaticContent: assets.Read("logo.png"),
	//}
	//
	//if desk, ok := myApp.(desktop.App); ok {
	//	desk.SetSystemTrayIcon(icon)
	//
	//	menu := fyne.NewMenu("我的应用",
	//		fyne.NewMenuItem("显示", func() { myWindow.Show() }),
	//		fyne.NewMenuItem("退出", func() { myApp.Quit() }),
	//	)
	//	desk.SetSystemTrayMenu(menu)
	//}

	onReady := func() {
		systray.SetIcon(assets.Read("tray.ico"))
		systray.SetTitle(config.AppName)
		systray.SetTooltip(config.AppName)

		systray.SetOnClick(func(menu systray.IMenu) {
			fyne.Do(func() {
				myWindow.Show()
			})
		})
		systray.SetOnRClick(func(menu systray.IMenu) {
			menu.ShowMenu()
		})

		mShow := systray.AddMenuItem("打开", "打开")
		mQuit := systray.AddMenuItem("退出", "退出")
		mShow.Click(func() {
			fyne.Do(func() {
				myWindow.Show()
			})
		})
		mQuit.Click(func() {
			system.SigChan <- syscall.SIGTERM
			fyne.Do(func() {
				myApp.Quit()
			})
			systray.Quit()
		})
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("%s not support tray", runtime.GOOS)
			}
		}()

		systray.Run(onReady, nil)
	}()

	// 拦截关闭事件
	myWindow.SetCloseIntercept(func() {
		if config.Conf.System.MinExit {
			myWindow.Hide()
		} else {
			system.SigChan <- syscall.SIGTERM
			myApp.Quit()
			systray.Quit()
		}
	})
}

// 选择的打包路径
var (
	Proxy_Status = binding.NewString()
	Https_Status = binding.NewString()
	caBtn        *widget.Button
)

func initView(myApp fyne.App, myWindow fyne.Window) *fyne.Container {
	setItem := widget.NewToolbarAction(theme.SettingsIcon(), nil)
	setItem.OnActivated = settingOnClick(myWindow, setItem)

	helpItem := widget.NewToolbarAction(theme.HelpIcon(), nil)
	helpItem.OnActivated = helpOnClick(myWindow, helpItem)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FileTextIcon(), editRuleOnClick(myApp)),
		widget.NewToolbarSpacer(),
		setItem,
		helpItem,
	)

	//myWindow.SetMainMenu(makeMenu(myApp, myWindow))

	//ca证书提示
	caBtn = widget.NewButton(CA_STATUS, nil)
	caBtn.OnTapped = caInsOnClick(myWindow, caBtn)
	caBtn.Importance = widget.WarningImportance
	caBtn.Hide()

	// 创建一个标签
	proxyTitle := widget.NewLabel(PROXY_TITLE)
	Proxy_Status.Set(PROXY_STATUS_OFF)
	// 创建一个标签(绑定数据,自动更新)
	proxyLabel := widget.NewLabelWithData(Proxy_Status)
	proxyLabel.Importance = widget.WarningImportance

	// 创建一个标签
	httpsTitle := widget.NewLabel(HTTPS_TITLE)
	Https_Status.Set(PROXY_STATUS_OFF)
	// 创建一个标签(绑定数据,自动更新)
	httpsLabel := widget.NewLabelWithData(Https_Status)
	httpsLabel.Importance = widget.WarningImportance

	//开始按钮
	startBtn := widget.NewButton(PROXY_BTN_START, nil)
	startBtn.OnTapped = startOnClick(myWindow, startBtn)

	//横线
	thickLine := canvas.NewRectangle(color.NRGBA{R: 128, G: 128, B: 128, A: 255})

	//entry := widget.NewMultiLineEntry()
	//entry.Wrapping = fyne.TextWrapWord //启用自动换行
	//server.LogFilter.Register(entry)

	top := container.NewVBox(
		toolbar,
		thickLine,
		cus.NewLayout(0, 10),
		container.NewHBox(layout.NewSpacer(), caBtn, layout.NewSpacer()),
		cus.NewLayout(0, 20),
		container.NewHBox(cus.NewLayout(120, 0), proxyTitle, proxyLabel), //嵌套一个水平布局,并且居中
		container.NewHBox(cus.NewLayout(120, 0), httpsTitle, httpsLabel), //嵌套一个水平布局,并且居中
		cus.NewLayout(0, 50),
	)
	content := container.NewBorder(
		top,
		//container.NewVBox(cus.NewLayout(0, 50), entry),
		cus.NewLayout(0, 100),
		cus.NewLayout(100, 0),
		cus.NewLayout(100, 0),
		startBtn,
	)
	return content
}
