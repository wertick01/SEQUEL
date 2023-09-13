package pkg

import (
	"github.com/kbinani/screenshot"
)

func GetDisplayParams() map[int]map[string]int {
	n := screenshot.NumActiveDisplays()

	var activeDisplays = map[int]map[string]int{}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		activeDisplays[i] = map[string]int{
			"X": bounds.Dx(),
			"Y": bounds.Dy(),
		}
	}

	return activeDisplays
}
