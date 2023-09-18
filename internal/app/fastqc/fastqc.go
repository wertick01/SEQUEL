package fastqc

import (
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

func (fastQC *FastQC) CreateFastQCForn(window fyne.Window, commandChan chan string) (*widget.Form, *widget.Label) {
	descriptionLabel := widget.NewLabel("")
	descriptionLabel.Alignment = fyne.TextAlignCenter
	inputFileLabel := widget.NewLabel("")

	// isSelectedChan := make(chan bool, 1)
	descriptionChan := make(chan string, 1)

	inputFileChan := make(chan bool, 1)    //
	outDirChan := make(chan bool, 1)       //
	casavaChan := make(chan bool, 1)       //
	nanoChan := make(chan bool, 1)         //
	noFilterChan := make(chan bool, 1)     //
	extractChan := make(chan bool, 1)      //
	javaChan := make(chan bool, 1)         //
	noExtractChan := make(chan bool, 1)    //
	noGroupChan := make(chan bool, 1)      //
	minLengthChan := make(chan bool, 1)    //
	formatChan := make(chan bool, 1)       //
	threadsChan := make(chan bool, 1)      //
	contaminantsChan := make(chan bool, 1) //
	adaptersChan := make(chan bool, 1)     //
	limitsChan := make(chan bool, 1)       //
	kMersChan := make(chan bool, 1)        //
	// quietChan := make(chan bool, 1)        //
	dirChan := make(chan bool, 1) //

	javaSelectedChan := make(chan bool, 1)         //
	contaminantsSelectedChan := make(chan bool, 1) //
	adaptersSelectedChan := make(chan bool, 1)     //
	limitsSelectedChan := make(chan bool, 1)       //

	exitRutine := make(chan bool, 1)

	isOutDir := false       //
	isCasava := false       //
	isNano := false         //
	isNoFilter := false     //
	isExtract := false      //
	isJava := false         //
	isNoExtract := false    //
	isNoGroup := false      //
	isMinLength := false    //
	isFormat := false       //
	isThreads := false      //
	isContaminants := false //
	isAdapters := false     //
	isLimits := false       //
	isKMers := false        //
	// isQuiet := false        //
	isDir := false //

	selectFileButton := widget.NewButton("Input File", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					inputFileLabel.SetText(r.URI().Path())
					if len(inputFileLabel.Text) > 0 {
						inputFileChan <- true
					} else {
						inputFileChan <- false
					}
				}
			},
			window,
		)
	})

	inputFileItem := widget.NewFormItem("Input File", selectFileButton)
	outputDirItem, outputDirButton, outputDirLabel := fastQC.ChoseOutputDir(window, outDirChan)
	threadsItem, threadsCountEntry := fastQC.ChoseThreads()
	casavaItem := fastQC.CreateCasavaItem(descriptionChan, casavaChan)
	nanoItem := fastQC.CreateNanoItem(descriptionChan, nanoChan)
	noFilterItem := fastQC.CreateNoFilterItem(descriptionChan, noFilterChan)
	extractItem := fastQC.CreateExtractItem(descriptionChan, extractChan)
	javaItem, javaPathToBin, JavaBinButton := fastQC.CreateJavaItem(window, descriptionChan, javaChan, javaSelectedChan)
	noExtractItem := fastQC.CreateNoExtractItem(descriptionChan, noExtractChan)
	noGroupItem := fastQC.CreateNoGroupItem(descriptionChan, noGroupChan)
	minLengthItem, minLenghtLabel := fastQC.CreateMinLengthItem(descriptionChan, minLengthChan)
	formatItem, formatLabel := fastQC.CreateFormatItem(descriptionChan, formatChan)
	contaminantsItem, contaminantsFilePath, contaminantsFileButton := fastQC.CreateContaminantsItem(window, descriptionChan, contaminantsChan, contaminantsSelectedChan)
	adaptersItem, adaptersFilePath, adaptersFileButton := fastQC.CreateAdaptersItem(window, descriptionChan, adaptersChan, adaptersSelectedChan)
	limitsItem, limitsFilePath, limitsFileButton := fastQC.CreateLimitsItem(window, descriptionChan, limitsChan, limitsSelectedChan)
	kMersItem, kMersSliderLabel := fastQC.CreateKMersItem(descriptionChan, kMersChan)
	suppDirItem, supDirLabel := fastQC.ChoseSupDir(window, dirChan)

	go func() {
		for {
			select {
			case <-inputFileChan:
				selectFileButton.SetIcon(fastQC.Icons.Submit)
			// case <-isSelectedChan:
			// 	selectFileButton.SetIcon(fastQC.Icons.Submit)
			case description := <-descriptionChan:
				descriptionLabel.SetText(description)
			case outDir := <-outDirChan:
				outputDirButton.SetIcon(fastQC.Icons.Submit)
				isOutDir = outDir
			case casava := <-casavaChan:
				isCasava = casava
			case nano := <-nanoChan:
				isNano = nano
			case noFilter := <-noFilterChan:
				isNoFilter = noFilter
			case extract := <-extractChan:
				isExtract = extract
			case java := <-javaChan:
				isJava = java
			case noExtract := <-noExtractChan:
				isNoExtract = noExtract
			case noGroup := <-noGroupChan:
				isNoGroup = noGroup
			case minLength := <-minLengthChan:
				isMinLength = minLength
			case format := <-formatChan:
				isFormat = format
			case threads := <-threadsChan:
				isThreads = threads
			case contaminants := <-contaminantsChan:
				isContaminants = contaminants
			case adapters := <-adaptersChan:
				isAdapters = adapters
			case limits := <-limitsChan:
				isLimits = limits
			case kMers := <-kMersChan:
				isKMers = kMers
			case <-javaSelectedChan:
				JavaBinButton.SetIcon(fastQC.Icons.Submit)
			case <-contaminantsSelectedChan:
				contaminantsFileButton.SetIcon(fastQC.Icons.Submit)
			case <-adaptersSelectedChan:
				adaptersFileButton.SetIcon(fastQC.Icons.Submit)
			case <-limitsSelectedChan:
				limitsFileButton.SetIcon(fastQC.Icons.Submit)

			// case quiet := <-quietChan:
			// 	isQuiet = quiet
			case dir := <-dirChan:
				isDir = dir

			// case progress := <-ratioChan:
			// 	progressBar.Value = progress
			// 	progressBar.Refresh()

			case <-exitRutine:
				fastQC.InputFile.Value = fmt.Sprintf("%v", inputFileLabel.Text)
				// close(isSelectedChan)
				close(descriptionChan)
				close(inputFileChan)
				close(outDirChan)
				close(casavaChan)
				close(nanoChan)
				close(noFilterChan)
				close(extractChan)
				close(javaChan)
				close(noExtractChan)
				close(noGroupChan)
				close(minLengthChan)
				close(formatChan)
				close(threadsChan)
				close(contaminantsChan)
				close(adaptersChan)
				close(limitsChan)
				close(kMersChan)
				// close(quietChan)
				close(dirChan)
				close(contaminantsSelectedChan) //
				close(adaptersSelectedChan)     //
				close(limitsSelectedChan)       //
				close(exitRutine)
				return
			}
		}
	}()

	fastQCForm := widget.NewForm(
		inputFileItem,
		outputDirItem,
		threadsItem,
		casavaItem,
		nanoItem,
		noFilterItem,
		extractItem,
		javaItem,
		noExtractItem,
		noGroupItem,
		minLengthItem,
		formatItem,
		contaminantsItem,
		adaptersItem,
		limitsItem,
		kMersItem,
		suppDirItem,
	)

	fastQCForm.SubmitText = "Run"
	fastQCForm.CancelText = "Close"
	// fastQCForm.Refresh()

	fastQCForm.OnSubmit = func() {
		switch {
		case len(inputFileLabel.Text) <= 0:
			selectFileButton.SetIcon(fastQC.Icons.Warning)
			descriptionLabel.SetText("Field is empty. Please, chose input data file.")
		case !isOutDir:
			outputDirButton.SetIcon(fastQC.Icons.Warning)
			descriptionLabel.SetText("Field is empty. Please, chose input data file.")
		default:
			fastQC.InputFile.Value = inputFileLabel.Text
			fastQC.InputFile.IsUsed = true

			if fastQC.OutputDir.IsUsed = isOutDir; isOutDir {
				fastQC.OutputDir.Value = outputDirLabel.Text
			}
			if fastQC.Casava.IsUsed = isCasava; isCasava {
				fastQC.Casava.Value = ""
			}
			if fastQC.Nano.IsUsed = isNano; isNano {
				fastQC.Nano.Value = ""
			}
			if fastQC.NoFilter.IsUsed = isNoFilter; isNoFilter {
				fastQC.NoFilter.Value = ""
			}
			if fastQC.Extract.IsUsed = isExtract; isExtract {
				fastQC.Extract.Value = ""
			}
			if fastQC.Java.IsUsed = isJava; isJava {
				fastQC.Java.Value = javaPathToBin.Text
			}
			if fastQC.NoExtract.IsUsed = isNoExtract; isNoExtract {
				fastQC.NoExtract.Value = ""
			}
			if fastQC.NoGroup.IsUsed = isNoGroup; isNoGroup {
				fastQC.NoGroup.Value = ""
			}
			if fastQC.MinLength.IsUsed = isMinLength; isMinLength {
				fastQC.MinLength.Value = minLenghtLabel.Text
			}
			if fastQC.Format.IsUsed = isFormat; isFormat {
				fastQC.Format.Value = formatLabel.Text
			}
			if fastQC.Threads.IsUsed = isThreads; isThreads {
				fastQC.Threads.Value = threadsCountEntry.Text
			}
			if fastQC.Contaminants.IsUsed = isContaminants; isContaminants {
				fastQC.Contaminants.Value = contaminantsFilePath.Text
			}
			if fastQC.Adapters.IsUsed = isAdapters; isAdapters {
				fastQC.Adapters.Value = adaptersFilePath.Text
			}
			if fastQC.Limits.IsUsed = isLimits; isLimits {
				fastQC.Limits.Value = limitsFilePath.Text
			}

			if fastQC.KMers.IsUsed = isKMers; isKMers {
				fastQC.KMers.Value = kMersSliderLabel.Text
			}

			if fastQC.Dir.IsUsed = isDir; isDir {
				fastQC.Dir.Value = supDirLabel.Text
			}

			fastQC.BuildCommand()
		}
	}

	fastQCForm.OnCancel = func() {
		exitRutine <- true
		window.Close()
	}

	return fastQCForm, descriptionLabel
}

