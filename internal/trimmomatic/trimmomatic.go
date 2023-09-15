package trimmomatic

import (
	"biolink-nipt-gui/internal/models"
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/cpu"
)

type Trimmomatic struct {
	Params  *models.TrimmomaticParams
	Path    *widget.Entry
	Input   *widget.Entry
	Output  *widget.Entry
	Phred   *widget.Entry
	Threads *widget.Entry
	Paired  *widget.RadioGroup
}

func (trimm *Trimmomatic) CreateTrimmomaticWindow() fyne.Window {

	return nil
}

func (trimm *Trimmomatic) SelectPairedReadsFiles(window fyne.Window, commandChan chan string, exitTerminat chan bool) (*widget.Label, *widget.Label, *widget.Form) {
	trimm.Params.Input = ""
	forwardInput, reverseInput := "", ""
	forwardSelected := widget.NewLabel("")
	reverseSelected := widget.NewLabel("")

	forwardChan := make(chan bool, 1)
	reverseChan := make(chan bool, 1)
	outputDirChan := make(chan bool, 1)
	exitRutine := make(chan bool, 1)

	forwardButton := widget.NewButton("Forward Reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					forwardInput += r.URI().Path()
					forwardSelected.SetText(r.URI().Path())
					if len(forwardInput) > 0 {
						forwardChan <- true
					} else {
						forwardChan <- false
					}
				}
			},
			window,
		)
	})

	reverseButton := widget.NewButton("Reverse Reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					reverseInput += r.URI().Path()
					reverseSelected.SetText(r.URI().Path())
					if len(reverseInput) > 0 {
						reverseChan <- true
					} else {
						reverseChan <- false
					}
				}
			},
			window,
		)
	})

	outputDir := trimm.ChoseOutputDir(window, outputDirChan)
	outputDirItem := widget.NewFormItem("Output directory", outputDir)
	outputFormatItem := trimm.ChoseOutputFormat()
	phredItems := trimm.ChosePhred()
	threadsEntry := trimm.ChoseThreads()
	threadsEntryItem := widget.NewFormItem("Threads", threadsEntry)
	logFileItem := trimm.SaveLogs()

	illuminaClip, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold, isClip := trimm.CreateIlluminaClip()
	slidingWindow, windowSize, requiredQuality, isSliding := trimm.CreateSlidingWindow()
	leadingItem, leading, isLeading := trimm.CreateLeading()
	trailingItem, trailing, isTrailing := trimm.CreateTrailing()
	cropItem, crop, isCrop := trimm.CreateCrop()
	headCropItem, headCrop, isHeadCrop := trimm.CreateHeadCrop()
	minLenItem, minLen, isMinLen := trimm.CreateMinLen()

	go func() {
		for {
			select {
			case <-forwardChan:
				buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
				if err != nil {
					log.Println(err)
				}
				forwardButton.SetIcon(buttonIcon)
			case <-reverseChan:
				buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
				if err != nil {
					log.Println(err)
				}
				reverseButton.SetIcon(buttonIcon)
			case <-outputDirChan:
				buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
				if err != nil {
					log.Println(err)
				}
				outputDir.SetIcon(buttonIcon)
			case <-exitRutine:
				trimm.Params.Input = fmt.Sprintf("%v %v", forwardInput, reverseInput)
				close(forwardChan)
				close(reverseChan)
				close(exitRutine)
				return
			}
		}
	}()

	choseReadsForm := widget.NewForm(
		widget.NewFormItem("Select input forward file", forwardButton),
		widget.NewFormItem("Select input reverse file", reverseButton),
		outputDirItem,
		outputFormatItem,
		phredItems,
		threadsEntryItem,
		logFileItem,
		illuminaClip,
		slidingWindow,
		leadingItem,
		trailingItem,
		cropItem,
		headCropItem,
		minLenItem,
	)

	choseReadsForm.OnSubmit = func() {
		switch {
		case len(forwardInput) <= 0:
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			forwardButton.SetIcon(buttonIcon)
			forwardSelected.SetText("Field is empty. Please, chose input data file.")
		case len(reverseInput) <= 0:
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			reverseButton.SetIcon(buttonIcon)
			reverseSelected.SetText("Field is empty. Please, chose input data file.")
		case len(trimm.Params.Output) <= 0:
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			outputDir.SetIcon(buttonIcon)
		case isClip:
			// Если тут не заработает, переставить в default
			trimm.Params.SubParams.IlluminaClip = fmt.Sprintf("ILLUMINACLIP:%v:%v:%v:%v", fastaWithAdaptersEtc.Text, seedMismatches.Text, palindromeClipThreshold.Text, simpleClipThreshold.Text)
		case isSliding:
			trimm.Params.SubParams.SlidingWindow = fmt.Sprintf("SLIDINGWINDOW:%v:%v", windowSize.Text, requiredQuality.Text)
		case isLeading:
			trimm.Params.SubParams.Leading = fmt.Sprintf("LEADING:%v", leading.Text)
		case isTrailing:
			trimm.Params.SubParams.Trailing = fmt.Sprintf("TRAILING:%v", trailing.Text)
		case isCrop:
			trimm.Params.SubParams.Crop = fmt.Sprintf("CROP:%v", crop.Text)
		case isHeadCrop:
			trimm.Params.SubParams.HeadCrop = fmt.Sprintf("HEADCROP:%v", headCrop.Text)
		case isMinLen:
			trimm.Params.SubParams.MinLen = fmt.Sprintf("MINLEN:%v", minLen.Text)
		default:
			exitRutine <- true
			trimm.Params.Input = fmt.Sprintf("%v %v", forwardInput, reverseInput)
			trimm.Params.Paired = "PE"

			thrds, err := strconv.Atoi(threadsEntry.Text)
			if err != nil {
				log.Println(err)
			}

			trimm.Params.Threads = thrds
			window.Close()
			cmnd, err := trimm.BuildMainCommand()
			if err != nil {
				log.Println(err)
			}
			commandChan <- cmnd
		}
	}

	// choseReadsForm.OnCancel = func() {
	// 	window.Close()
	// }

	return forwardSelected, reverseSelected, choseReadsForm
}

