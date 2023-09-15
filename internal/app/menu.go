package app

import (
	"biolink-nipt-gui/internal/trimmomatic"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateMainMenu(window fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, exitTerminal chan bool, displays map[int]map[string]int) *fyne.MainMenu {
	fileMenu := CreateFileItems(window, trimm, newApp)
	analysisMenu := CreateAnalysisItems(window, trimm, newApp, commandChan, exitTerminal)

	// commandsChan := make(chan string)
	// closeTerminal := make(chan bool, 1)
	// defer close(commandsChan)
	// defer close(closeTerminal)

	// newVBox2 := container.NewVBox()

	terminalMenu := CreateTerminalItems(window, newApp, commandChan, exitTerminal, float32(displays[0]["X"])*0.9, float32(displays[0]["Y"])*0.2)

	timeMenu := CreateTimeItems(newApp)

	menu := fyne.NewMainMenu(fileMenu, analysisMenu, terminalMenu, timeMenu)
	return menu
}

func CreateNewResearchForm(newApp *App, window fyne.Window, trimm *trimmomatic.Trimmomatic) *widget.Form {
	userName := widget.NewEntry()
	// userNameLabel := widget.NewLabel("")

	currentTime := time.Now().Format(time.RFC3339)
	currentDateTime := widget.NewEntryWithData(binding.BindString(&currentTime))
	// currentDateTimeLabel := widget.NewLabel("")

	researchName := widget.NewEntry()
	// researchNameLabel := widget.NewLabel("")

	outputDataFolderPath := ""
	outputDataFolder := widget.NewButton("Chose Output Dir", func() {
		dialog.ShowFolderOpen(
			func(r fyne.ListableURI, err error) {
				if r != nil {
					outputDataFolderPath = r.Path()
					// forwardSelected.SetText(r.URI().Path())
				}
			},
			window,
		)
	})
	// outputDataFolderLabel := widget.NewLabel("")

	form := widget.NewForm(
		widget.NewFormItem("Research Name", researchName),
		widget.NewFormItem("Current Date", currentDateTime),
		widget.NewFormItem("User Name", userName),
		widget.NewFormItem("Output Data Folder", outputDataFolder),
	)

	form.OnSubmit = func() {
		switch {
		case len(researchName.Text) <= 0:
			researchName.SetText("Research Name Is Empty.")
		case len(currentDateTime.Text) <= 0:
			currentDateTime.SetText("Date Is Empty.")
		case len(userName.Text) <= 0:
			userName.SetText("User Name Is Empty.")
		case len(outputDataFolderPath) <= 0:
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			outputDataFolder.SetIcon(buttonIcon)
		default:
			log.Println(trimm.BuildMainCommand())
			window.Close()
		}
		log.Println(researchName.Text, currentDateTime.Text, userName.Text)
	}

	form.OnCancel = func() {
		newApp.App.Quit()
	}

	return form
}
