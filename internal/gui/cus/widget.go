package cus

import (
	"fyne.io/fyne/v2/widget"
)

type Button struct {
	widget.Button
}

func NewButton(label string, tapped func()) *Button {
	button := &Button{}
	button.Text = label
	button.OnTapped = tapped

	button.ExtendBaseWidget(button)
	return button
}
