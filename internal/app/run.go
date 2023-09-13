package app

import (
	"biolink-nipt-gui/internal/models"
	"biolink-nipt-gui/internal/pkg"
	"biolink-nipt-gui/internal/trimmomatic"
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
	displays := pkg.GetDisplayParams()
	mainWindow.Resize(fyne.NewSize(float32(displays[0]["X"])*0.9, float32(displays[0]["Y"])*0.9))

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

	// trimCont := CreateTrimmomaticContainer(chosePairedReadsBotton, choseSingleReadsBotton)
	// anotherCont := CreateAnotherContainer(chosePairedReadsBotton, choseSingleReadsBotton)
	// tabs := container.NewAppTabs(
	// 	container.NewTabItem("Trimmomatic", trimCont),
	// 	container.NewTabItem("Another Trimmomatic", anotherCont),
	// )
	// tabs.SetTabLocation(container.TabLocationLeading)

	menu := CreateMainMenu(mainWindow, &trimm, newApp)

	mainWindow.SetMainMenu(menu)
	mainWindow.Show()
	mainWindow.SetMaster()
	newApp.App.Run()
}

func CreateTrimmomaticContainer(paired, single *widget.Button) *fyne.Container {
	container := container.NewVBox(
		widget.NewLabel("Trimmomatic"),
		paired,
		single,
	)

	return container
}

func CreateAnotherContainer(paired, single *widget.Button) *fyne.Container {
	container := container.NewVBox(
		widget.NewLabel("Another trimmomatic"),
		// paired,
		// single,
	)

	return container
}
