package app

import (
	"biolink-gui/internal/models"
	"biolink-gui/internal/trimmomatic"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
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

// func Run() {
// 	newApp := NewAppStruct()
// 	windowName := newApp.App.NewWindow("App")
// 	windowName.Resize(fyne.NewSize(400, 320))

// 	icon, err := fyne.LoadResourceFromPath("images/icon.png")
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	windowName.SetIcon(icon)

// 	label := widget.NewLabel("Hello from FYNE")

// 	trimm := trimmomatic.Trimmomatic{
// 		Params: &models.TrimmomaticParams{},
// 	}

// 	// btn, for, rev :=
// 	a, b, c := trimm.ReadParams(windowName)

// 	windowName.SetContent(container.NewVBox(
// 		label,
// 		trimm.Path,
// 		// trimm.Input,
// 		trimm.Output,
// 		trimm.Paired,
// 		trimm.Phred,
// 		trimm.Threads,
// 		a, b, c,
// 	))

// 	// log.Println(trimm.Params.Path)

// 	windowName.ShowAndRun()
// }

func Run() {
	newApp := NewAppStruct()

	newApp.BuildMainWindow()

	newApp.App.Run()
}

func (newApp *App) BuildMainWindow() {
	mainWindow := newApp.App.NewWindow("App")
	mainWindow.Resize(fyne.NewSize(600, 420))

	icon, err := fyne.LoadResourceFromPath("images/icon.png")
	if err != nil {
		log.Println(err)
	}

	mainWindow.SetIcon(icon)

	trimm := trimmomatic.Trimmomatic{
		Params: &models.TrimmomaticParams{
			Prefix: "java -jar",
			Path:   "internal/tools/Trimmomatic-0.35/trimmomatic-0.35.jar",
		},
	}

	var subWindow fyne.Window

	chosePairedReadsBotton := widget.NewButton("Paired", func() {
		subWindow = newApp.App.NewWindow("Choose paired reads")
		subWindow.Resize(fyne.NewSize(500, 100))

		fov, rev, reads := trimm.SelectPairedReadsFiles(subWindow)
		subWindow.SetContent(container.NewVBox(
			reads,
			fov, rev,
		))
		subWindow.Show()
	})

	choseSingleReadsBotton := widget.NewButton("Single", func() {
		subWindow = newApp.App.NewWindow("Choose paired reads")
		subWindow.Resize(fyne.NewSize(500, 100))

		selected, frm := trimm.SelectSingleReadsFiles(subWindow)
		subWindow.SetContent(container.NewVBox(
			frm,
			selected,
		))
		subWindow.Show()
	})

	trimCont := CreateTrimmomaticContainer(chosePairedReadsBotton, choseSingleReadsBotton)
	anotherCont := CreateAnotherContainer(chosePairedReadsBotton, choseSingleReadsBotton)
	tabs := container.NewAppTabs(
		container.NewTabItem("Trimmomatic", trimCont),
		container.NewTabItem("Another Trimmomatic", anotherCont),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	newForm := CreateNewResearchForm(newApp, mainWindow)

	// a, err := trimm.ReadParams()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	mainWindow.SetContent(container.NewVBox(
		// label,
		// a,
		tabs,
		newForm,
	))
	mainWindow.Show()
	mainWindow.SetMaster()
}

func CreateTrimmomaticContainer(paired, single *widget.Button) *fyne.Container {
	container := container.NewVBox(
		widget.NewLabel("Trimmomatic"),
		paired,
		single,
	)

	return container
}

func CreateNewResearchForm(newApp *App, window fyne.Window) *widget.Form {
	userName := widget.NewEntry()
	// userNameLabel := widget.NewLabel("")

	currentTime := time.Now().Format(time.DateTime)
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
			// default:
			// 	window.Close()
		}
		log.Println(researchName.Text, currentDateTime.Text, userName.Text)
	}

	form.OnCancel = func() {
		newApp.App.Quit()
	}

	return form
}

func CreateAnotherContainer(paired, single *widget.Button) *fyne.Container {
	container := container.NewVBox(
		widget.NewLabel("Another trimmomatic"),
		// paired,
		// single,
	)

	return container
}
