package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"time"
)

func main() {
	color.Cyan("Prints text in cyan.")
	color.Blue("Prints %s in blue.", "text")
	color.Red("We have red")
	color.Magenta("And many others ..")
	bar := progressbar.NewOptions(1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	for i := 0; i < 100; i++ {
		if i == 0 {
			bar.ChangeMax(50)
		}
		if i == 49 {
			bar.ChangeMax(100)
		}
		bar.Add(1)
		bar.Describe(fmt.Sprintf(`%d`, i*2))
		time.Sleep(10 * time.Millisecond)
	}
}
