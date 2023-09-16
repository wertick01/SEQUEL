package trimmomatic

import (
	"biolink-nipt-gui/internal/models"
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
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
	Icons   *models.Icons
}

func (trimm *Trimmomatic) CreateTrimmomaticWindow() fyne.Window {

	return nil
}

func (trimm *Trimmomatic) SelectPairedReadsFiles(window fyne.Window, commandChan chan string, exitTerminat chan bool) (*widget.Label, *widget.Label, *widget.Label, *widget.Form) {
	trimm.Params.Input = ""
	forwardInput, reverseInput := "", ""
	forwardSelected := widget.NewLabel("")
	reverseSelected := widget.NewLabel("")
	descriptionLabel := widget.NewLabel("")

	forwardChan := make(chan bool)
	reverseChan := make(chan bool)
	outputDirChan := make(chan bool)
	descriptionChan := make(chan string)
	exitRutine := make(chan bool, 1)

	isClipChan := make(chan bool, 1)
	isSlidingChan := make(chan bool, 1)
	isLeadingChan := make(chan bool, 1)
	isTrailingChan := make(chan bool, 1)
	isCropChan := make(chan bool, 1)
	isHeadCropChan := make(chan bool, 1)
	isMinLenChan := make(chan bool, 1)

	isClip := false
	isSliding := false
	isLeading := false
	isTrailing := false
	isCrop := false
	isHeadCrop := false
	isMinLen := false

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

	illuminaClip, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold := trimm.CreateIlluminaClip(descriptionChan, isClipChan)
	slidingWindow, windowSize, requiredQuality := trimm.CreateSlidingWindow(descriptionChan, isSlidingChan)
	leadingItem, leading := trimm.CreateLeading(descriptionChan, isLeadingChan)
	trailingItem, trailing := trimm.CreateTrailing(descriptionChan, isTrailingChan)
	cropItem, crop := trimm.CreateCrop(descriptionChan, isCropChan)
	headCropItem, headCrop := trimm.CreateHeadCrop(descriptionChan, isHeadCropChan)
	minLenItem, minLen := trimm.CreateMinLen(descriptionChan, isMinLenChan)

	go func() {
		for {
			select {
			case <-forwardChan:
				forwardButton.SetIcon(trimm.Icons.Submit)
			case <-reverseChan:
				reverseButton.SetIcon(trimm.Icons.Submit)
			case <-outputDirChan:
				outputDir.SetIcon(trimm.Icons.Submit)
			case description := <-descriptionChan:
				descriptionLabel.SetText(description)
			case clipVal := <-isClipChan:
				isClip = clipVal
			case slidingVal := <-isSlidingChan:
				isSliding = slidingVal
			case leadingVal := <-isLeadingChan:
				isLeading = leadingVal
			case trailingVal := <-isTrailingChan:
				isTrailing = trailingVal
			case cropVal := <-isCropChan:
				isCrop = cropVal
			case headCropVal := <-isHeadCropChan:
				isHeadCrop = headCropVal
			case minLenVal := <-isMinLenChan:
				isMinLen = minLenVal
			case <-exitRutine:
				trimm.Params.Input = fmt.Sprintf("%v %v", forwardInput, reverseInput)
				close(forwardChan)
				close(reverseChan)
				close(descriptionChan)
				close(isClipChan)
				close(isSlidingChan)
				close(isLeadingChan)
				close(isTrailingChan)
				close(isCropChan)
				close(isHeadCropChan)
				close(isMinLenChan)
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
			forwardButton.SetIcon(trimm.Icons.Warning)
			forwardSelected.SetText("Field is empty. Please, chose input data file.")
		case len(reverseInput) <= 0:
			reverseButton.SetIcon(trimm.Icons.Warning)
			reverseSelected.SetText("Field is empty. Please, chose input data file.")
		case len(trimm.Params.Output) <= 0:
			outputDir.SetIcon(trimm.Icons.Warning)
		default:
			exitRutine <- true
			trimm.Params.Input = fmt.Sprintf("%v %v", forwardInput, reverseInput)
			trimm.Params.Paired = "PE"

			if isClip {
				// Если тут не заработает, переставить в default
				trimm.Params.SubParams.IlluminaClip = fmt.Sprintf("ILLUMINACLIP:%v:%v:%v:%v", fastaWithAdaptersEtc.Text, seedMismatches.Text, palindromeClipThreshold.Text, simpleClipThreshold.Text)
			}
			if isSliding {
				trimm.Params.SubParams.SlidingWindow = fmt.Sprintf("SLIDINGWINDOW:%v:%v", windowSize.Text, requiredQuality.Text)
			}
			if isLeading {
				trimm.Params.SubParams.Leading = fmt.Sprintf("LEADING:%v", leading.Text)
			}
			if isTrailing {
				trimm.Params.SubParams.Trailing = fmt.Sprintf("TRAILING:%v", trailing.Text)
			}
			if isCrop {
				trimm.Params.SubParams.Crop = fmt.Sprintf("CROP:%v", crop.Text)
			}
			if isHeadCrop {
				trimm.Params.SubParams.HeadCrop = fmt.Sprintf("HEADCROP:%v", headCrop.Text)
			}
			if isMinLen {
				trimm.Params.SubParams.MinLen = fmt.Sprintf("MINLEN:%v", minLen.Text)
			}
			log.Println(trimm.Params.SubParams)

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

	return forwardSelected, reverseSelected, descriptionLabel, choseReadsForm
}

func (trimm *Trimmomatic) SelectSingleReadsFiles(window fyne.Window, commandChan chan string, exitTerminat chan bool) (*widget.Label, *widget.Label, *widget.Form) {
	trimm.Params.Input = ""
	selected := widget.NewLabel("")
	descriptionLabel := widget.NewLabel("")

	readsChan := make(chan bool)
	outputDirChan := make(chan bool)
	descriptionChan := make(chan string)
	exitRutine := make(chan bool, 1)

	isClipChan := make(chan bool, 1)
	isSlidingChan := make(chan bool, 1)
	isLeadingChan := make(chan bool, 1)
	isTrailingChan := make(chan bool, 1)
	isCropChan := make(chan bool, 1)
	isHeadCropChan := make(chan bool, 1)
	isMinLenChan := make(chan bool, 1)

	isClip := false
	isSliding := false
	isLeading := false
	isTrailing := false
	isCrop := false
	isHeadCrop := false
	isMinLen := false

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
	phredItems := trimm.ChosePhred()
	threadsEntry := trimm.ChoseThreads()
	threadsEntryItem := widget.NewFormItem("Threads", threadsEntry)
	logFileItem := trimm.SaveLogs()

	illuminaClip, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold := trimm.CreateIlluminaClip(descriptionChan, isClipChan)
	slidingWindow, windowSize, requiredQuality := trimm.CreateSlidingWindow(descriptionChan, isSlidingChan)
	leadingItem, leading := trimm.CreateLeading(descriptionChan, isLeadingChan)
	trailingItem, trailing := trimm.CreateTrailing(descriptionChan, isTrailingChan)
	cropItem, crop := trimm.CreateCrop(descriptionChan, isCropChan)
	headCropItem, headCrop := trimm.CreateHeadCrop(descriptionChan, isHeadCropChan)
	minLenItem, minLen := trimm.CreateMinLen(descriptionChan, isMinLenChan)

	go func() {
		for {
			select {
			case <-readsChan:
				button.SetIcon(trimm.Icons.Submit)
			case <-outputDirChan:
				outputDir.SetIcon(trimm.Icons.Submit)
			case description := <-descriptionChan:
				descriptionLabel.SetText(description)
			case clipVal := <-isClipChan:
				isClip = clipVal
			case slidingVal := <-isSlidingChan:
				isSliding = slidingVal
			case leadingVal := <-isLeadingChan:
				isLeading = leadingVal
			case trailingVal := <-isTrailingChan:
				isTrailing = trailingVal
			case cropVal := <-isCropChan:
				isCrop = cropVal
			case headCropVal := <-isHeadCropChan:
				isHeadCrop = headCropVal
			case minLenVal := <-isMinLenChan:
				isMinLen = minLenVal
			case <-exitRutine:
				close(readsChan)
				close(descriptionChan)
				close(isClipChan)
				close(isSlidingChan)
				close(isLeadingChan)
				close(isTrailingChan)
				close(isCropChan)
				close(isHeadCropChan)
				close(isMinLenChan)
				close(exitRutine)
				return
			}
		}
	}()

	choseReadsForm := widget.NewForm(
		widget.NewFormItem("Select input file", button),
		outputDirItem,
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
			button.SetIcon(trimm.Icons.Warning)
			selected.SetText("Field is empty. Please, chose input data file.")
		case len(trimm.Params.Output) <= 0:
			outputDir.SetIcon(trimm.Icons.Warning)
			// selected.SetText("Field is empty. Please, chose input data file.")
		default:
			if isClip {
				// Если тут не заработает, переставить в default
				trimm.Params.SubParams.IlluminaClip = fmt.Sprintf("ILLUMINACLIP:%v:%v:%v:%v", fastaWithAdaptersEtc.Text, seedMismatches.Text, palindromeClipThreshold.Text, simpleClipThreshold.Text)
			}
			if isSliding {
				trimm.Params.SubParams.SlidingWindow = fmt.Sprintf("SLIDINGWINDOW:%v:%v", windowSize.Text, requiredQuality.Text)
			}
			if isLeading {
				trimm.Params.SubParams.Leading = fmt.Sprintf("LEADING:%v", leading.Text)
			}
			if isTrailing {
				trimm.Params.SubParams.Trailing = fmt.Sprintf("TRAILING:%v", trailing.Text)
			}
			if isCrop {
				trimm.Params.SubParams.Crop = fmt.Sprintf("CROP:%v", crop.Text)
			}
			if isHeadCrop {
				trimm.Params.SubParams.HeadCrop = fmt.Sprintf("HEADCROP:%v", headCrop.Text)
			}
			if isMinLen {
				trimm.Params.SubParams.MinLen = fmt.Sprintf("MINLEN:%v", minLen.Text)
			}

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

	return selected, descriptionLabel, choseReadsForm
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
	selectThreadsCount.SetText(strconv.Itoa(len(cpuInfo) / 2))

	return selectThreadsCount
}

func (trimm *Trimmomatic) SaveLogs() *widget.FormItem {
	logFile := widget.NewCheck("logfile.log", func(b bool) {
		if b {
			trimm.Params.Logfile = trimm.Params.Output + "/logfile.log"
			// trimm.CountLogFileLength("/home/mrred/Загрузки/ERR9792312.fastq.gz", trimm.Params.Logfile, "/home/mrred/Загрузки/")
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

func (trimm *Trimmomatic) CreateIlluminaClip(descrChan chan string, isClipChan chan bool) (*widget.FormItem, *widget.SelectEntry, *widget.Entry, *widget.Entry, *widget.Entry) {
	fastaWithAdaptersEtc := widget.NewSelectEntry([]string{"TruSeq3-SE"})
	fastaWithAdaptersEtc.SetPlaceHolder(trimm.Params.SubParams.Names.IlluminaClip[0])

	seedMismatches := widget.NewEntry()
	seedMismatches.SetPlaceHolder(trimm.Params.SubParams.Names.IlluminaClip[1])

	palindromeClipThreshold := widget.NewEntry()
	palindromeClipThreshold.SetPlaceHolder(trimm.Params.SubParams.Names.IlluminaClip[2])

	simpleClipThreshold := widget.NewEntry()
	simpleClipThreshold.SetPlaceHolder(trimm.Params.SubParams.Names.IlluminaClip[3])

	// isClip := false
	isIlluminaClip := widget.NewCheck("Use", func(b bool) {
		if b {
			isClipChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.IlluminaClip
	})

	cols := container.NewGridWithColumns(6, isIlluminaClip, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold, describe)
	formItem := widget.NewFormItem("ILLUMINACLIP", cols)
	return formItem, fastaWithAdaptersEtc, seedMismatches, palindromeClipThreshold, simpleClipThreshold
}

func (trimm *Trimmomatic) CreateSlidingWindow(descrChan chan string, isSlidingChan chan bool) (*widget.FormItem, *widget.Entry, *widget.Entry) {
	windowSize := widget.NewEntry()
	windowSize.SetPlaceHolder(trimm.Params.SubParams.Names.SlidingWindow[0])

	requiredQuality := widget.NewEntry()
	requiredQuality.SetPlaceHolder(trimm.Params.SubParams.Names.SlidingWindow[1])

	isSlidingWindow := widget.NewCheck("Use", func(b bool) {
		if b {
			isSlidingChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.SlidingWindow
	})

	cols := container.NewGridWithColumns(6, isSlidingWindow, windowSize, requiredQuality, layout.NewSpacer(), layout.NewSpacer(), describe)
	// rows := container.NewGridWithRows(2, cols, isSlidingWindow)
	formItem := widget.NewFormItem("SLIDINGWINDOW", cols)
	return formItem, windowSize, requiredQuality
}

func (trimm *Trimmomatic) CreateLeading(descrChan chan string, isLeadingChan chan bool) (*widget.FormItem, *widget.Entry) {
	quality := widget.NewEntry()
	quality.SetPlaceHolder(trimm.Params.SubParams.Names.Leading)

	isQualityCheck := widget.NewCheck("Use", func(b bool) {
		if b {
			isLeadingChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.Leading
	})

	cols := container.NewGridWithColumns(6, isQualityCheck, quality, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem("LEADING", cols)
	return formItem, quality
}

func (trimm *Trimmomatic) CreateTrailing(descrChan chan string, isTrailingChan chan bool) (*widget.FormItem, *widget.Entry) {
	quality := widget.NewEntry()
	quality.SetPlaceHolder(trimm.Params.SubParams.Names.Trailing)

	isQualityCheck := widget.NewCheck("Use", func(b bool) {
		if b {
			isTrailingChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.Trailing
	})

	cols := container.NewGridWithColumns(6, isQualityCheck, quality, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem("TRAILING", cols)
	return formItem, quality
}

func (trimm *Trimmomatic) CreateCrop(descrChan chan string, isCropChan chan bool) (*widget.FormItem, *widget.Entry) {
	length := widget.NewEntry()
	length.SetPlaceHolder(trimm.Params.SubParams.Names.Crop)

	isLengthCheck := widget.NewCheck("Use", func(b bool) {
		if b {
			isCropChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.Crop
	})

	cols := container.NewGridWithColumns(6, isLengthCheck, length, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem("CROP", cols)
	return formItem, length
}

func (trimm *Trimmomatic) CreateHeadCrop(descrChan chan string, isHeadCropChan chan bool) (*widget.FormItem, *widget.Entry) {
	length := widget.NewEntry()
	length.SetPlaceHolder(trimm.Params.SubParams.Names.HeadCrop)

	isLengthCheck := widget.NewCheck("Use", func(b bool) {
		if b {
			isHeadCropChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.HeadCrop
	})

	cols := container.NewGridWithColumns(6, isLengthCheck, length, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem("HEADCROP", cols)
	return formItem, length
}

func (trimm *Trimmomatic) CreateMinLen(descrChan chan string, isMinLenChan chan bool) (*widget.FormItem, *widget.Entry) {
	length := widget.NewEntry()
	length.SetPlaceHolder(trimm.Params.SubParams.Names.MinLen)

	isLengthCheck := widget.NewCheck("Use", func(b bool) {
		if b {
			isMinLenChan <- b
		}
	})

	describe := widget.NewButtonWithIcon("", trimm.Icons.Question, func() {
		descrChan <- trimm.Params.SubParams.Description.MinLen
	})
	// cntr := container.NewGridWrap(fyne.NewSize(trimm.Params.SubParams.QuestionWidth, length.MinSize().Height), describe)
	// cntr := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), describe)

	cols := container.NewGridWithColumns(6, isLengthCheck, length, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem("MINLEN", cols)
	return formItem, length
}

// func (trimm *Trimmomatic) CountLogFileLength(inputPath, logfilePath, target string) {
// 	// inputFile, err := os.ReadFile(inputPath)
// 	// file, err := os.Open(inputPath)

// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// reader := fastq.NewReader(file, linear.NewQSeq("", nil, alphabet.DNAgapped, alphabet.Illumina1_3))
// 	// r, err := reader.Read()
// 	// log.Println(r, err)
// 	go func() {
// 		for {
// 			time.Sleep(50 * time.Millisecond)
// 			file, err := os.Stat(logfilePath)
// 			log.Println(file.Size(), err)
// 		}
// 	}()
// 	return
// }
