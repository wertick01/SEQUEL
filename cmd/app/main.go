package main

import (
	"biolink-gui/internal/models"
	"biolink-gui/internal/trimmomatic"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	App fyne.App
}

func NewAppStruct() *App {
	newApp := app.New()
	return &App{
		App: newApp,
	}
}

func main() {
	newApp := NewAppStruct()
	windowName := newApp.App.NewWindow("App")
	windowName.Resize(fyne.NewSize(400, 320))

	label := widget.NewLabel("Hello from FYNE")
	// entry := widget.NewEntry()
	// entryCopy := widget.NewEntry()

	// button := widget.NewButton("Print", func() {
	// 	log.Println(entry.Text, entryCopy.Text)
	// })

	trimm := trimmomatic.Trimmomatic{
		Params: &models.TrimmomaticParams{},
	}
	btn := trimm.ReadParams()
	log.Println(label, btn)

	windowName.SetContent(container.NewVBox(
		label,
		trimm.Path,
		trimm.Input,
		trimm.Output,
		trimm.Paired,
		trimm.Phred,
		trimm.Threads,
		btn,
	))

	// log.Println(trimm.Params.Path)

	windowName.ShowAndRun()
}
