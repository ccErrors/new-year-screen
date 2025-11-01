package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/fogleman/gg"
)

const (
	width      = 1280
	height     = 720
	fontPath   = "assets/IBMPlexSans-Bold.ttf"
	bgImage    = "assets/background.png"
	newYearImg = "assets/newyear.png"
)

func main() {
	// Установим целевую дату (Новый год)
	target := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.Local)

	// Загружаем картинки
	bg := loadImage(bgImage)
	ny := loadImage(newYearImg)

	const fontSize = 72
	dc := gg.NewContext(width, height)
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Fatalf("failed to load font: %v", err)
	}

	for {
		now := time.Now()
		remaining := target.Sub(now)

		var img image.Image
		if remaining <= 0 {
			img = ny
		} else {
			img = renderCountdown(dc, bg, remaining)
		}

		var buf bytes.Buffer
		png.Encode(&buf, img)
		os.Stdout.Write(buf.Bytes())

		time.Sleep(time.Second)
	}
}

func loadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open image %s: %v", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}
	return img
}

func renderCountdown(dc *gg.Context, bg image.Image, remaining time.Duration) image.Image {
	days := int(remaining.Hours() / 24)
	hours := int(remaining.Hours()) % 24
	minutes := int(remaining.Minutes()) % 60

	text := formatCountdown(days, hours, minutes)

	dc.Clear()
	dc.DrawImage(bg, 0, 0)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
	return dc.Image()
}

func formatCountdown(d, h, m int) string {
	return fmt.Sprintf("%02d дней %02d часов %02d минут", d, h, m)
}
