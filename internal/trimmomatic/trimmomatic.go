package trimmomatic

import (
	"biolink-gui/internal/models"
	"log"
	"strconv"

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

func (trimm *Trimmomatic) ReadParams() *widget.Button {

	reads := map[string]string{
		"Paired": "PE",
		"Single": "SE",
	}

	trimm.Path = widget.NewEntry()
	trimm.Path.SetPlaceHolder("Trimmomatic path:")

	trimm.Input = widget.NewEntry()
	trimm.Input.SetPlaceHolder("Input file path:")

	trimm.Output = widget.NewEntry()
	trimm.Output.SetPlaceHolder("Output file path:")

	trimm.Paired = widget.NewRadioGroup([]string{"Paired", "Single"}, func(b string) {})
	// entryPaired.
	// entryPaired.SetPlaceHolder("Is paired:")
	// entryPaired.

	trimm.Phred = widget.NewEntry()
	trimm.Phred.SetPlaceHolder("Phred:")

	trimm.Threads = widget.NewEntry()
	trimm.Threads.SetPlaceHolder("Threads:")

	// res := widget.NewLabel("")
	log.Println(1)

	createRunCommand := widget.NewButton("GO", func() {
		phred, err := strconv.Atoi(trimm.Phred.Text)
		if err != nil {
			log.Println("ERR")
		}
		trimm.Params.Phred = phred

		threads, err := strconv.Atoi(trimm.Threads.Text)
		if err != nil {
			log.Println("ERR")
		}
		trimm.Params.Threads = threads

		trimm.Params.Path = trimm.Path.Text
		trimm.Params.Input = trimm.Input.Text
		trimm.Params.Output = trimm.Output.Text
		trimm.Params.Paired = reads[trimm.Paired.Selected]

		log.Println(trimm.Params)
	})

	log.Println(3)

	return createRunCommand
}