func (trimm *Trimmomatic) SelectSingleReadsFiles(window fyne.Window, commandChan chan string, exitTerminat chan bool) (*widget.Label, *widget.Form) {
	trimm.Params.Input = ""
	selected := widget.NewLabel("")

	readsChan := make(chan bool, 1)
	outputDirChan := make(chan bool, 1)
	exitRutine := make(chan bool, 1)

	button := widget.NewButton("Reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					trimm.Params.Input += r.URI().Path()
					selected.SetText(r.URI().Path())
					if len(selected.Text) > 0 {
						readsChan <- true
					} else {
						readsChan <- false
					}
				}
			},
			window,
		)
	})

	outputDir := trimm.ChoseOutputDir(window, outputDirChan)
	outputDirItem := widget.NewFormItem("Output directory", outputDir)
	outputFormatItem := trimm.ChoseOutputFormat()
	phredItems := trimm.ChosePhred()
	threadsEntry := trimm.ChoseThreads()
	threadsEntryItem := widget.NewFormItem("Threads", threadsEntry)
	logFileItem := trimm.SaveLogs()

	illuminaClip, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold, isClip := trimm.CreateIlluminaClip()
	slidingWindow, windowSize, requiredQuality, isSliding := trimm.CreateSlidingWindow()
	leadingItem, leading, isLeading := trimm.CreateLeading()
	trailingItem, trailing, isTrailing := trimm.CreateTrailing()
	cropItem, crop, isCrop := trimm.CreateCrop()
	headCropItem, headCrop, isHeadCrop := trimm.CreateHeadCrop()
	minLenItem, minLen, isMinLen := trimm.CreateMinLen()

	go func() {
		for {
			select {
			case <-readsChan:
				buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
				if err != nil {
					log.Println(err)
				}
				button.SetIcon(buttonIcon)
			case <-outputDirChan:
				buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
				if err != nil {
					log.Println(err)
				}
				outputDir.SetIcon(buttonIcon)
			case <-exitRutine:
				close(readsChan)
				close(exitRutine)
				return
			}
		}
	}()

	choseReadsForm := widget.NewForm(
		widget.NewFormItem("Select input file", button),
		outputDirItem,
		outputFormatItem,
		phredItems,
		threadsEntryItem,
		logFileItem,
		illuminaClip,
		slidingWindow,
		leadingItem,
		trailingItem,
		cropItem,
		headCropItem,
		minLenItem,
	)

	choseReadsForm.OnSubmit = func() {
		switch {
		case len(trimm.Params.Input) <= 0:
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			button.SetIcon(buttonIcon)
			selected.SetText("Field is empty. Please, chose input data file.")
		case len(trimm.Params.Output) <= 0:
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			outputDir.SetIcon(buttonIcon)
			// selected.SetText("Field is empty. Please, chose input data file.")
		case isClip:
			// Если тут не заработает, переставить в default
			trimm.Params.SubParams.IlluminaClip = fmt.Sprintf("ILLUMINACLIP:%v:%v:%v:%v", fastaWithAdaptersEtc.Text, seedMismatches.Text, palindromeClipThreshold.Text, simpleClipThreshold.Text)
		case isSliding:
			trimm.Params.SubParams.SlidingWindow = fmt.Sprintf("SLIDINGWINDOW:%v:%v", windowSize.Text, requiredQuality.Text)
		case isLeading:
			trimm.Params.SubParams.Leading = fmt.Sprintf("LEADING:%v", leading.Text)
		case isTrailing:
			trimm.Params.SubParams.Trailing = fmt.Sprintf("TRAILING:%v", trailing.Text)
		case isCrop:
			trimm.Params.SubParams.Crop = fmt.Sprintf("CROP:%v", crop.Text)
		case isHeadCrop:
			trimm.Params.SubParams.HeadCrop = fmt.Sprintf("HEADCROP:%v", headCrop.Text)
		case isMinLen:
			trimm.Params.SubParams.MinLen = fmt.Sprintf("MINLEN:%v", minLen.Text)
		default:
			exitRutine <- true
			window.Close()
			trimm.Params.Paired = "SE"
			thrds, err := strconv.Atoi(threadsEntry.Text)
			if err != nil {
				log.Println(err)
			}

			trimm.Params.Threads = thrds
			cmnd, err := trimm.BuildMainCommand()
			if err != nil {
				log.Println(err)
			}
			commandChan <- cmnd
		}
	}

	return selected, choseReadsForm
}