func (fastQC *FastQC) ChoseOutputDir(window fyne.Window, dirChan chan bool) (*widget.FormItem, *widget.Button, *widget.Entry) {
	outputDirEntry := widget.NewEntry()
	outputDirButton := widget.NewButton("Change output dir", func() {
		dialog.ShowFolderOpen(
			func(r fyne.ListableURI, err error) {
				if r != nil {
					outputDirEntry.SetText(r.Path())
					if len(outputDirEntry.Text) > 0 {
						dirChan <- true
					} else {
						dirChan <- false
					}
				}
			},
			window,
		)
	})

	return widget.NewFormItem("Output Directory", outputDirButton), outputDirButton, outputDirEntry
}

func (fastQC *FastQC) ChoseThreads() (*widget.FormItem, *widget.SelectEntry) {
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

	selectThreadsItem := widget.NewFormItem(fastQC.Threads.Name, selectThreadsCount)

	return selectThreadsItem, selectThreadsCount
}

func (fastQC *FastQC) CreateCasavaItem(descrChan chan string, casavaChan chan bool) *widget.FormItem {
	isCasava := widget.NewCheck("Use", func(b bool) {
		// if b {
		casavaChan <- b
		// }
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Casava.Description
	})

	cols := container.NewGridWithColumns(6, isCasava, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Casava.Name, cols)
	return formItem
}

