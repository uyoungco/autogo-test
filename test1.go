package main

import (
	"app/core"
	"app/util"
	"fmt"
	"time"
)

// 寻找名称使用ocr
func main3() {
	//template, _ := assets.ImageFile.ReadFile("img/gerenzhongxin.png")
	//img, _, _ := image.Decode(bytes.NewReader(template))
	//img2 := core.ImageToNRGBA(img)
	//
	//clippedImg := images.Clip(img2, 957, 588, 1138, 609)
	//images.Save(clippedImg, "/sdcard/autogotemp/clippedImg.png", 100)
	var x1, y1 = core.OpenCV.FindImage(956, 585, 1258, 611, "img/editicon.png", false, 1, 0.8)
	fmt.Println("图像1位置: ", x1, y1)
	//img := images.ReadFromPath("/sdcard/autogotemp/clippedImg.png")
	// results := ppocr.Ocr(958, 587, x1, 610, "EAEBEB-101010", 0)
	results := core.OCR.DetectAllText(958, 587, x1, 610)
	//img := images.CaptureScreen(958, 587, x1, 610, 0)
	//images.Save(img, "/sdcard/autogotemp/ocr-test.png", 100)

	fmt.Print(results)
}

// 查找哈弗币
func main4() {

	//img := images.ReadFromPath("/sdcard/autogotemp/clippedImg.png")
	// results := ppocr.Ocr(958, 587, x1, 610, "EAEBEB-101010", 0)
	results := core.OCR.DetectAllText(1088, 19, 1170, 43)
	//img := images.CaptureScreen(958, 587, x1, 610, 0)
	//images.Save(img, "/sdcard/autogotemp/ocr-test.png", 100)

	fmt.Print(results)
}

// 查找编号
func main5() {
	results := core.OCR.DetectAllText(990, 673, 1180, 692)
	fmt.Print("查找编号", results)
}

func main6() {
	fmt.Print("开始执行获取哈夫币")
	// 获取哈夫币图标位置
	x1, y1, _ := core.OpenCV.WaitFor(947, 4, 1279, 44, "img/仓库-哈夫币-icon.png", false, 1, 0.8, 2*time.Second, 60, "等待获取哈夫币图标")
	x2, _ := core.OpenCV.FindImage(947, 4, 1279, 44, "img/仓库-三角币-icon.png", false, 1, 0.8)

	x3 := x1 + 22
	y3 := y1

	x4 := x2
	y4 := 42

	results := core.OCR.DetectAllText(x3, y3, x4, y4)
	fmt.Print("查找编号", results)
}

func main7() {
	data := map[string]interface{}{
		"deviceCode":      "003",
		"accountId":       "186009030988940728143",
		"gameName":        "8苗条的d子弹",
		"havalCoinAmount": "87,725K",
	}
	fmt.Print("data", data, "\n")

	code, data2, err := util.HttpRequest.PostJSON(
		"http://175.24.153.109:4000/device/update",
		data,
		nil,
		1000*30,
	)

	fmt.Print(code, "\n", data2, "\n", err, "\n")
	if err != nil {
		fmt.Print(err)
	}

}
