package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"proxy-dev/internal/config"
	"proxy-dev/internal/util"
)

func editRuleOnClick(myApp fyne.App) func() {
	return func() {
		//打开一个新窗口
		newWin := myApp.NewWindow(config.AppName)
		newWin.Resize(fyne.NewSize(APP_WIDTH, APP_HEIGHT))
		newWin.CenterOnScreen() //居中显示

		entry := widget.NewMultiLineEntry()
		entry.SetText(util.PrettyJSON(config.Conf.Rule))
		entry.Wrapping = fyne.TextWrapWord //启用自动换行

		saveBtn := widget.NewButton("保存", func() {
			err := config.WriteJson(entry.Text)
			if err != nil {
				dialog.ShowError(err, newWin)
			} else {
				dialog.ShowInformation("success", "", newWin)
			}
		})
		refreshBtn := widget.NewButton("格式化", func() {
			entry.SetText(util.PrettyJSON(entry.Text))
		})

		btn := container.NewHBox(saveBtn, refreshBtn)

		content := container.NewBorder(btn, nil, nil, nil, entry)
		newWin.SetContent(content)
		//newWin.SetContent(makeListTab())
		newWin.Show()
	}
}

// list列表
func makeListTab() fyne.CanvasObject {
	data := binding.BindStringList(
		&[]string{"Item 1", "Item 2", "Item 3"},
	)

	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			box := container.NewVBox(
				container.NewHBox(widget.NewLabel("原始地址"), widget.NewEntry()),
				container.NewHBox(widget.NewLabel("目标地址"), widget.NewEntry()),
			)
			return container.NewHBox(widget.NewCheck("", nil), box)
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			//o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	selectItem := ""
	list.OnSelected = func(id widget.ListItemID) {
		selectItem, _ = data.GetValue(id)
		fmt.Println(data.GetValue(id))
	}

	saveBtn := widget.NewButton("添加", func() {
		data.Append("1")
	})
	refreshBtn := widget.NewButton("删除", func() {
		data.Remove(selectItem)
	})

	btn := container.NewHBox(saveBtn, refreshBtn)

	content := container.NewBorder(btn, nil, nil, nil, list)
	return content
}
