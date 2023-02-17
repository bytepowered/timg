package canvas

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
	"wxbot-shenbi/pkg"
	"wxbot-shenbi/text"
)

var (
	defaultPadding = Padding{
		Left:   10,
		Right:  10,
		Top:    10,
		Bottom: 10,
	}
)

type Canvas struct {
	_fontDpi float64
	_font    *truetype.Font
	_canvas  *image.RGBA
	_padding Padding
}

func NewCanvas(width, height int) (*Canvas, error) {
	rect := image.Rect(0, 0, width, height)
	canvas := image.NewRGBA(rect)
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: pkg.White}, image.Point{}, draw.Src)
	// text
	if _font, err := text.Load("./resources/msyh.ttf"); err != nil {
		return nil, err
	} else {
		return &Canvas{_canvas: canvas, _font: _font, _padding: defaultPadding}, nil
	}
}

func (c *Canvas) DrawText(opts text.Option, text string) {
	c.drawTextLines(opts, strings.Split(text, "\n"), func(drawer *font.Drawer, inLines []string) []string {
		charWidth := drawer.MeasureString("中").Ceil()
		output := make([]string, 0, len(inLines))
		// 检查每行文字的长度，如果超过了画布的宽度，则进行换行
		cutLine := func(line string) {
			// line: 一行文字
			overflow := drawer.MeasureString(line).Ceil() - c.ContentWidth()
			if overflow > charWidth {
				output = append(output, line)
				output = append(output, "====================")
			} else {
				output = append(output, line)
			}
		}
		for _, line := range inLines {
			cutLine(line)
		}
		return output
	})
}

func (c *Canvas) DrawLines(opts text.Option, text []string) {
	c.drawTextLines(opts, text, func(drawer *font.Drawer, lines []string) []string {
		return lines
	})
}

func (c *Canvas) drawTextLines(opts text.Option, lines []string, transform func(*font.Drawer, []string) []string) {
	drawer := &font.Drawer{
		Dst: c._canvas,
		Src: image.NewUniform(opts.Color),
		Face: truetype.NewFace(c._font, &truetype.Options{
			Size:    opts.Size,
			DPI:     opts.DPI,
			Hinting: font.HintingVertical,
		}),
	}
	// 绘制文字
	height := measureTextHeight(opts)
	drawer.Dot.Y = fixed.I(c._padding.Top)
	drawer.Dot.X = fixed.I(c._padding.Left)
	spacing := fixed.I(int(float64(height) * opts.Spacing))
	for _, str := range transform(drawer, lines) {
		drawer.DrawString(str)
		drawer.Dot.Y = drawer.Dot.Y + spacing
		drawer.Dot.X = fixed.I(c._padding.Left)
	}
}

func measureTextHeight(opts text.Option) int {
	return int(math.Ceil(opts.Size * opts.DPI / 72))
}

func (c *Canvas) DrawRect(pos Position, color color.Color) {
	rect := image.Rect(
		pos.X+c._padding.Left,
		pos.Y+c._padding.Top,
		pos.Width-(c._padding.Left+c._padding.Right),
		pos.Height-(c._padding.Top+c._padding.Bottom),
	)
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

func (c *Canvas) ContentWidth() int {
	return c.Width() - (c._padding.Left + c._padding.Right)
}

func (c *Canvas) ContentHeight() int {
	return c.Width() - (c._padding.Top + c._padding.Bottom)
}

type Position struct {
	X, Y          int
	Width, Height int
}

type Padding struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}
