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
	c.drawTextLines(opts, strings.Split(text, "\n"),
		func(drawer *font.Drawer, lines []string) []string {
			output := make([]string, 0, len(lines))
			for _, line := range lines {
				output = c.cutText(drawer, line, output)
			}
			return output
		})
}

func (c *Canvas) DrawLines(opts text.Option, text []string) {
	c.drawTextLines(opts, text,
		func(drawer *font.Drawer, lines []string) []string {
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
		bound := c._canvas.Bounds()
		if bound.Dy() < drawer.Dot.Y.Ceil() {
			canvas := image.NewRGBA(image.Rect(0, 0, bound.Dx(), bound.Dy()+100))
			draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: pkg.White}, image.Point{}, draw.Src)
			draw.Draw(canvas, bound, c._canvas, image.Point{}, draw.Src)
			c._canvas = canvas
			drawer.Dst = c._canvas
		}
	}
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

func measureTextHeight(opts text.Option) int {
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
