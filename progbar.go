package smbget

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)


func GetProgressBar(size int64, descr string, padding string) *progressbar.ProgressBar {
    return progressbar.NewOptions64(size, 
        progressbar.OptionEnableColorCodes(true), 
        progressbar.OptionShowBytes(true), 
        progressbar.OptionShowCount(), 
        progressbar.OptionShowElapsedTimeOnFinish(), 
        progressbar.OptionSetVisibility(true), 
        progressbar.OptionSetElapsedTime(true),
        progressbar.OptionOnCompletion(func() {
            fmt.Printf("\n")
        }),
        progressbar.OptionSetDescription(descr + padding),
        progressbar.OptionSetTheme(progressbar.Theme{
            Saucer:        "[green]=[reset]",
            SaucerHead:    "[green]>[reset]",
            SaucerPadding: "[white]•[reset]",
            BarStart:       "[green][",
            BarEnd:       "[green]][reset]",
        }))
}