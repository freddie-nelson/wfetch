package main

import "fmt"

// AccentAnsii the ansii code for the 2nd most prominent color in the users desktop wallpaper
var AccentAnsii string

// ResetAnsii the ansii code to reset terminal output to its default color
var ResetAnsii string

func init() {
	accent := GetAccentColor()
	AccentAnsii = fmt.Sprintf("\033[38;2;%v;%v;%vm", accent.R, accent.G, accent.B)
	ResetAnsii = "\033[0m"
}