func (fastQC *FastQC) CreateNanoItem(descrChan chan string, nanoChan chan bool) *widget.FormItem {
	isNano := widget.NewCheck("Use", func(b bool) {
		// if b {
		nanoChan <- b
		// }
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Nano.Description
	})

	cols := container.NewGridWithColumns(6, isNano, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Nano.Name, cols)
	return formItem
}

func (fastQC *FastQC) CreateNoFilterItem(descrChan chan string, noFilterChan chan bool) *widget.FormItem {
	isNoFilter := widget.NewCheck("Use", func(b bool) {
		// if b {
		noFilterChan <- b
		// }
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.NoFilter.Description
	})

	cols := container.NewGridWithColumns(6, isNoFilter, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.NoFilter.Name, cols)
	return formItem
}

func (fastQC *FastQC) CreateExtractItem(descrChan chan string, extractChan chan bool) *widget.FormItem {
	isExtract := widget.NewCheck("Use", func(b bool) {
		// if b {
		extractChan <- b
		// }
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Extract.Description
	})

	cols := container.NewGridWithColumns(6, isExtract, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Extract.Name, cols)
	return formItem
}

func (fastQC *FastQC) CreateJavaItem(window fyne.Window, descrChan chan string, javaChan chan bool, selectedChan chan bool) (*widget.FormItem, *widget.Entry, *widget.Button) {
	isJava := widget.NewCheck("Use", func(b bool) {
		// if b {
		javaChan <- b
		// }
	})

	javaPath := widget.NewEntry()

	selectJavaBinFileButton := widget.NewButton("Chose file", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					javaPath.SetText(r.URI().Path())
					if len(javaPath.Text) > 0 {
						isJava.Checked = true // При выборе файла галочка должна проставиться автоматически
						isJava.Refresh()      // затем её можно будет снять вручную
						selectedChan <- true
					} else {
						selectedChan <- false
					}
				}
			},
			window,
		)
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Java.Description
	})

	cols := container.NewGridWithColumns(6, isJava, selectJavaBinFileButton, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Java.Name, cols)
	return formItem, javaPath, selectJavaBinFileButton
}