func (trimm *Trimmomatic) ChosePhred() *widget.FormItem {
	phredGroup := widget.NewRadioGroup([]string{"33", "64"}, func(s string) {
		phred, err := strconv.Atoi(s)
		if err != nil {
			log.Println(err)
		}
		trimm.Params.Phred = phred
	})
	phredGroup.Horizontal = true
	phredItem := widget.NewFormItem("Phred", phredGroup)
	return phredItem
}

func (trimm *Trimmomatic) ChoseThreads() *widget.SelectEntry {
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Println(err)
	}

	cpuSlice := []string{}
	for i := 1; i <= len(cpuInfo); i++ {
		cpuSlice = append(cpuSlice, strconv.Itoa(i))
	}

	selectThreadsCount := widget.NewSelectEntry(cpuSlice)
	selectThreadsCount.SetText(strconv.Itoa(len(cpuInfo)))

	return selectThreadsCount
}

func (trimm *Trimmomatic) SaveLogs() *widget.FormItem {
	logFile := widget.NewCheck("logfile.log", func(b bool) {
		if b {
			trimm.Params.Logfile = trimm.Params.Output + "/logfile.log"
			log.Println(trimm.Params.Logfile)
		}
	})
	logFileItem := widget.NewFormItem("Save logs", logFile)
	return logFileItem
}

var outputFormatMap = map[string]bool{
	"Base":     true,
	"Unpaired": false,
}

func (trimm *Trimmomatic) ChoseOutputFormat() *widget.FormItem {
	outputGroup := widget.NewRadioGroup([]string{"Base", "Paired+Unpaired"}, func(s string) {
		trimm.Params.BaseOutput = outputFormatMap[s]
	})
	outputGroup.Horizontal = true
	outputGroupItem := widget.NewFormItem("Output format", outputGroup)
	return outputGroupItem
}

