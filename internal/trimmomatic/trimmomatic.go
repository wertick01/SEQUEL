package trimmomatic

import (
	"biolink-gui/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

// func (trimm *Trimmomatic) ReadParams() (*widget.Button, error) {
// 	reads := map[string]string{
// 		"Paired": "PE",
// 		"Single": "SE",
// 	}

// 	var newErr error

// 	botton := widget.NewButton("GO", func() {
// 		phred, err := strconv.Atoi(trimm.Phred.Text)
// 		if err != nil {
// 			newErr = err
// 		}
// 		trimm.Params.Phred = phred

// 		threads, err := strconv.Atoi(trimm.Threads.Text)
// 		if err != nil {
// 			newErr = err
// 		}
// 		trimm.Params.Threads = threads

// 		trimm.Params.Path = trimm.Path.Text
// 		trimm.Params.Input = trimm.Input.Text
// 		trimm.Params.Output = trimm.Output.Text
// 		trimm.Params.Paired = reads[trimm.Paired.Selected]
// 		log.Println(trimm.Params)
// 	})

// 	return botton, newErr
// }

func (trimm *Trimmomatic) SelectPairedReadsFiles(window fyne.Window) (*widget.Button, *widget.Button) {
	trimm.Params.Input = ""

	choseForwardReadsBotton := widget.NewButton("Forward reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				trimm.Params.Input += r.URI().Path() + " "
			},
			window,
		)
	})

	choseReverseReadsBotton := widget.NewButton("Reverse reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				trimm.Params.Input += r.URI().Path()
			},
			window,
		)
	})

	return choseForwardReadsBotton, choseReverseReadsBotton
}

func (trimm *Trimmomatic) SelectSingleReadsFiles(window fyne.Window) *widget.Button {
	trimm.Params.Input = ""

	choseReadsBotton := widget.NewButton("Reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				trimm.Params.Input += r.URI().Path() + " "
			},
			window,
		)
	})

	return choseReadsBotton
}
