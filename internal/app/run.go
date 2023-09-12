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

	// label := widget.NewLabel("Hello from FYNE")

	trimm := trimmomatic.Trimmomatic{
		Params: &models.TrimmomaticParams{},
	}

	subWindow := newApp.App.NewWindow("Choose paired reads")
	subWindow.Resize(fyne.NewSize(500, 500))

	chosePairedReadsBotton := widget.NewButton("Paired", func() {
		forw, rev := trimm.SelectPairedReadsFiles(subWindow)
		subWindow.SetContent(container.NewVBox(
			forw, rev,
		))
		subWindow.Show()
	})

	choseSingleReadsBotton := widget.NewButton("Single", func() {
		single := trimm.SelectSingleReadsFiles(subWindow)
		subWindow.SetContent(container.NewVBox(
			single,
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

	newForm := CreateNewResearchForm(newApp)

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

func CreateNewResearchForm(newApp *App) *widget.Form {
	userName := widget.NewEntry()
	currentTime := time.Now().Format("2017-09-07 17:06:06")
	currentDateTime := widget.NewEntryWithData(binding.BindString(&currentTime))
	researchName := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("Research Name", researchName),
		widget.NewFormItem("Current Date", currentDateTime),
		widget.NewFormItem("User Name", userName),
	)

	form.OnSubmit = func() {
		log.Println("Submited")
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
