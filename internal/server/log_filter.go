package server

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const Max_Line = 10

var LogFilter = new(logFilter)

type logFilter struct {
	entry *widget.Entry
}

// 注册组件
func (h *logFilter) Register(entry *widget.Entry) {
	h.entry = entry
}

func (h *logFilter) Log(method, msg string) {
	go func() {
		if h.entry != nil {
			msg = fmt.Sprintf("%s: %s", method, msg)

			text := h.entry.Text
			if len(text) > Max_Line {
				text = text[:Max_Line]
			}

			fyne.Do(func() {
				h.entry.SetText(text + "\n" + msg)
			})
			h.entry.CursorRow = len(h.entry.Text)
		}
	}()
}
