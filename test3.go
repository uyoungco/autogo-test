package main

import (
	"github.com/Dasongzi1366/AutoGo/images"
)

func main99() {
	var colors = "bec0c0-101010,1,8,bec0c0-101010,-6,7,bbbdbe-101010,7,7,b0b3b4-101010,-25,-2,aeb1b2-101010,-25,3,aeb1b2-101010,13,-4,bcbebf-101010,13,1,b6b8b9-101010,19,-1,bec0c0-101010,19,8,979a9b-101010"
	x, y := images.FindMultiColors(0, 0, 0, 0, colors, 0.9, 0, 0)
	println("多点找色", x, y)

	//x1, y1, 出售按钮 := core.OpenCV.WaitFor(0, 0, 0, 0, "img/仓库_物品出售图标.png", true, 1, 0.8, 2*time.Second, 10, "仓库_物品出售图标")
	//println("出售按钮", x1, y1)
	//if !出售按钮 {
	//	println("未找到出售按钮")
	//
	//}
	//x1, y1, 出售按钮2 := core.OpenCV.WaitFor(0, 0, 0, 0, "img/仓库_物品出售图标2.png", true, 1, 0.8, 2*time.Second, 10, "仓库_物品出售图标")
	//println("出售按钮2", x1, y1)
	//if !出售按钮2 {
	//	println("未找到出售按钮2")
	//
	//}
}
