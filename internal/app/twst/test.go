package main

import (
	"biolink-nipt-gui/internal/app/fastqc"
	"fmt"
	"log"
	"reflect"
)

// import (
// 	"biolink-nipt-gui/internal/pkg"
// 	protocol "biolink-nipt-gui/internal/protocols"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"path/filepath"
// 	"strconv"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	newApp := app.New()
// 	mainWindow := newApp.NewWindow("App")
// 	displays := pkg.GetDisplayParams()
// 	mainWindow.Resize(fyne.NewSize(float32(displays[0]["X"])*0.9, float32(displays[0]["Y"])*0.9))

// 	btn1 := widget.NewMultiLineEntry()
// 	btn1Grid := container.NewGridWrap(fyne.NewSize(800, 200), btn1)

// 	buttonIcon := widget.NewMultiLineEntry()
// 	buttonIconGrid := container.NewGridWrap(fyne.NewSize(800, 600), buttonIcon)

// 	topBox := container.NewGridWithRows(2, buttonIconGrid, btn1Grid)
// 	// topBox := container.NewGridWrap(fyne.NewSize(800, 800), buttonIconGrid, btn1Grid)

// 	// btn3 := widget.NewButton("Button 3", func() {})
// 	form := OpenProtocolForm(newApp)
// 	// formBox := container.NewVBox(form)

// 	cnt := container.NewGridWithColumns(2, topBox, form)
// 	// cnt := container.NewGridWrap(fyne.NewSize(float32(displays[0]["X"])*0.9, float32(displays[0]["Y"])*0.9), topBox, form)

// 	mainWindow.SetContent(cnt)
// 	mainWindow.Show()
// 	mainWindow.SetMaster()
// 	newApp.Run()
// }

// func CreateNewProtocolForm() *widget.Form {
// 	currentTime := time.Now().Format(time.RFC3339)
// 	currentDateTime := widget.NewEntryWithData(binding.BindString(&currentTime))

// 	name := widget.NewEntry()
// 	user := widget.NewSelectEntry([]string{"Pasha", "Dasha"})

// 	form := widget.NewForm(
// 		widget.NewFormItem("Protocol Name", name),
// 		widget.NewFormItem("Current Date", currentDateTime),
// 		widget.NewFormItem("User Name", user),
// 		// widget.NewFormItem("Output Data Folder", outputDataFolder),
// 	)

// 	form.OnSubmit = func() {
// 		log.Println(name.Text, user.Text)
// 	}

// 	return form
// }

// type Protocol struct {
// 	Version   *widget.FormItem
// 	Author    *widget.FormItem
// 	Stages    []*Stage
// 	CreatedAt *widget.FormItem
// 	UpdatedAt *widget.FormItem
// }

// type Stage struct {
// 	Number      *widget.FormItem
// 	Name        *widget.FormItem
// 	Params      *widget.FormItem
// 	Tool        *widget.FormItem
// 	Description *widget.FormItem
// 	ToolPath    *widget.FormItem
// }

// func OpenProtocolForm(newApp fyne.App) *widget.Form {
// 	prt := protocol.NewProtocol()
// 	filename, err := filepath.Abs("internal/protocols/test_protocol.yaml")
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	yamlFile, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	protocol, err := protocol.ParseProtocol(yamlFile, prt)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	newEntryProtocol := new(Protocol)

// 	version := widget.NewLabel(protocol.Version)
// 	newEntryProtocol.Version = widget.NewFormItem("Version", version)

// 	author := widget.NewLabel(protocol.Author)
// 	newEntryProtocol.Author = widget.NewFormItem("Author", author)

// 	createdAt := widget.NewLabel(protocol.CreatedAt.Format(time.RFC3339))
// 	newEntryProtocol.CreatedAt = widget.NewFormItem("Created at", createdAt)

// 	updatedAt := widget.NewLabel(protocol.UpdatedAt.Format(time.RFC3339))
// 	newEntryProtocol.UpdatedAt = widget.NewFormItem("Updated at", updatedAt)

// 	updateTitle := widget.NewButton("Change", func() {
// 		version := widget.NewEntry()
// 		version.SetText(protocol.Version)
// 		versionItem := widget.NewFormItem("Version", version)

// 		author := widget.NewEntry()
// 		author.SetText(protocol.Author)
// 		authorItem := widget.NewFormItem("Author", author)

// 		createdAt := widget.NewLabel(protocol.CreatedAt.Format(time.RFC3339))
// 		createdAtItem := widget.NewFormItem("Created at", createdAt)

// 		updatedAt := widget.NewLabel(protocol.UpdatedAt.Format(time.RFC3339))
// 		updatedAtItem := widget.NewFormItem("Updated at", updatedAt)

// 		updForm := widget.NewForm(
// 			versionItem,
// 			authorItem,
// 			createdAtItem,
// 			updatedAtItem,
// 		)

// 		// New Window
// 		changeTitleWindow := newApp.NewWindow("Change protocol")
// 		updForm.OnCancel = func() {
// 			changeTitleWindow.Close()
// 		}

// 		updForm.OnSubmit = func() {
// 			log.Println("SUBMITED")
// 			prt, _ := ChangeProtocolTitle(map[string]*widget.Entry{
// 				"Version": version,
// 				"Author":  author,
// 			})
// 			frm := widget.NewForm(prt.Version, prt.Author, createdAtItem, prt.UpdatedAt)
// 			subW := newApp.NewWindow("changed")
// 			subW.SetContent(frm)
// 			subW.Show()
// 			changeTitleWindow.Close()
// 		}
// 		changeTitleWindow.Resize(fyne.NewSize(300, 200))
// 		changeTitleWindow.SetContent(updForm)
// 		changeTitleWindow.Show()
// 	})
// 	updateTitle.Resize(fyne.NewSize(30, 10))
// 	updateItem := widget.NewFormItem("", updateTitle)

