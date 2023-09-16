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

	icons := new(models.Icons)
	submitIcon, err := fyne.LoadResourceFromPath("images/accept.png")
	if err != nil {
		log.Println(err)
	}
	warningIcon, err := fyne.LoadResourceFromPath("images/warning.png")
	if err != nil {
		log.Println(err)
	}
	questionIcon, err := fyne.LoadResourceFromPath("images/question.svg")
	if err != nil {
		log.Println(err)
	}
	icons.Submit = submitIcon
	icons.Warning = warningIcon
	icons.Question = questionIcon

	trimm := trimmomatic.Trimmomatic{
		Params: &models.TrimmomaticParams{
			Prefix:      "java -jar",
			Path:        "/home/mrred/'Рабочий стол'/Работа/nipt-gui/internal/tools/Trimmomatic-0.35/trimmomatic-0.35.jar",
			Description: models.NewMainParamsDescription(),
			SubParams: &models.TrimmomaticSubParams{
				Description:   models.NewDescription(),
				Names:         models.NewSubparamsNames(),
				QuestionWidth: 50,
			},
		},
		Icons: icons,
	}

	commandsChan := make(chan string)
	exitRootine := make(chan bool, 1)
	terminalCloseChan := make(chan bool, 1)
	defer close(commandsChan)
	defer close(exitRootine)
	defer close(terminalCloseChan)

	menu := CreateMainMenu(mainWindow, &trimm, newApp, commandsChan, exitRootine, terminalCloseChan, displays)

	mainWindow.SetMainMenu(menu)
	// mainWindow.SetContent(cnt)
	mainWindow.Show()
	mainWindow.SetMaster()
	newApp.App.Run()
	terminalCloseChan <- true
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
