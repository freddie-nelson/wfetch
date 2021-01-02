package main

import (
	"fmt"
	"os"
)

// AccentAnsii the ansii code for the 2nd most prominent color in the users desktop wallpaper
var AccentAnsii string

// ResetAnsii the ansii code to reset terminal output to its default color
var ResetAnsii string

func init() {
	ResetAnsii = "\033[0m"

	if len(os.Args) > 1 && os.Args[1] == "--wall" {
		accent := GetAccentColor()
		AccentAnsii = fmt.Sprintf("\033[38;2;%v;%v;%vm", accent.R, accent.G, accent.B)
	} else {
		AccentAnsii = fmt.Sprintf("\033[95m")
	}
}