func (fastQC *FastQC) CreateNoExtractItem(descrChan chan string, noExtractChan chan bool) *widget.FormItem {
	isNoExtract := widget.NewCheck("Use", func(b bool) {
		// if b {
		noExtractChan <- b
		// }
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.NoExtract.Description
	})

	cols := container.NewGridWithColumns(6, isNoExtract, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.NoExtract.Name, cols)
	return formItem
}

func (fastQC *FastQC) CreateNoGroupItem(descrChan chan string, noGroupChan chan bool) *widget.FormItem {
	isNoGroup := widget.NewCheck("Use", func(b bool) {
		// if b {
		noGroupChan <- b
		// }
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.NoGroup.Description
	})

	cols := container.NewGridWithColumns(6, isNoGroup, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.NoGroup.Name, cols)
	return formItem
}

func (fastQC *FastQC) CreateMinLengthItem(descrChan chan string, minLengthChan chan bool) (*widget.FormItem, *widget.Entry) {
	isMinLength := widget.NewCheck("Use", func(b bool) {
		// if b {
		minLengthChan <- b
		// }
	})

	minLengthEntry := widget.NewEntry()
	minLengthEntry.SetPlaceHolder(fastQC.MinLength.Flag)

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.MinLength.Description
	})

	cols := container.NewGridWithColumns(6, isMinLength, minLengthEntry, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.MinLength.Name, cols)
	return formItem, minLengthEntry
}

func (fastQC *FastQC) CreateFormatItem(descrChan chan string, formatChan chan bool) (*widget.FormItem, *widget.SelectEntry) {
	isFormat := widget.NewCheck("Use", func(b bool) {
		// if b {
		formatChan <- b
		// }
	})

	formatEntry := widget.NewSelectEntry(fastQC.Format.Formats)
	formatEntry.SetPlaceHolder(fastQC.Format.Flag)

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Format.Description
	})

	cols := container.NewGridWithColumns(6, isFormat, formatEntry, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Format.Name, cols)
	return formItem, formatEntry
}

// need to create one more channel for "accept/warning" file selection file
func (fastQC *FastQC) CreateContaminantsItem(window fyne.Window, descrChan chan string, contaminantsChan chan bool, selectedChan chan bool) (*widget.FormItem, *widget.Entry, *widget.Button) {
	isContaminants := widget.NewCheck("Use", func(b bool) {
		// if b {
		contaminantsChan <- b
		// }
	})

	contaminantFile := widget.NewEntry()

	selectContaminantsFileButton := widget.NewButton("Chose file", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					contaminantFile.SetText(r.URI().Path())
					if len(contaminantFile.Text) > 0 {
						isContaminants.Checked = true // При выборе файла галочка должна проставиться автоматически
						isContaminants.Refresh()      // затем её можно будет снять вручную
						selectedChan <- true
					} else {
						selectedChan <- false
					}
				}
			},
			window,
		)
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Contaminants.Description
	})

	cols := container.NewGridWithColumns(6, isContaminants, selectContaminantsFileButton, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Contaminants.Name, cols)
	return formItem, contaminantFile, selectContaminantsFileButton
}

