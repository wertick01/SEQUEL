package main

import (
	"biolink-nipt-gui/internal/pkg"
	protocol "biolink-nipt-gui/internal/protocols"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	newApp := app.New()
	mainWindow := newApp.NewWindow("App")
	displays := pkg.GetDisplayParams()
	mainWindow.Resize(fyne.NewSize(float32(displays[0]["X"])*0.9, float32(displays[0]["Y"])*0.9))

	btn1 := widget.NewMultiLineEntry()
	// btn1 := widget.NewButton("Button 1", func() {})
	// btn1.Move(fyne.NewPos(0, 100))

	// btn2 := widget.NewButton("Button 2", func() {})
	// // btn2.Move(fyne.NewPos(0, 100))
	// btn2.Resize(fyne.NewSize(30, 10))

	buttonIcon := canvas.NewImageFromFile("images/warning.png")

	topBox := container.NewGridWithRows(2, btn1, buttonIcon)

	// btn3 := widget.NewButton("Button 3", func() {})
	form := OpenProtocolForm()
	// formBox := container.NewVBox(form)

	cnt := container.New(layout.NewGridLayout(2), topBox, form)
	// cnt := container.New

	mainWindow.SetContent(cnt)
	mainWindow.Show()
	mainWindow.SetMaster()
	newApp.Run()
}

func CreateNewProtocolForm() *widget.Form {
	currentTime := time.Now().Format(time.RFC3339)
	currentDateTime := widget.NewEntryWithData(binding.BindString(&currentTime))

	name := widget.NewEntry()
	user := widget.NewSelectEntry([]string{"Pasha", "Dasha"})

	form := widget.NewForm(
		widget.NewFormItem("Protocol Name", name),
		widget.NewFormItem("Current Date", currentDateTime),
		widget.NewFormItem("User Name", user),
		// widget.NewFormItem("Output Data Folder", outputDataFolder),
	)

	form.OnSubmit = func() {
		log.Println(name.Text, user.Text)
	}

	return form
}

type Protocol struct {
	Version   *widget.FormItem
	Author    *widget.FormItem
	Stages    []*Stage
	CreatedAt *widget.FormItem
	UpdatedAt *widget.FormItem
}

type Stage struct {
	Number      *widget.FormItem
	Name        *widget.FormItem
	Params      *widget.FormItem
	Tool        *widget.FormItem
	Description *widget.FormItem
	ToolPath    *widget.FormItem
}

func OpenProtocolForm() *widget.Form {
	prt := protocol.NewProtocol()
	filename, err := filepath.Abs("internal/protocols/test_protocol.yaml")
	if err != nil {
		log.Println(err)
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	protocol, err := protocol.ParseProtocol(yamlFile, prt)
	if err != nil {
		log.Println(err)
	}

	newEntryProtocol := new(Protocol)

	version := widget.NewEntry()
	version.SetText(protocol.Version)
	newEntryProtocol.Version = widget.NewFormItem("Version", version)

	author := widget.NewEntry()
	author.SetText(protocol.Author)
	newEntryProtocol.Author = widget.NewFormItem("Author", author)

	createdAt := widget.NewEntry()
	createdAt.SetText(protocol.CreatedAt.Format(time.RFC3339))
	newEntryProtocol.CreatedAt = widget.NewFormItem("Created at", createdAt)

	updatedAt := widget.NewEntry()
	updatedAt.SetText(protocol.UpdatedAt.Format(time.RFC3339))
	newEntryProtocol.UpdatedAt = widget.NewFormItem("Updated at", updatedAt)

	newEntryProtocol.Stages = make([]*Stage, 0)

	for _, stage := range protocol.Stages {
		newStage := new(Stage)

		number := widget.NewEntry()
		strNum := strconv.Itoa(stage.Number)
		number.SetText(strNum)
		newStage.Number = widget.NewFormItem("Number", number)

		name := widget.NewEntry()
		name.SetText(stage.Name)
		newStage.Name = widget.NewFormItem("Name", name)

		params := widget.NewEntry()
		params.SetText(stage.Params)
		newStage.Params = widget.NewFormItem("Params", params)

		tool := widget.NewEntry()
		tool.SetText(stage.Tool)
		newStage.Tool = widget.NewFormItem("Tool", tool)

		description := widget.NewEntry()
		description.SetText(stage.Description)
		newStage.Description = widget.NewFormItem("Description", description)

		toolPath := widget.NewEntry()
		toolPath.SetText(stage.ToolPath)
		newStage.ToolPath = widget.NewFormItem("ToolPath", toolPath)

		newEntryProtocol.Stages = append(newEntryProtocol.Stages, newStage)
	}

	// stages := widget.NewLabel("stages")
	stages := container.NewVBox()

	for _, val := range newEntryProtocol.Stages {
		stages.Add(widget.NewForm(
			val.Number,
			val.Name,
			val.Params,
			val.Tool,
			val.Description,
			val.ToolPath,
		))
	}
	// log.Println(version.Text, createdAt.Text)

	stagesItem := widget.NewFormItem("Stages", stages)

	protocolForm := widget.NewForm(
		newEntryProtocol.Version,
		newEntryProtocol.Author,
		newEntryProtocol.UpdatedAt,
		newEntryProtocol.CreatedAt,
		stagesItem,
	)

	return protocolForm
}
