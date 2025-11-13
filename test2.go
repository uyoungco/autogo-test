package main

import (
	"app/util"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/opencv"
)

// GameItem 对应 Python 的 GameItem 类
type GameItem struct {
	X, Y  int
	W, H  int
	Size  int
	Color string
}

func (gi GameItem) String() string {
	return fmt.Sprintf("GameItem(x:%d, y:%d, w:%d, h:%d, %s, 占用:%d格)", gi.X, gi.Y, gi.W, gi.H, gi.Color, gi.Size)
}

// createCellArray 根据 cellSize 生成面积区间（默认容差 5%，长度 9）
func createCellArray(cellSize, length int, tolerance float64) [][2]float64 {
	base := float64(cellSize * cellSize) // 110*110 = 12100
	res := make([][2]float64, 0, length)
	for i := 1; i <= length; i++ {
		current := base * float64(i)
		lower := current * (1 - tolerance)
		upper := current * (1 + tolerance)
		res = append(res, [2]float64{lower, upper})
	}
	return res
}

// calculateOccupiedCells 依据矩形面积匹配“占用格数”；未命中区间返回 0
func calculateOccupiedCells(x, y, w, h int, cellSize int) int {
	area := float64(w * h)
	ranges := createCellArray(cellSize, 9, 0.05)
	for i, r := range ranges { // i 从 0 开始，需 +1
		if area >= r[0] && area <= r[1] {
			return i + 1
		}
	}
	return 0
}

// findByColor 按 HSV 范围筛选轮廓并返回符合条件的 GameItem
func findByColor(img opencv.Mat) ([]GameItem, error) {
	// img := opencv.IMRead(imgPath, opencv.IMReadColor)
	// imgtest := images.CaptureScreen(0, 0, 0, 0, 0)
	//if img.Empty() {
	//	return nil, fmt.Errorf("无法读取图片: %s", imgPath)
	//}
	//defer img.Close()

	// 高斯模糊
	//blurred := opencv.NewMat()
	//defer blurred.Close()
	//opencv.GaussianBlur(img, &blurred, image.Pt(5, 5), 0, 0, opencv.BorderDefault)

	// BGR -> HSV
	hsv := opencv.NewMat()
	defer hsv.Close()
	opencv.CvtColor(img, &hsv, opencv.ColorBGRToHSV)

	// 调试：保存原始图像和HSV图像
	orig_img := util.MatToImage(img)
	images.Save(orig_img, "/sdcard/debug_original.png", 100)

	hsv_img := util.Mat3ToImage(hsv)
	images.Save(hsv_img, "/sdcard/debug_hsv.png", 100)

	yellowRange := util.HSVRange{
		LowerH: 12, LowerS: 42, LowerV: 50,
		UpperH: 16, UpperS: 90, UpperV: 67,
	}

	lower, upper := util.CreateHSVScalars(yellowRange)
	// InRange 生成掩膜
	mask := opencv.NewMat()
	opencv.InRangeWithScalar(hsv, lower, upper, &mask)

	// 调试：检查掩膜中的白色像素数量
	maskData := mask.ToBytes()
	whitePixels := 0
	for _, pixel := range maskData {
		if pixel == 255 {
			whitePixels++
		}
	}
	fmt.Printf("掩膜中检测到的白色像素数量: %d / %d (%.2f%%)\n", whitePixels, len(maskData), float64(whitePixels)/float64(len(maskData))*100)

	mask_img := util.MaskToImage(mask)
	images.Save(mask_img, "/sdcard/mask_img.png", 100)
	// 查找外部轮廓
	contours := opencv.FindContours(mask, opencv.RetrievalExternal, opencv.ChainApproxNone)

	out := make([]GameItem, 0)
	for i := 0; i < contours.Size(); i++ {
		c := contours.At(i)
		rect := opencv.BoundingRect(c)
		w, h := rect.Dx(), rect.Dy()
		fmt.Print("i:", i, "w:", w, "h:", h, "\n")
		if w > 65 && h > 65 {
			size := calculateOccupiedCells(rect.Min.X, rect.Min.Y, w, h, 68)
			if size > 0 {
				out = append(out, GameItem{
					X: rect.Min.X, Y: rect.Min.Y, W: w, H: h, Size: size, Color: "",
				})
			}
		}
	}
	return out, nil
}

// visualizeResults 在图上绘制矩形与编号，返回绘制后的图像
func visualizeResults(img opencv.Mat, boxes []GameItem) (opencv.Mat, error) {
	//img := opencv.IMRead(imgPath, opencv.IMReadColor)
	//if img.Empty() {
	//	return img, fmt.Errorf("无法读取图片: %s", imgPath)
	//}

	green := color.RGBA{0, 255, 0, 255}
	black := color.RGBA{0, 0, 0, 255}

	for i, it := range boxes {
		rect := image.Rect(it.X, it.Y, it.X+it.W, it.Y+it.H)
		opencv.Rectangle(&img, rect, green, 2)

		// 文本居中
		text := fmt.Sprintf("%d", i+1)
		size := opencv.GetTextSize(text, opencv.FontHersheySimplex, 0.6, 2)
		cx := it.X + it.W/2
		cy := it.Y + it.H/2
		tx := cx - size.X/2
		ty := cy + size.Y/2

		// 文字（不绘制白底，简化实现）
		opencv.PutText(&img, text, image.Pt(tx, ty), opencv.FontHersheySimplex, 0.6, black, 2)
	}

	return img, nil
}

// colorRGB 构造 BGR/Scalar（opencv 使用 BGR 排列）
func colorRGB(r, g, b uint8) opencv.Scalar {
	return opencv.NewScalar(float64(b), float64(g), float64(r), 0)
}

func main() {
	// imgPath := "./zhuangbei.png"
	img := images.CaptureScreen(570, 78, 1201, 638, 0)
	Matimg := util.ImageToMat(img)
	defer Matimg.Close()
	// fmt.Println("Matimg", Matimg)
	// 紫色
	// lower := opencv.NewScalar(112, 40, 46, 0)
	// upper := opencv.NewScalar(125, 79, 83, 0)
	// 黄色 - 扩大范围以便调试
	//lower := opencv.NewScalar(10, 30, 30, 0)  // H: 10-30, S: 30-255, V: 30-255
	//upper := opencv.NewScalar(30, 255, 255, 0)

	boxes, err := findByColor(Matimg)
	if err != nil {
		log.Fatalf("查找失败: %v", err)
	}

	for _, it := range boxes {
		fmt.Println(it)
	}

	vis, err := visualizeResults(Matimg, boxes)
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
