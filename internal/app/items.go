package app

import (
	"biolink-nipt-gui/internal/app/fastqc"
	"biolink-nipt-gui/internal/trimmomatic"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

func CreateAnalysisItems(window fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, exitTerminal chan bool, w, h float32) *fyne.Menu {
	var subWindow fyne.Window

	trymmomaticTool := CreateTrimmomaticAnalysisItems(subWindow, trimm, newApp, commandChan, w, h)

	fastQCConfig := fastqc.New()
	fastQCConfig.Icons = trimm.Icons
	fastQCItems := CreateFastQCAnalysisItems(window, newApp, fastQCConfig, commandChan, w, h)

	analysisMenu := fyne.NewMenu("Analysis", trymmomaticTool, fastQCItems)

	return analysisMenu
}

func CreateTrimmomaticAnalysisItems(subWindow fyne.Window, trimm *trimmomatic.Trimmomatic, newApp *App, commandChan chan string, w, h float32) *fyne.MenuItem {
	pairedReads := fyne.NewMenuItem("Paired reads", func() {
		subWindow = newApp.App.NewWindow("Choose paired reads")
		subWindow.Resize(fyne.NewSize(w, h))

		fov, rev, description, reads := trimm.SelectPairedReadsFiles(subWindow, commandChan)
		subWindow.SetContent(container.NewVBox(
			reads,
			fov, rev,
			description,
		))
		subWindow.CenterOnScreen()
		subWindow.Show()
	})

	singleReads := fyne.NewMenuItem("Single reads", func() {
		subWindow = newApp.App.NewWindow("Choose single reads")
		subWindow.Resize(fyne.NewSize(w, h))

		selected, description, frm := trimm.SelectSingleReadsFiles(subWindow, commandChan)
		subWindow.SetContent(container.NewVBox(
			frm,
			selected,
			description,
		))
		// rect := canvas.NewRectangle(color.White)
		// subWindow.SetContent(rect)
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

func CreateFastQCAnalysisItems(subWindow fyne.Window, newApp *App, fastQC *fastqc.FastQC, commandChan chan string, w, h float32) *fyne.MenuItem {
	FastQCItem := fyne.NewMenuItem("FastQC", func() {
		subWindow = newApp.App.NewWindow("FastQC")
		subWindow.Resize(fyne.NewSize(w, h))

		fastQCForm, description := fastQC.CreateFastQCForn(subWindow, commandChan)
		subWindow.SetContent(container.NewVBox(
			fastQCForm,
			description,
		))
		subWindow.CenterOnScreen()
		subWindow.Show()
	})

	fastQCV_0_11_9_Tool := fyne.NewMenuItem(fmt.Sprintf("%v %v", "FastQC", fastQC.Version.Value), nil)
	fastQCV_0_11_9_Tool.ChildMenu = fyne.NewMenu("", FastQCItem)

	FastQCTool := fyne.NewMenuItem("FastQC", nil)
	FastQCTool.ChildMenu = fyne.NewMenu("", fastQCV_0_11_9_Tool)

	return FastQCTool
}

func CreateTimeItems(newApp *App) *fyne.Menu {
	var timeWindow fyne.Window

	createNewClock := fyne.NewMenuItem("Watch current time", func() {
		clock := widget.NewLabel("")
		clock.TextStyle = fyne.TextStyle{
			Monospace: true,
			TabWidth:  20,
		}
		updateTime(clock)

		go func() {
			for range time.Tick(time.Second) {
				updateTime(clock)
			}
		}()

		timeWindow = newApp.App.NewWindow("Time")
		timeWindow.SetContent(clock)
		timeWindow.Show()
	})

	clockMenu := fyne.NewMenu("Time", createNewClock)

	return clockMenu
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}

func CreateTerminalItems(window fyne.Window, newApp *App, commandChan chan string, exitTerminal chan bool, X, Y float32) *fyne.Menu {
	terminal := CreateTerminalWindow(newApp, commandChan, exitTerminal, 0, 150, window)
	terminal.Resize(fyne.NewSize(X, Y))
	createNewTerminal := fyne.NewMenuItem("Create", func() {
		terminal.Show()
	})

	selTerminalParams := fyne.NewMenuItem("Create (TODO)", func() {
		terminal.Show()
	})

	terminalMenuItem := fyne.NewMenuItem("Terminal", func() {})
	terminalMenuItem.ChildMenu = fyne.NewMenu("Terminal (TODO)", createNewTerminal, selTerminalParams)

	terminalMenu := fyne.NewMenu("Terminal", terminalMenuItem)

	return terminalMenu
}
