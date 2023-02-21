package timg

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
)

var (
	defaultPadding = Padding{
		Left:   30,
		Right:  30,
		Top:    30,
		Bottom: 0,
	}
)

type Canvas struct {
	dpi     float64
	fpath   string
	font    *truetype.Font
	rgba    *image.RGBA
	padding Padding
}

func NewCanvas(width, height int) (*Canvas, error) {
	rect := image.Rect(0, 0, width, height)
	canvas := image.NewRGBA(rect)
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: White}, image.Point{}, draw.Src)
	// text
	if _font, err := Load("./resources/msyh.ttf"); err != nil {
		return nil, err
	} else {
		return &Canvas{rgba: canvas, font: _font, padding: defaultPadding}, nil
	}
}

func (c *Canvas) DrawText(opts FontOption, text string) {
	tf := func(drawer *font.Drawer, lines []string) []string {
		output := make([]string, 0, len(lines))
		for _, line := range lines {
			output = c.cutText(drawer, line, output)
		}
		return output
	}
	c.drawLines(opts, strings.Split(text, "\n"), tf)
}

func (c *Canvas) DrawLines(opts FontOption, text []string) {
	c.drawLines(opts, text, func(drawer *font.Drawer, lines []string) []string {
		return lines
	})
}

func (c *Canvas) drawLines(fontOpts FontOption, lines []string, transform func(*font.Drawer, []string) []string) {
	drawer := &font.Drawer{
		Dst: c.rgba,
		Src: image.NewUniform(fontOpts.Color),
		Face: truetype.NewFace(c.font, &truetype.Options{
			Size:    fontOpts.Size,
			DPI:     fontOpts.DPI,
			Hinting: font.HintingVertical,
		}),
	}
	// 绘制文字
	height := measureTextHeight(fontOpts)
	drawer.Dot.Y = fixed.I(c.padding.Top)
	drawer.Dot.X = fixed.I(c.padding.Left)
	spacing := fixed.I(int(float64(height) * fontOpts.Spacing))
	for _, str := range transform(drawer, lines) {
		drawer.DrawString(str)
		drawer.Dot.Y = drawer.Dot.Y + spacing
		drawer.Dot.X = fixed.I(c.padding.Left)
		bound := c.rgba.Bounds()
		if bound.Dy() < drawer.Dot.Y.Ceil() {
			srcimg := c.rgba
			c.resize(bound.Dx(), bound.Dy()+height*2)
			draw.Draw(c.rgba, bound, srcimg, image.Point{}, draw.Src)
			drawer.Dst = c.rgba
		}
	}
}

func (c *Canvas) resize(width, height int) {
	_new := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(_new, _new.Bounds(), &image.Uniform{C: White}, image.Point{}, draw.Src)
	c.rgba = _new
}

func (c *Canvas) cutText(drawer *font.Drawer, text string, output []string) []string {
	oneCharWidth := drawer.MeasureString("中").Ceil()
	maxCharCount := c.ContentWidth() / oneCharWidth
	// 检查每行文本的长度，如果超过了画布的宽度，则进行换行
	overflow := drawer.MeasureString(text).Ceil() - c.ContentWidth()
	if overflow > oneCharWidth {
		line := []rune(text)
		// 溢出，需要按rune来裁剪到最大宽度
		index := min(len(line), maxCharCount)
		// 尝试裁剪到最大宽度
		for i := 0; i < maxCharCount; i++ {
			index += 1
			displayWidth := drawer.MeasureString(string(line[:index])).Ceil()
			if displayWidth >= c.ContentWidth() {
				index--
				break
			}
		}
		output = append(output, string(line[:index]))
		return c.cutText(drawer, string(line[index:]), output)
	} else {
		output = append(output, text)
	}
	return output
}

func measureTextHeight(opts FontOption) int {
	return int(math.Ceil(opts.Size * opts.DPI / 72))
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func (c *Canvas) DrawRect(pos Position, color color.Color) {
	rect := image.Rect(
		pos.X+c.padding.Left,
		pos.Y+c.padding.Top,
		pos.Width-(c.padding.Left+c.padding.Right),
		pos.Height-(c.padding.Top+c.padding.Bottom),
	)
	draw.Draw(c.rgba, rect, &image.Uniform{C: color}, image.Point{}, draw.Src)
}

func (c *Canvas) Canvas() *image.RGBA {
	return c.rgba
}

func (c *Canvas) Width() int {
	return c.rgba.Bounds().Dx()
}

func (c *Canvas) Height() int {
	return c.rgba.Bounds().Dy()
}

func (c *Canvas) ContentWidth() int {
	return c.Width() - (c.padding.Left + c.padding.Right)
}

func (c *Canvas) ContentHeight() int {
	return c.Width() - (c.padding.Top + c.padding.Bottom)
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
