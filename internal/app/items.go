package app

import (
	"biolink-nipt-gui/internal/trimmomatic"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

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

	trymmomaticTool := CreateTrimmomaticAnalysisItems(subWindow, trimm, newApp, commandChan, exitTerminal)

	analysisMenu := fyne.NewMenu("Analysis", trymmomaticTool)

	return analysisMenu
}

func CreateTrimmomaticAnalysisItems(subWindow fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, exitTerminal chan bool) *fyne.MenuItem {
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

	trymmomatic035Tool := fyne.NewMenuItem("Trimmomatic-0.35", nil)
	trymmomatic035Tool.ChildMenu = fyne.NewMenu("", pairedReads, singleReads)

	trymmomatic039Tool := fyne.NewMenuItem("Trimmomatic-0.39", nil)
	trymmomatic039Tool.ChildMenu = fyne.NewMenu("", pairedReads)

	trymmomaticTool := fyne.NewMenuItem("Trimmomatic", nil)
	trymmomaticTool.ChildMenu = fyne.NewMenu("", trymmomatic035Tool, trymmomatic039Tool)

	return trymmomaticTool
}

func CreateTerminalItems(window fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, exitTerminal chan bool, X, Y float32) *fyne.Menu {
	createNewTerminal := fyne.NewMenuItem("Create", func() {
		terminal := CreateTerminalWindow(newApp, commandChan, exitTerminal, 0, 150, window)
		terminal.Resize(fyne.NewSize(X, Y))
		terminal.Show()
	})

	selTerminalParams := fyne.NewMenuItem("Create (TODO)", func() {
		terminal := CreateTerminalWindow(newApp, commandChan, exitTerminal, 0, 150, window)
		terminal.Resize(fyne.NewSize(X, Y))
	})

	terminalMenuItem := fyne.NewMenuItem("Terminal", func() {})
	terminalMenuItem.ChildMenu = fyne.NewMenu("Terminal (TODO)", createNewTerminal, selTerminalParams)

	terminalMenu := fyne.NewMenu("Terminal", terminalMenuItem)

	return terminalMenu
}
