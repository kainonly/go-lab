package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

func main() {
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
	for i := 0; i < 200; i++ {
		if i == 0 {
			bar.ChangeMax(100)
		}
		if i == 99 {
			bar.ChangeMax(200)
		}
		bar.Add(1)
		bar.Describe(fmt.Sprintf(`%d`, i*2))
		time.Sleep(40 * time.Millisecond)
	}
}
