package timg

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"math"
	"strings"
)

var (
	defaultPadding = Padding{
		Left:   10,
		Right:  10,
		Top:    10,
		Bottom: 0,
	}
)

type Canvas struct {
	debug    bool
	dpi      float64
	fontPath string
	font     *truetype.Font
	rgba     *image.RGBA
	padding  Padding
}

type CanvasOption func(canvas *Canvas)

func (c *Canvas) DrawText(fontOpts FontOption, text string) {
	_lines := strings.Split(text, "\n")
	lines := make([]string, 0, len(_lines)*2)
	fontFace := c.newFontFace(fontOpts)
	for _, line := range _lines {
		lines = c.cut(fontFace, line, lines)
	}
	c.draw(fontOpts, fontFace, lines)
}

func (c *Canvas) draw(fontOpts FontOption, fontFace font.Face, lines []string) {
	charHeight := measureTextHeight(fontOpts)
	lineHeight := int(float64(charHeight) * fontOpts.Spacing)
	if c.debug {
		fmt.Printf("[DEBUG] draw char height: %d, line height: %d, lines: %d \n", charHeight, lineHeight, len(lines))
	}
	// 确保图片高度大于文本行数高度
	if pageHeight := (len(lines) + 1) * lineHeight; c.Height() < pageHeight {
		c.resize(c.Width(), pageHeight)
		if c.debug {
			fmt.Printf("[DEBUG] draw resize canvas height: %d\n", pageHeight)
		}
	}
	if c.debug {
		for i := 0; i < len(lines)+1; i++ {
			y := i*lineHeight + c.padding.Top
			draw.Draw(c.rgba,
				image.Rect(0, y, c.Width(), y+1),
				&image.Uniform{C: NiceGray},
				image.Point{
					X: 0, Y: y,
				},
				draw.Src)
		}
	}
	drawer := &font.Drawer{
		Dst:  c.rgba,
		Src:  image.NewUniform(fontOpts.Color),
		Face: fontFace,
	}
	// 绘制文字
	drawer.Dot.Y = fixed.I(c.padding.Top + lineHeight)
	drawer.Dot.X = fixed.I(c.padding.Left)
	for _, line := range lines {
		drawer.DrawString(line)
		drawer.Dot.Y = drawer.Dot.Y + fixed.I(lineHeight)
		drawer.Dot.X = fixed.I(c.padding.Left)
	}
}

func (c *Canvas) resize(width, height int) {
	_new := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(_new, _new.Bounds(), &image.Uniform{C: White}, image.Point{}, draw.Src)
	if c.debug {
		draw.Draw(_new, image.Rect(0, height-1, width, height), &image.Uniform{C: NiceRed}, image.Point{
			X: 0, Y: height - 1,
		}, draw.Src)
	}
	c.rgba = _new
}

func (c *Canvas) cut(face font.Face, text string, output []string) []string {
	oneCharWidth := font.MeasureString(face, "中").Ceil()
	maxCharCount := c.ContentWidth() / oneCharWidth
	// 检查每行文本的长度，如果超过了画布的宽度，则进行换行
	overflow := font.MeasureString(face, text).Ceil() - c.ContentWidth()
	if overflow > oneCharWidth {
		line := []rune(text)
		// 溢出，需要按rune来裁剪到最大宽度
		index := min(len(line), maxCharCount)
		// 尝试裁剪到最大宽度
		for i := 0; i < maxCharCount; i++ {
			index += 1
			displayWidth := font.MeasureString(face, string(line[:index])).Ceil()
			if displayWidth >= c.ContentWidth() {
				index--
				break
			}
		}
		output = append(output, string(line[:index]))
		return c.cut(face, string(line[index:]), output)
	} else {
		output = append(output, text)
	}
	return output
}

func (c *Canvas) newFontFace(fontOpts FontOption) font.Face {
	return truetype.NewFace(c.font, &truetype.Options{
		Size:    fontOpts.Size,
		DPI:     fontOpts.DPI,
		Hinting: font.HintingVertical,
	})
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

func WithDebug(debugEnabled bool) CanvasOption {
	return func(canvas *Canvas) {
		canvas.debug = debugEnabled
	}
}

func WithDPI(dpi float64) CanvasOption {
	return func(canvas *Canvas) {
		canvas.dpi = dpi
	}
}

func WithFontPath(fpath string) CanvasOption {
	return func(canvas *Canvas) {
		canvas.fontPath = fpath
	}
}

func WithPadding(padding Padding) CanvasOption {
	return func(canvas *Canvas) {
		canvas.padding = padding
	}
}

func NewDefaultCanvas() (*Canvas, error) {
	return NewCanvas(1032, 100,
		WithPadding(defaultPadding),
		WithDPI(FontOptionDPI),
		WithDebug(false),
	)
}

func NewCanvas(width, height int, opts ...CanvasOption) (*Canvas, error) {
	box := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(box)
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{C: White}, image.Point{}, draw.Src)
	canvas := &Canvas{rgba: rgba, padding: defaultPadding}
	// options
	for _, opt := range opts {
		opt(canvas)
	}
	// init font
	fontPath := canvas.fontPath
	if fontPath == "" {
		fontPath = "./resources/msyh.ttf"
	}
	if _font, err := LoadFont(fontPath); err != nil {
		return nil, fmt.Errorf("load font failed: %w", err)
	} else {
		canvas.font = _font
	}
	return canvas, nil
}

func measureTextHeight(opts FontOption) int {
	return int(math.Ceil(opts.Size * opts.DPI / 72))
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
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
