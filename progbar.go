package smbget

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(size int64, descr string, padding string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetVisibility(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetDescription("[light_magenta]"+descr+padding+"[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[light_green]=[reset]",
			SaucerHead:    "[light_green]>[reset]",
			SaucerPadding: "[dark_gray]•[reset]",
			BarStart:      "[cyan][",
			BarEnd:        "[cyan]][reset]",
		}))
}
