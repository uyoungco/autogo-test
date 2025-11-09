package core

import (
	"app/assets"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/console"
	"github.com/Dasongzi1366/AutoGo/imgui"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/utils"
)

func Log(message string) {
	// 将日志推送到SLS
	// tool.SLS_Log(message)

	// 打印到控制台
	fmt.Println(message)

	// 打印到imgui界面
	//imgui.Toast(message)
}

func Print(message string) {
	console.Println(message)
}

func Click(x, y int) {
	motion.Click(x, y, 1, 0)
}

// 在指定区域内随机点击
func RandomClickInArea(x1, y1, x2, y2 int) {
	// 确保坐标顺序正确
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// 在区域内生成随机坐标
	randomX := utils.Random(x1, x2)
	randomY := utils.Random(y1, y2)

	// 点击随机位置
	Click(randomX, randomY)
}

func Swipe(x1, y1, x2, y2, duration int) {
	motion.Swipe(x1, y1, x2, y2, duration, 0, 0)
}

func Sleep(ms int) {
	utils.Sleep(ms)
}

// 毫秒
func RandomSleep(min, max int) {
	randNum := utils.Random(min, max)
	Sleep(randNum)
}

func Toast(context string) {
	imgui.Toast(context)
	fmt.Println(context)
}

func ClickMap() {
	RandomClickInArea(60, 20, 160, 35)
}

// 关闭所有框
func CloseAllWindows() {
	// 获取 img/close 目录下的所有文件
	closeFiles, err := assets.ImageFile.ReadDir("img/close")
	if err != nil {
		return
	}

	// 循环检测最多4次
	for i := 0; i < 4; i++ {
		startTime := time.Now() // 记录开始时间
		found := false

		for _, file := range closeFiles {
			fileName := file.Name()
			// 只处理png图片文件
			if strings.HasSuffix(strings.ToLower(fileName), ".png") {

				// 构建相对路径
				imgPath := "img/close/" + fileName
				// 使用 OpenCV 检测图片，检测范围是全屏 (0,0,0,0)
				x, y := OpenCV.FindImage(0, 0, 0, 0, imgPath, false, 1, 0.8)

				if x != -1 && y != -1 {
					// 检测到图片，点击识别到的位置 (x+5, y+5)
					Click(x+5, y+5)
					found = true
					Sleep(500) // 点击后稍等一下，让界面反应
					break      // 跳出内层循环，继续下一次检测
				}
			}
		}

		// 计算本次循环耗时
		duration := time.Since(startTime)
		fmt.Println("第" + strconv.Itoa(i+1) + "次检测耗时: " + duration.String())

		// 如果本次没检测到任何关闭按钮，则退出循环
		if !found {
			break
		}
	}

}
