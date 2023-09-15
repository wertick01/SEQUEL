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
			Path:   "/home/mrred/'Рабочий стол'/Работа/nipt-gui/internal/tools/Trimmomatic-0.35/trimmomatic-0.35.jar",
		},
	}

	commandsChan := make(chan string)
	closeTerminal := make(chan bool, 1)
	defer close(commandsChan)
	defer close(closeTerminal)

	menu := CreateMainMenu(mainWindow, &trimm, newApp, commandsChan, closeTerminal, displays)

	mainWindow.SetMainMenu(menu)
	// mainWindow.SetContent(cnt)
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
