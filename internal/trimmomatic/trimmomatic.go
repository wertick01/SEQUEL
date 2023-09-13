package trimmomatic

import (
	"biolink-nipt-gui/internal/models"
	"fmt"
	"log"

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

func (trimm *Trimmomatic) SelectPairedReadsFiles(window fyne.Window) (*widget.Label, *widget.Label, *widget.Form) {
	trimm.Params.Input = ""
	forwardInput, reverseInput := "", ""
	forwardSelected := widget.NewLabel("")
	reverseSelected := widget.NewLabel("")

	forwardButton := widget.NewButton("Forward Reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					forwardInput += r.URI().Path()
					forwardSelected.SetText(r.URI().Path())
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
				}
			},
			window,
		)
	})

	choseReadsForm := widget.NewForm(
		widget.NewFormItem("Select input forward file", forwardButton),
		widget.NewFormItem("Select input reverse file", reverseButton),
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
		default:
			trimm.Params.Input = fmt.Sprintf("%v %v", forwardInput, reverseInput)
			window.Close()
		}
	}

	// choseReadsForm.OnCancel = func() {
	// 	window.Close()
	// }

	return forwardSelected, reverseSelected, choseReadsForm
}

func (trimm *Trimmomatic) SelectSingleReadsFiles(window fyne.Window) (*widget.Label, *widget.Form) {
	trimm.Params.Input = ""
	selected := widget.NewLabel("")

	button := widget.NewButton("Reads", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					trimm.Params.Input += r.URI().Path()
					selected.SetText(r.URI().Path())
				}
			},
			window,
		)
	})

	choseReadsForm := widget.NewForm(
		widget.NewFormItem("Select input file", button),
	)

	choseReadsForm.OnSubmit = func() {
		if len(trimm.Params.Input) <= 0 {
			buttonIcon, err := fyne.LoadResourceFromPath("images/warning.png")
			if err != nil {
				log.Println(err)
			}
			button.SetIcon(buttonIcon)
			selected.SetText("Field is empty. Please, chose input data file.")
		} else {
			window.Close()
		}
	}

	// choseReadsForm.OnCancel = func() {
	// 	window.Close()
	// }

	return selected, choseReadsForm
}
