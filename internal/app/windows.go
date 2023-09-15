package app

import (
	"fyne.io/fyne/v2"
	"github.com/fyne-io/terminal"
)

func CreateTerminalWindow(newApp *App, commandsChan chan string, exitRutine chan bool, X, Y float32, window fyne.Window) fyne.Window {
	terminalWindow := newApp.App.NewWindow("Terminal")
	t := terminal.New()
	terminalWindow.SetContent(t)

	go func() {
		_ = t.RunLocalShell()
		newApp.App.Quit()
	}()

	// command := ""

	go func() {
		for {
			select {
			case command := <-commandsChan:
				if len(command) > 0 {
					// commandLabel.Text = fmt.Sprintf("%v> %v\n", commandLabel.Text, command)
					t.Write([]byte(command))
					// t.
				}
			// case <-reverseChan:
			// 	buttonIcon, err := fyne.LoadResourceFromPath("images/accept.png")
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	reverseButton.SetIcon(buttonIcon)
			case <-exitRutine:
				close(commandsChan)
				close(exitRutine)
				return
			}
		}
	}()

	// newVBox1 := container.NewVBox(t)
	// newVBox1.Resize(fyne.NewSize(X, Y))
	// newVBox1.Move(fyne.NewPos(0, 100))
	// window.SetContent(newVBox1)

	return terminalWindow
}
