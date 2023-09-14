package app

import (
	"biolink-nipt-gui/internal/trimmomatic"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

	terminalMenu := CreateTerminalItems(window, trimm, newApp, commandChan, exitTerminal, float32(displays[0]["X"])*0.9, float32(displays[0]["Y"])*0.2)

	menu := fyne.NewMainMenu(fileMenu, analysisMenu, terminalMenu)
	return menu
}

func CreateFileItems(window fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App) *fyne.Menu {
	var subWindow fyne.Window

	openResearch := fyne.NewMenuItem("Open Research", func() {
		dialog.ShowFolderOpen(
			func(r fyne.ListableURI, err error) {
				if r != nil {
					log.Println("The research path is", r.Path())
					// forwardSelected.SetText(r.URI().Path())
				}
			},
			window,
		)
	})

	openProtocol := fyne.NewMenuItem("Open Protocol", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					log.Println("The protocol path is", r.URI().Path())
					// forwardSelected.SetText(r.URI().Path())
				}
			},
			window,
		)
	})
	open := fyne.NewMenuItem("Open", nil)
	open.ChildMenu = fyne.NewMenu("Open", openResearch, openProtocol)

	newResearch := fyne.NewMenuItem("New Research", func() {
		subWindow = newApp.App.NewWindow("Choose paired reads")
		subWindow.Resize(fyne.NewSize(500, 100))

		newResearchForm := CreateNewResearchForm(newApp, subWindow, trimm)

		subWindow.SetContent(container.NewVBox(
			newResearchForm,
		))
		subWindow.CenterOnScreen()
		subWindow.Show()
	})

	newProtocol := fyne.NewMenuItem("New Protocol", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					log.Println("The protocol path is", r.URI().Path())
					// forwardSelected.SetText(r.URI().Path())
				}
			},
			window,
		)
	})
	newItem := fyne.NewMenuItem("New", nil)
	newItem.ChildMenu = fyne.NewMenu("New", newResearch, newProtocol)

	fileMenu := fyne.NewMenu("File", open, newItem)

	return fileMenu
}

func CreateAnalysisItems(window fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, exitTerminal chan bool) *fyne.Menu {
	var subWindow fyne.Window

	pairedReads := fyne.NewMenuItem("Paired reads", func() {
		subWindow = newApp.App.NewWindow("Choose paired reads")
		subWindow.Resize(fyne.NewSize(500, 300))

		fov, rev, reads := trimm.SelectPairedReadsFiles(subWindow, commandChan, exitTerminal)
		subWindow.SetContent(container.NewVBox(
			reads,
			fov, rev,
		))
		subWindow.CenterOnScreen()
		subWindow.Show()
	})

	singleReads := fyne.NewMenuItem("Single reads", func() {
		subWindow = newApp.App.NewWindow("Choose single reads")
		subWindow.Resize(fyne.NewSize(500, 300))

		selected, frm := trimm.SelectSingleReadsFiles(subWindow, commandChan, exitTerminal)
		subWindow.SetContent(container.NewVBox(
			frm,
			selected,
		))
		subWindow.CenterOnScreen()
		subWindow.Show()
	})

	trymmomaticTool := fyne.NewMenuItem("Trimmomatic", nil)
	trymmomaticTool.ChildMenu = fyne.NewMenu("Trimmomatic-0.35", pairedReads, singleReads)

	analysisMenu := fyne.NewMenu("Analysis", trymmomaticTool)

	return analysisMenu
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

func CreateTerminalItems(window fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, exitTerminal chan bool, X, Y float32) *fyne.Menu {
	createNewTerminal := fyne.NewMenuItem("Create", func() {
		terminal := CreateNewTerminalWindow(newApp, commandChan, exitTerminal)
		terminal.Resize(fyne.NewSize(X, Y))
		terminal.Show()
	})

	selTerminalParams := fyne.NewMenuItem("Create (TODO)", func() {
		terminal := CreateNewTerminalWindow(newApp, commandChan, exitTerminal)
		terminal.Resize(fyne.NewSize(X, Y))
	})

	terminalMenuItem := fyne.NewMenuItem("Terminal", func() {})
	terminalMenuItem.ChildMenu = fyne.NewMenu("Terminal (TODO)", createNewTerminal, selTerminalParams)

	terminalMenu := fyne.NewMenu("Terminal", terminalMenuItem)

	return terminalMenu
}

func CreateNewTerminalWindow(newApp *App, commandsChan chan string, exitRutine chan bool) fyne.Window {
	terminal := newApp.App.NewWindow("Terminal")
	// terminal.
	commandLabel := widget.NewMultiLineEntry()

	go func() {
		for {
			select {
			case command := <-commandsChan:
				if len(command) > 0 {
					commandLabel.Text = fmt.Sprintf("%v> %v\n", commandLabel.Text, command)
					terminal.SetContent(commandLabel)
				}
			// case <-reverseChan:
			// 	buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	reverseButton.SetIcon(buttonIcon)
			case <-exitRutine:
				close(commandsChan)
				close(exitRutine)
				return
			}
		}
	}()

	return terminal
}