func (trimm *Trimmomatic) ChoseOutputDir(window fyne.Window, dirChan chan bool) *widget.Button {
	openResearch := widget.NewButton("Change output dir", func() {
		dialog.ShowFolderOpen(
			func(r fyne.ListableURI, err error) {
				if r != nil {
					trimm.Params.Output = r.Path()
					if len(trimm.Params.Output) > 0 {
						dirChan <- true
					} else {
						dirChan <- false
					}
				}
			},
			window,
		)
	})
	return openResearch
}

func (trimm *Trimmomatic) CreateIlluminaClip() (*widget.FormItem, *widget.SelectEntry, *widget.Entry, *widget.Entry, *widget.Entry, bool) {
	fastaWithAdaptersEtc := widget.NewSelectEntry([]string{"TruSeq3-SE"})
	seedMismatches := widget.NewEntry()
	palindromeClipThreshold := widget.NewEntry()
	simpleClipThreshold := widget.NewEntry()
	isClip := false
	isIlluminaClip := widget.NewCheck("Use", func(b bool) {
		isClip = true
	})
	cols := container.NewGridWithColumns(5, isIlluminaClip, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold)
	// rows := container.NewGridWithRows(2, cols, isIlluminaClip)
	formItem := widget.NewFormItem("ILLUMINACLIP", cols)
	return formItem, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold, isClip
}

func (trimm *Trimmomatic) CreateSlidingWindow() (*widget.FormItem, *widget.Entry, *widget.Entry, bool) {
	windowSize := widget.NewEntry()
	requiredQuality := widget.NewEntry()
	isSliding := false
	isSlidingWindow := widget.NewCheck("Use", func(b bool) {
		isSliding = true
	})
	cols := container.NewGridWithColumns(3, isSlidingWindow, windowSize, requiredQuality)
	// rows := container.NewGridWithRows(2, cols, isSlidingWindow)
	formItem := widget.NewFormItem("SLIDINGWINDOW", cols)
	return formItem, windowSize, requiredQuality, isSliding
}

func (trimm *Trimmomatic) CreateLeading() (*widget.FormItem, *widget.Entry, bool) {
	quality := widget.NewEntry()
	isQuality := false
	isQualityCheck := widget.NewCheck("Use", func(b bool) {
		isQuality = true
	})
	cols := container.NewGridWithColumns(2, isQualityCheck, quality)
	formItem := widget.NewFormItem("LEADING", cols)
	return formItem, quality, isQuality
}

func (trimm *Trimmomatic) CreateTrailing() (*widget.FormItem, *widget.Entry, bool) {
	quality := widget.NewEntry()
	isQuality := false
	isQualityCheck := widget.NewCheck("Use", func(b bool) {
		isQuality = true
	})
	cols := container.NewGridWithColumns(2, isQualityCheck, quality)
	formItem := widget.NewFormItem("TRAILING", cols)
	return formItem, quality, isQuality
}

func (trimm *Trimmomatic) CreateCrop() (*widget.FormItem, *widget.Entry, bool) {
	length := widget.NewEntry()
	isLength := false
	isLengthCheck := widget.NewCheck("Use", func(b bool) {
		isLength = true
	})
	cols := container.NewGridWithColumns(2, isLengthCheck, length)
	formItem := widget.NewFormItem("CROP", cols)
	return formItem, length, isLength
}

func (trimm *Trimmomatic) CreateHeadCrop() (*widget.FormItem, *widget.Entry, bool) {
	length := widget.NewEntry()
	isLength := false
	isLengthCheck := widget.NewCheck("Use", func(b bool) {
		isLength = true
	})
	cols := container.NewGridWithColumns(2, isLengthCheck, length)
	formItem := widget.NewFormItem("HEADCROP", cols)
	return formItem, length, isLength
}

func (trimm *Trimmomatic) CreateMinLen() (*widget.FormItem, *widget.Entry, bool) {
	length := widget.NewEntry()
	isLength := false
	isLengthCheck := widget.NewCheck("Use", func(b bool) {
		isLength = true
	})
	cols := container.NewGridWithColumns(2, isLengthCheck, length)
	formItem := widget.NewFormItem("MINLEN", cols)
	return formItem, length, isLength
}
