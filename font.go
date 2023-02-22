package timg

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image/color"
	"os"
)

type FontOption struct {
	Size    float64
	Spacing float64
	DPI     float64
	Color   color.Color
}

const (
	FontOptionDPI = 180
)

var (
	FontOptionDefault = FontOption{
		Size:    12,
		DPI:     FontOptionDPI,
		Spacing: 1.5,
		Color:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
)

func FontOptionOf(size float64, color color.Color) FontOption {
	return FontOption{
		Size:  size,
		DPI:   FontOptionDPI,
		Color: color,
	}
}

func LoadFont(path string) (fnt *truetype.Font, err error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read text file error: %s", err)
	}
	font, err := freetype.ParseFont(bytes)
	if err != nil {
		return nil, fmt.Errorf("parse text error: %s", err)
	}
	return font, nil
}
