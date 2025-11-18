package main

import (
	"app/core"
	"app/scene"
	"app/util"
	"fmt"
	"log"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
)

// 预定义的颜色范围配置
var defaultColorRanges = []util.ColorRange{
	{
		Name: "黄色",
		HSV: util.HSVRange{
			LowerH: 11, LowerS: 45, LowerV: 58,
			UpperH: 14, UpperS: 66, UpperV: 77,
		},
		Color: "yellow",
	},
	{
		Name: "黄色",
		HSV: util.HSVRange{
			LowerH: 11, LowerS: 45, LowerV: 50,
			UpperH: 15, UpperS: 90, UpperV: 81,
		},
		Color: "yellow",
	},
	{
		Name: "绿色",
		HSV: util.HSVRange{
			LowerH: 83, LowerS: 52, LowerV: 49,
			UpperH: 88, UpperS: 68, UpperV: 65,
		},
		Color: "green",
	},
	{
		Name: "红色",
		HSV: util.HSVRange{
			LowerH: 174, LowerS: 72, LowerV: 61,
			UpperH: 178, UpperS: 109, UpperV: 86,
		},
		Color: "red",
	},
	{
		Name: "紫色",
		HSV: util.HSVRange{
			LowerH: 116, LowerS: 45, LowerV: 49,
			UpperH: 124, UpperS: 77, UpperV: 81,
		},
		Color: "purple",
	},
}

func main() {
	// 判断场景值执行获取仓库物品价值
	for {
		sceneName := scene.Identify()
		core.Log("当前场景: " + sceneName)
		if sceneName == "主界面" || sceneName == "主界面2" {
			// 点击仓库
			// core.RandomClickInArea(23, 647, 133, 682)
			core.Click(74, 665)
			core.Sleep(2000)
		}
		if sceneName == "界面仓库" {
			break
		}
		//core.Sleep(2000)
	}
	// 随机延迟
	// core.RandomSleep(2000, 3000)
	matched := images.CmpColor(45, 427, "9EA2A7-101010|A3A3A8-101010|9FA1A5-101010", 0.9, 0)

	// var x1, y1 = core.Color.WaitForColor(32, 413, 58, 430, "9DA4A9|A3A3A8|9CA1A5", 0.9, 10)
	println("判断左侧图标是否被点击", matched)
	if matched {
		// 点击保险箱图标
		core.Click(43, 421)
		core.Sleep(500)
	}

	// 计算安全箱范围
	// x1固定开始: 189
	// y1需要计算: 通过安全箱问号图标获取y值 y1 = y1 - 37
	// x2, y2 400 639
	// opencv检测安全箱物品
	_, y1, 检查图标 := core.OpenCV.WaitFor(83, 212, 332, 644, "img/仓库_安全箱_问号.png", true, 1, 0.8, 2*time.Second, 10, "仓库_物品出售图标")
	if !检查图标 {
		println("检测到仓库_安全箱_问号")
		return
	}
	y1 = y1 + 37
	var rect = util.Rect{X1: 189, Y1: y1, X2: 400, Y2: 639}
	img := images.CaptureScreen(rect.X1, rect.Y1, rect.X2, rect.Y2, 0)
	// images.Save(img, "/sdcard/保险箱截图.png", 100)

	// 使用默认颜色范围查找物品
	boxes, err := util.FindByColor(img, defaultColorRanges, 65, 68)
	if err != nil {
		log.Fatalf("查找失败: %v", err)
	}
	if len(boxes) == 0 {
		println("未找到仓库物品")
		return
	}

	vis, err := util.VisualizeResults(img, boxes)
	if err != nil {
		log.Fatalf("可视化失败: %v", err)
	}
	vis_img := util.MatToImage(vis)

	images.Save(vis_img, "/sdcard/saved.png", 100)
	defer vis.Close()

	// 将相对坐标转换为全局坐标
	for idx := range boxes {
		// 转换左上角坐标
		boxes[idx].X, boxes[idx].Y = util.ToGlobalCoords(rect, boxes[idx].X, boxes[idx].Y)
		// 转换中心点坐标
		boxes[idx].CenterX, boxes[idx].CenterY = util.ToGlobalCoords(rect, boxes[idx].CenterX, boxes[idx].CenterY)
	}

	// 打印转换后的全局坐标
	for _, it := range boxes {
		fmt.Println(it)
	}

	// 循环出售物品
	for idx, _ := range boxes {
		item := boxes[idx]
		// 点击物品
		core.Click(item.CenterX, item.CenterY)
		core.Sleep(1000)
		// 检查出售按钮
		x1, y1, 出售按钮 := core.OpenCV.WaitFor(76, 71, 1206, 649, "img/仓库_出售文本.png", true, 1, 0.7, 2*time.Second, 10, "仓库_出售文本")
		println("出售按钮", x1, y1)
		if !出售按钮 {
			println("未找到出售按钮", x1, y1)
			return
		}
		core.Click(x1, y1)
		// 检查出售弹框是否打开
		var colors = "939699-101010,0,9,939699-101010,1,23,939699-101010,-11,23,93979a-101010,11,22,94979a-101010,35,2,939699-101010,34,12,94979b-101010,23,25,939699-101010,43,24,939699-101010"
		_, _, 出售弹框 := core.Color.WaitFindMultiColors(89, 65, 153, 100, colors, 0.8, 10)
		if !出售弹框 {
			println("未找到出售弹框")
			return
		}
		// 检查交易行按钮颜色如果是灰色代表不能上架直接出售
		// 为真表示不能上架
		var 是否可以上架 = core.Color.CmpColor(1092, 533, "626769", 0.9)
		if 是否可以上架 {
			println("不能上架直接出售")
			// 点击出售跳出循环
			core.Click(1083, 294)
			core.Sleep(1000)
			continue
			// return
		}
		// ocr获取系统价格和交易行价格
		// 系统
		systemPriceOcr := core.OCR.DetectAllText(992, 214, 1168, 245)
		systemPriceInt, _ := util.StringToIntMap(systemPriceOcr[0])
		println("systemPriceInt", systemPriceInt)

		// 交易行
		marketPriceOcr := core.OCR.DetectAllText(980, 444, 1171, 480)
		marketPriceInt, _ := util.StringToIntMap(marketPriceOcr[0])
		println("marketPriceInt", marketPriceInt)
		println("差值", marketPriceInt-systemPriceInt)
		if (marketPriceInt - systemPriceInt) < 4000 {
			println("小于4000直接出售")
			// 小于4000直接出售
			core.Click(1083, 294)
			core.Sleep(1000)
			continue
			// return
		}
		// 点击上架打开交易行窗口
		core.Click(1084, 529)
		core.Sleep(1000)
		var colors2 = "939699-101010,-33,-16,939699-101010,-33,6,939699-101010,-9,-16,939699-101010,24,-13,939699-101010,25,-2,939699-101010,45,-14,95989b-101010,40,7,939699-101010,63,-17,939699-101010,61,6,939699-101010,68,5,979a9d-101010"
		_, _, 交易行弹框 := core.Color.WaitFindMultiColors(82, 64, 223, 101, colors2, 0.8, 10)
		if !交易行弹框 {
			println("未找到交易行弹框")
			return
		}
		core.Click(810, 473)
		core.Sleep(500)
		core.Click(991, 577)
		core.Sleep(1000)
		println("正常出售完成")
		//fmt.Println(it)
		//return
	}

}
