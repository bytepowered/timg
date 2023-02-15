package pkg

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
)

var (
	Black = color.RGBA{A: 255}
)

type FontOption struct {
	Font     string
	FontSize float64
	Spacing  float64
	DPI      float64
}

func LoadFont(path string) (fnt *truetype.Font, err error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read font file error: %s", err)
	}
	font, err := freetype.ParseFont(bytes)
	if err != nil {
		return nil, fmt.Errorf("parse font error: %s", err)
	}
	return font, nil
}

func LoadBackgroundImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image file error: %s", err)
	}
	defer file.Close()
	if img, err = jpeg.Decode(file); err != nil {
		return nil, fmt.Errorf("decode image file error: %s", err)
	} else {
		return img, nil
	}
}

func LoadBackgroundImageRGBA(path string) (bg *image.RGBA, err error) {
	img, err := LoadBackgroundImage(path)
	if err != nil {
		return nil, err
	}
	background := image.NewRGBA(img.Bounds())
	draw.Draw(background, background.Bounds(), img, image.Point{}, draw.Src)
	return background, nil
}

//func DrawText(img *image.RGBA, text string, font *truetype.Font, opt FontOption) error {
//	c := freetype.NewContext()
//	c.SetDPI(opt.DPI)
//	c.SetFont(font)
//	c.SetFontSize(opt.FontSize)
//	c.SetClip(img.Bounds())
//	c.SetDst(img)
//	c.SetSrc(image.NewUniform(Black))
//	pt := freetype.Pt(0, 0)
//	_, err := c.DrawText(text, pt)
//	if err != nil {
//		return fmt.Errorf("draw string error: %s", err)
//	}
//	return nil
//}
