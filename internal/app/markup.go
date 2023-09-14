package app

import (
	"log"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateMarkup() {
	btn1 := widget.NewButton("Button 1", func() {})
	// btn1.Move()
	btn2 := widget.NewButton("Button 2", func() {})
	hBox1 := container.NewVScroll(btn1)
	hBox2 := container.NewVScroll(btn2)
	topBox := container.NewVBox(hBox1, hBox2)

	btn3 := widget.NewButton("Button 1", func() {})
	btn4 := widget.NewButton("Button 2", func() {})
	hBox3 := container.NewVScroll(btn3)
	hBox4 := container.NewVScroll(btn4)
	lowBox := container.NewVBox(hBox3, hBox4)

	log.Println(topBox, lowBox)
}
