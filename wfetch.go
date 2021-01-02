package main

import (
	"fmt"
	"strings"
)

func main() {
	info := GetInfo()
	formattedInfo := strings.Split(FormatInfo(info), "\n")

	artLines := strings.Split(Art, "\n")
	numOfPaddingLines := (len(artLines) - len(formattedInfo)) / 2
	accent := GetAccentColor()
	accentAnsii := fmt.Sprintf("\033[38;2;%v;%v;%vm", accent.R, accent.G, accent.B)
	resetAnsii := "\033[0m"

	fmt.Println()
	for i := 0; i < len(artLines); i++ {
		if i < numOfPaddingLines || i-numOfPaddingLines > len(formattedInfo)-1 {
			fmt.Println(accentAnsii + artLines[i] + resetAnsii)
		} else {
			fmt.Printf("%s%s     %s\n", accentAnsii, artLines[i], formattedInfo[i-numOfPaddingLines])
		}
	}
	fmt.Println()
}
