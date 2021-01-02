package main

import (
	"image"
	"os"
	"os/exec"

	"github.com/EdlinOrg/prominentcolor"
)

// GetAccentColor : gets the most prominent color in the users desktop wallpaper
func GetAccentColor() prominentcolor.ColorRGB {
	wallpapers := getValuesFromList(exec.Command("wmic", "desktop", "get", "wallpaper", "/value").Output())
	var wallpaper string
	for _, w := range wallpapers {
		if w != "" {
			wallpaper = w
		}
	}

	img, err := loadImage(wallpaper)
	if err != nil {
		panic(err)
	}

	colors, err := prominentcolor.Kmeans(img)
	if err != nil {
		panic(err)
	}

	return colors[1].Color
}

func loadImage(fileInput string) (image.Image, error) {
	f, err := os.Open(fileInput)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
