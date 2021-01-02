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

	for i := 0; i < len(artLines); i++ {
		if i < numOfPaddingLines || i-numOfPaddingLines > len(formattedInfo)-1 {
			fmt.Println(artLines[i])
		} else {
			fmt.Printf("%s%s\n", artLines[i], formattedInfo[i-numOfPaddingLines])
		}
	}
}
