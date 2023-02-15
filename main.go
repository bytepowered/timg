//package main
//
//import (
//	"fmt"
//	"image"
//	"image/color"
//	"image/draw"
//	"image/jpeg"
//	"os"
//	"wxbot-fuzi/pkg"
//
//	"github.com/golang/freetype"
//)
//
//var (
//	utf8FontFile = "./fonts/micr.ttf"
//	utf8FontSize = float64(14.0)
//	spacing      = float64(1.0)
//	dpi          = float64(360)
//	black        = color.RGBA{0, 0, 0, 255}
//)
//
//func main() {
//	var text = []string{
//		`怒发冲冠，凭栏处，潇潇雨歇。`,
//		`抬望眼，仰天长啸，壮怀激烈。`,
//		`三十功名尘与土，八千里路云和月。`,
//		`莫等闲，白了少年头，空悲切。`,
//		`靖康耻，犹未雪；臣子恨，何时灭！`,
//		`驾长车踏破贺兰山缺。`,
//		`壮志饥餐胡虏肉，笑谈渴饮匈奴血。`,
//		`待从头，收拾旧山河，朝天阙。`,
//		``,
//		`关关雎鸠，在河之洲，窈窕淑女，君子好逑。`,
//		`蒹葭苍苍，白露为霜。所谓伊人，在水一方。`,
//	}
//
//	font, err := pkg.LoadFont("./fonts/micr.ttf")
//	if err != nil {
//		panic(err)
//	}
//
//	canvas := image.NewRGBA(image.Rectangle{})
//	draw.Draw(canvas, canvas.Bounds(), img, image.Point{}, draw.Src)
//	//canvas, err := pkg.LoadBackgroundImageRGBA("./background.jpeg")
//	//if err != nil {
//	//	panic(err)
//	//}
//	fontColor := image.NewUniform(black)
//
//	// draw footer to canvas
//	//footer :=
//
//	ctx := freetype.NewContext()
//	ctx.SetDPI(dpi) //screen resolution in Dots Per Inch
//	ctx.SetFont(font)
//	ctx.SetFontSize(utf8FontSize) //font size in points
//	ctx.SetClip(canvas.Bounds())
//	ctx.SetDst(canvas)
//	ctx.SetSrc(fontColor)
//
//	// 绘制文字
//	pt := freetype.Pt(100, 100+int(ctx.PointToFixed(utf8FontSize)>>6))
//	for _, str := range text {
//		_, err := ctx.DrawText(str, pt)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		pt.Y += ctx.PointToFixed(utf8FontSize * spacing)
//	}
//
//	// 将绘制好的图片保存为JPEG文件
//	out, err := os.Create("output.jpg")
//	if err != nil {
//		panic(err)
//	}
//	defer out.Close()
//
//	jpeg.Encode(out, canvas, &jpeg.Options{Quality: 80})
//}

package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

var (
	Black = color.RGBA{A: 255}
	White = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Blue  = color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

func main() {
	// 创建一个白色背景的矩形图像
	canvas, err := NewCanvas(Size{
		Width:  1032,
		Height: 2000,
	})
	if err != nil {
		panic(err)
	}
	red := color.RGBA{R: 204, G: 68, B: 60, A: 255}
	canvas.DrawRect(Position{
		Width: canvas.Width(), Height: 300,
	}, red)

	_ = canvas.DrawText(Position{
		Width:  canvas.Width(),
		Height: 500,
	}, FontOptionOf(24, White), "我拿水兑水")
	// 在矩形图像的底部画一个蓝色小矩形
	// DrawFooterBG(canvas, 200)

	// 将图像输出为PNG格式的文件
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, canvas.Canvas())
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

type Canvas struct {
	_fontDpi float64
	_font    *truetype.Font
	_canvas  *image.RGBA
}

func NewCanvas(size Size) (*Canvas, error) {
	// canvas
	rect := image.Rect(0, 0, size.Width, size.Height)
	canvas := image.NewRGBA(rect)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)
	// font
	font, err := LoadFont("./fonts/micr.ttf")
	if err != nil {
		return nil, err
	}
	return &Canvas{_canvas: canvas, _font: font}, nil
}

func (c *Canvas) DrawText(pos Position, opts FontOption, text string) error {
	face := truetype.NewFace(c._font, &truetype.Options{
		Size:    opts.Size,
		DPI:     opts.DPI,
		Hinting: font.HintingFull,
	})
	drawer := &font.Drawer{
		Dst:  c._canvas,
		Src:  image.NewUniform(opts.Color),
		Face: face,
	}
	// centerX := X: (fixed.I(c.Width()) - drawer.MeasureString(text)) / 2,
	// dy := int(math.Ceil(opts.Size * opts.Spacing * opts.DPI / 72))
	// 居中对齐
	charHeight := (face.Metrics().Ascent + face.Metrics().Descent).Ceil()
	drawer.Dot = fixed.Point26_6{
		//X: (fixed.I(c.Width()) - drawer.MeasureString(text)) / 2,
		Y: fixed.I(charHeight),
	}
	drawer.DrawString(text)
	return nil
}

func measureTextHeight(opts FontOption) int {
	return int(math.Ceil(opts.Size * opts.DPI / 72))
}

func (c *Canvas) DrawRect(pos Position, color color.Color) {
	rect := image.Rect(pos.X, pos.Y, pos.Width, pos.Height)
	draw.Draw(c._canvas, rect, &image.Uniform{C: color}, image.Point{}, draw.Src)
}

func (c *Canvas) Canvas() *image.RGBA {
	return c._canvas
}

func (c *Canvas) Width() int {
	return c._canvas.Bounds().Dx()
}

func (c *Canvas) Height() int {
	return c._canvas.Bounds().Dy()
}

type FontOption struct {
	Size    float64
	Spacing float64
	DPI     float64
	Color   color.Color
}

func DefaultFontOption() FontOption {
	return FontOption{
		Size:  12,
		DPI:   600,
		Color: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

func FontOptionOf(size float64, color color.Color) FontOption {
	return FontOption{
		Size:  size,
		DPI:   600,
		Color: color,
	}
}

type Position struct {
	X, Y          int
	Width, Height int
}

type Size struct {
	Width, Height int
}
