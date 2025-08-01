package cus

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewLayout(w, h float32) *fyne.Container {
	return container.NewGridWrap(fyne.NewSize(w, h))
}
