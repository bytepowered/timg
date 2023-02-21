package timg

import (
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

func (c *Canvas) DrawText(fontOpts FontOption, text string) {
	rawlines := strings.Split(text, "\n")
	lines := make([]string, 0, len(rawlines))
	fontFace := c.newFontFace(fontOpts)
	for _, line := range rawlines {
		lines = c.cut(fontFace, line, lines)
	}
	c.draw(fontOpts, fontFace, lines)
}

func (c *Canvas) DrawLines(fontOpts FontOption, lines []string) {
	fontFace := c.newFontFace(fontOpts)
	c.draw(fontOpts, fontFace, lines)
}

func (c *Canvas) draw(fontOpts FontOption, fontFace font.Face, lines []string) {
	charHeight := measureTextHeight(fontOpts)
	lineHeight := int(float64(charHeight) * fontOpts.Spacing)
	// 确保图片高度大于文本行数高度
	if c.rgba.Bounds().Dy() < len(lines)*lineHeight {
		c.resize(c.Width(), len(lines)*lineHeight)
	}
	drawer := &font.Drawer{
		Dst:  c.rgba,
		Src:  image.NewUniform(fontOpts.Color),
		Face: fontFace,
	}
	// 绘制文字
	drawer.Dot.Y = fixed.I(c.padding.Top)
	drawer.Dot.X = fixed.I(c.padding.Left)
	for _, str := range lines {
		drawer.DrawString(str)
		drawer.Dot.Y = drawer.Dot.Y + fixed.I(lineHeight)
		drawer.Dot.X = fixed.I(c.padding.Left)
		bound := c.rgba.Bounds()
		if bound.Dy() < drawer.Dot.Y.Ceil() {
			srcimg := c.rgba
			c.resize(bound.Dx(), bound.Dy()+lineHeight)
			draw.Draw(c.rgba, bound, srcimg, image.Point{}, draw.Src)
			drawer.Dst = c.rgba
		}
	}
}

func (c *Canvas) resize(width, height int) {
	_new := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(_new, _new.Bounds(), &image.Uniform{C: White}, image.Point{}, draw.Src)
	draw.Draw(_new, image.Rect(0, height-1, width, height), &image.Uniform{C: NiceRed}, image.Point{
		X: 0, Y: height - 1,
	}, draw.Src)
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