// need to create one more channel for "accept/warning" file selection file
func (fastQC *FastQC) CreateAdaptersItem(window fyne.Window, descrChan chan string, adaptersChan chan bool, selectedChan chan bool) (*widget.FormItem, *widget.Entry, *widget.Button) {
	isAdapters := widget.NewCheck("Use", func(b bool) {
		// if b {
		adaptersChan <- b
		// }
	})

	adaptersFile := widget.NewEntry()

	selectAdaptersFileButton := widget.NewButton("Chose file", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					adaptersFile.SetText(r.URI().Path())
					if len(adaptersFile.Text) > 0 {
						isAdapters.Checked = true // При выборе файла галочка должна проставиться автоматически
						isAdapters.Refresh()      // затем её можно будет снять вручную
						selectedChan <- true
					} else {
						selectedChan <- false
					}
				}
			},
			window,
		)
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Adapters.Description
	})

	cols := container.NewGridWithColumns(6, isAdapters, selectAdaptersFileButton, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Adapters.Name, cols)
	return formItem, adaptersFile, selectAdaptersFileButton
}

// need to create one more channel for "accept/warning" file selection file
func (fastQC *FastQC) CreateLimitsItem(window fyne.Window, descrChan chan string, limitsChan chan bool, selectedChan chan bool) (*widget.FormItem, *widget.Entry, *widget.Button) {
	isLimits := widget.NewCheck("Use", func(b bool) {
		// if b {
		limitsChan <- b
		// }
	})

	limitsFile := widget.NewEntry()

	selectLimitsFileButton := widget.NewButton("Chose file", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					limitsFile.SetText(r.URI().Path())
					if len(limitsFile.Text) > 0 {
						isLimits.Checked = true // При выборе файла галочка должна проставиться автоматически
						isLimits.Refresh()      // затем её можно будет снять вручную
						selectedChan <- true
					} else {
						selectedChan <- false
					}
				}
			},
			window,
		)
	})

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.Limits.Description
	})

	cols := container.NewGridWithColumns(6, isLimits, selectLimitsFileButton, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), describe)
	formItem := widget.NewFormItem(fastQC.Limits.Name, cols)
	return formItem, limitsFile, selectLimitsFileButton
}

func (fastQC *FastQC) CreateKMersItem(descrChan chan string, isKMerChan chan bool) (*widget.FormItem, *widget.Label) {
	kMerSlider := widget.NewSlider(2, 10)
	kMerSlider.Step = 1
	kMerSlider.SetValue(7)
	kMerSlider.Refresh()

	kMersLabel := widget.NewLabel("")
	sliderValueLabel := widget.NewLabel("")
	sliderValueLabel.SetText("7 - K-Mers length")

	isQualityCheck := widget.NewCheck("Use", func(b bool) {
		// if b {
		isKMerChan <- b
		kMerSlider.OnChangeEnded = func(f float64) {
			kMersVal := strconv.FormatFloat(f, 'f', -1, 64)
			kMersLabel.SetText(kMersVal)
			// }
		}
	})

	kMerSlider.OnChanged = func(f float64) {
		kMersVal := strconv.FormatFloat(f, 'f', -1, 64)
		sliderValueLabel.SetText(fmt.Sprintf("%v - K-Mers length", kMersVal))
	}

	describe := widget.NewButtonWithIcon("", fastQC.Icons.Question, func() {
		descrChan <- fastQC.KMers.Description
	})

	// checkCol := container.NewGridWithColumns(2, isQualityCheck, layout.NewSpacer())
	describeCol := container.NewGridWithColumns(2, sliderValueLabel, describe)
	// sliderCol := container.NewGridWithRows(2, sliderValueLabel, kMerSlider)

	cols := container.NewGridWithColumns(3, isQualityCheck, kMerSlider, describeCol)
	formItem := widget.NewFormItem(fastQC.KMers.Name, cols)
	return formItem, kMersLabel
}

func (fastQC *FastQC) ChoseSupDir(window fyne.Window, supDirChan chan bool) (*widget.FormItem, *widget.Label) {
	supDirLabel := widget.NewLabel("")
	aupDirButton := widget.NewButton("Change Directory", func() {
		dialog.ShowFolderOpen(
			func(r fyne.ListableURI, err error) {
				if r != nil {
					supDirLabel.SetText(r.Path())
					if len(fastQC.OutputDir.Value) > 0 {
						supDirChan <- true
					} else {
						supDirChan <- false
					}
				}
			},
			window,
		)
	})

	return widget.NewFormItem("Supplementary dir", aupDirButton), supDirLabel
}
