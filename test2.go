package main

import (
	"app/util"
	"fmt"
	"log"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/opencv"
)

// 预定义的颜色范围配置
var defaultColorRanges1 = []util.ColorRange{
	{
		Name: "黄色",
		HSV: util.HSVRange{
			LowerH: 11, LowerS: 45, LowerV: 58,
			UpperH: 14, UpperS: 66, UpperV: 77,
		},
		Color: "yellow",
	},
	{
		Name: "绿色",
		HSV: util.HSVRange{
			LowerH: 84, LowerS: 58, LowerV: 40,
			UpperH: 88, UpperS: 91, UpperV: 63,
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

// colorRGB 构造 BGR/Scalar（opencv 使用 BGR 排列）
func colorRGB(r, g, b uint8) opencv.Scalar {
	return opencv.NewScalar(float64(b), float64(g), float64(r), 0)
}

func main8() {
	// 截取屏幕区域
	img := images.CaptureScreen(570, 78, 1201, 638, 0)

	// 使用默认颜色范围查找物品
	boxes, err := util.FindByColor(img, defaultColorRanges1, 65, 68)
	if err != nil {
		log.Fatalf("查找失败: %v", err)
	}

	// 打印找到的物品
	for _, it := range boxes {
		fmt.Println(it)
	}

	// 可视化结果
	vis, err := util.VisualizeResults(img, boxes)
	if err != nil {
		log.Fatalf("可视化失败: %v", err)
	}
	vis_img := util.MatToImage(vis)
	images.Save(vis_img, "/sdcard/saved.png", 100)
	defer vis.Close()

	//win := opencv.NewWindow("image1")
	//defer win.Close()
	//win.IMShow(vis)
	//opencv.WaitKey(0)
}
