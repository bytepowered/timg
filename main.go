package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"wxbot-fuzi/pkg"

	"github.com/golang/freetype"
)

var (
	utf8FontFile = "./fonts/micr.ttf"
	utf8FontSize = float64(14.0)
	spacing      = float64(1.0)
	dpi          = float64(360)
	black        = color.RGBA{0, 0, 0, 255}
)

func main() {
	var text = []string{
		`怒发冲冠，凭栏处，潇潇雨歇。`,
		`抬望眼，仰天长啸，壮怀激烈。`,
		`三十功名尘与土，八千里路云和月。`,
		`莫等闲，白了少年头，空悲切。`,
		`靖康耻，犹未雪；臣子恨，何时灭！`,
		`驾长车踏破贺兰山缺。`,
		`壮志饥餐胡虏肉，笑谈渴饮匈奴血。`,
		`待从头，收拾旧山河，朝天阙。`,
		``,
		`关关雎鸠，在河之洲，窈窕淑女，君子好逑。`,
		`蒹葭苍苍，白露为霜。所谓伊人，在水一方。`,
	}

	font, err := pkg.LoadFont("./fonts/micr.ttf")
	if err != nil {
		panic(err)
	}
	background, err := pkg.LoadBackgroundImageRGBA("./background.jpeg")
	if err != nil {
		panic(err)
	}
	fontForeGroundColor := image.NewUniform(black)

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi) //screen resolution in Dots Per Inch
	ctx.SetFont(font)
	ctx.SetFontSize(utf8FontSize) //font size in points
	ctx.SetClip(background.Bounds())
	ctx.SetDst(background)
	ctx.SetSrc(fontForeGroundColor)

	// 绘制文字
	pt := freetype.Pt(100, 100+int(ctx.PointToFixed(utf8FontSize)>>6))
	for _, str := range text {
		_, err := ctx.DrawString(str, pt)
		if err != nil {
			fmt.Println(err)
			return
		}
		pt.Y += ctx.PointToFixed(utf8FontSize * spacing)
	}

	// 将绘制好的图片保存为JPEG文件
	out, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	jpeg.Encode(out, background, &jpeg.Options{Quality: 80})
}
