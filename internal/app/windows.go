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

	go func() {
		for {
			select {
			case command := <-commandsChan:
				if len(command) > 0 {
					t.Write([]byte(command + "\n"))
					// t.
				}
			case <-exitRutine:
				close(commandsChan)
				close(exitRutine)
				return
			}
		}
	}()

	return terminalWindow
}