// 	newEntryProtocol.Stages = make([]*Stage, 0)

// 	for _, stage := range protocol.Stages {
// 		newStage := new(Stage)

// 		number := widget.NewEntry()
// 		strNum := strconv.Itoa(stage.Number)
// 		number.SetText(strNum)
// 		newStage.Number = widget.NewFormItem("Number", number)

// 		name := widget.NewEntry()
// 		name.SetText(stage.Name)
// 		newStage.Name = widget.NewFormItem("Name", name)

// 		params := widget.NewEntry()
// 		params.SetText(stage.Params)
// 		newStage.Params = widget.NewFormItem("Params", params)

// 		tool := widget.NewEntry()
// 		tool.SetText(stage.Tool)
// 		newStage.Tool = widget.NewFormItem("Tool", tool)

// 		description := widget.NewMultiLineEntry()
// 		description.SetText(stage.Description)
// 		newStage.Description = widget.NewFormItem("Description", description)

// 		toolPath := widget.NewEntry()
// 		toolPath.SetText(stage.ToolPath)
// 		newStage.ToolPath = widget.NewFormItem("ToolPath", toolPath)

// 		newEntryProtocol.Stages = append(newEntryProtocol.Stages, newStage)
// 	}

// 	// stages := widget.NewLabel("stages")
// 	itms := make([]*widget.FormItem, 0)

// 	a := 1

// 	for _, val := range newEntryProtocol.Stages {
// 		it1 := widget.NewFormItem(fmt.Sprintf("Stage №%v", a), widget.NewForm(
// 			val.Number,
// 			val.Name,
// 			val.Params,
// 			val.Tool,
// 			val.Description,
// 			val.ToolPath,
// 		))
// 		itms = append(itms, it1)

// 		a++
// 	}

// 	protocolForm := widget.NewForm(
// 		newEntryProtocol.Version,
// 		newEntryProtocol.Author,
// 		newEntryProtocol.CreatedAt,
// 		newEntryProtocol.UpdatedAt,
// 		updateItem,
// 		// stagesItem,
// 		itms[0],
// 		itms[1],
// 	)

// 	return protocolForm
// }

// func ChangeProtocolTitle(prt map[string]*widget.Entry) (*Protocol, error) {
// 	// for _, item := range prt.Items {
// 	// 	e := widget.NewEntry()
// 	// 	log.Println(item.Widget)
// 	// }

// 	changedProtocol := &Protocol{
// 		Version:   widget.NewFormItem("Version", widget.NewLabel(prt["Version"].Text)),
// 		Author:    widget.NewFormItem("Author", widget.NewLabel(prt["Author"].Text)),
// 		UpdatedAt: widget.NewFormItem("UpdatedAt", widget.NewLabel(time.Now().Format(time.RFC3339))),
// 	}

// 	return changedProtocol, nil
// }

// func main() {
// 	inputPath := "/home/mrred/Загрузки/ERR9792312.fastq.gz"

// 	file, err := os.Open(inputPath)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	reader := fastq.NewReader(file, linear.NewQSeq("", nil, alphabet.DNA, alphabet.Illumina1_8))
// 	r, err := reader.Read()
// 	log.Println(r, err)
// 	// go func() {
// 	// 	time.Sleep(50 * time.Millisecond)
// 	// 	file, err := os.Stat(path)

// 	// }()
// 	return
// }

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// func main() {
// 	file, err := ioutil.ReadFile("/home/mrred/Рабочий стол/Учёба/logfile.log")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	splt := strings.Split(string(file), "\n")
// 	a := 0
// 	for _, val := range splt {
// 		a++
// 		if len(val) == 0 {
// 			a--
// 		}
// 	}
// 	log.Println(a, splt[len(splt)-2])

// 	fp, r := pkg.Xopen("/home/mrred/Загрузки/ERR9792312.fastq.gz")
// 	defer fp.Close()

// 	n := 0
// 	var fqr pkg.FqReader
// 	qual, seq := []string{}, []string{}

// 	fqr.R = r
// 	for r, done := fqr.Iter(); !done; r, done = fqr.Iter() {
// 		n += 1
// 		seq = append(seq, string(r.Seq))
// 		qual = append(qual, string(r.Seq))
// 		// sLen += int64(len(r.Seq)
// 		// qLen += int64(len(r.Qual))
// 	}
// 	fmt.Println(n, "\t", len(seq), "\t", len(qual))
// 	log.Println(seq[10:20])
// }

func main() {
	fastQC := *fastqc.New()
	fastQC.Version.IsUsed = true
	v := reflect.ValueOf(fastQC)
	log.Println(v)

	// v := reflect.ValueOf(fastQC)
	// names := make([]string, 0, v.NumField())
	// v.FieldByNameFunc(func(fieldName string) bool {
	// 	names = append(names, fieldName)
	// 	return false
	// })
	// log.Println(names)

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			if v.Field(i).Elem().FieldByName("IsUsed").Interface().(bool) {
				log.Println(reflect.TypeOf(fastQC).Field(i).Name)
				log.Println(fmt.Sprintf("%v %v", v.Field(i).Elem().FieldByName("Flag").Interface(), v.Field(i).Elem().FieldByName("Value").Interface()))
			}
			// log.Println(v.Field(i).
		}
	}
	return
}
