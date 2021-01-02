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

	fmt.Println()
	for i := 0; i < len(artLines); i++ {
		if i < numOfPaddingLines || i-numOfPaddingLines > len(formattedInfo)-1 {
			fmt.Println(AccentAnsii + artLines[i] + ResetAnsii)
		} else {
			fmt.Printf("%s%s     %s\n", AccentAnsii, artLines[i], formattedInfo[i-numOfPaddingLines])
		}
	}
	fmt.Println()
}
