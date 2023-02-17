package main

import (
	"image/jpeg"
	"os"
	"wxbot-shenbi/canvas"
	"wxbot-shenbi/text"
)

func main() {
	// 创建一个白色背景的矩形图像
	c, err := canvas.NewCanvas(1032, 200)
	if err != nil {
		panic(err)
	}
	//var lines = []string{
	//	`怒发冲冠，凭栏处，潇潇雨歇。`,
	//	`抬望眼，仰天长啸，壮怀激烈。`,
	//	`三十功名尘与土，八千里路云和月。`,
	//	`莫等闲，白了少年头，空悲切。`,
	//	`靖康耻，犹未雪；臣子恨，何时灭！`,
	//	`驾长车踏破贺兰山缺。`,
	//	`壮志饥餐胡虏肉，笑谈渴饮匈奴血。`,
	//	`待从头，收拾旧山河，朝天阙。`,
	//	``,
	//	`关关雎鸠，在河之洲，窈窕淑女，君子好逑。`,
	//	`蒹葭苍苍，白露为霜。所谓伊人，在水一方。`,
	//}
	//c.DrawLines(text.DefaultOption, lines)
	var str = `
科技公司巨头YANDEX，全公司源代码被恶意泄漏，总共45GB，全世界黑客正在疯狂下载。本文末尾附下载地址。

YANDEX IS A RUSSIAN INTERNET COMPANY BASED IN MOSCOW

本次源代码泄漏，涉及公司所有产品，规模是前所未有的，以往同类事件都是某公司的部分产品和项目。根据技术分析，本次数据泄漏是YANDEX前员工在去年窃取公司全部源代码，在网上公开发布BT磁力链下载地址，其作恶动机与政治有关——暗示该员工与乌克兰军事方面存在联系。在前美国总统特朗普的批准下，美国国防部建立一支王牌网络作战部队，不排除本次事件与之有关。

1、乌克兰网络部队入侵俄罗斯中央银行，泄漏大量敏感数据。
2、境外黑客控制网约车，造成俄罗斯首都交通严重拥堵。
3、最大银行被黑客攻击，普京称这是一场“信息空间的战争”。

根据BT下载原理，下载的人越多，下载速度越快。当前下载速度达到惊人的5MB/s，本人上网冲浪10多个，从来没见过BT下载速度如此高速，可见全世界有多少怀着各种各样目的和意图的黑客、程序员和其他机构目前正在极度疯狂地下载这份源码！
`
	c.DrawText(text.DefaultOption, str)
	// 将图像输出为PNG格式的文件
	file, err := os.Create("output.jpeg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jpeg.Encode(file, c.Canvas(), &jpeg.Options{Quality: 50})
}
