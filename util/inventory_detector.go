package util

import (
	"fmt"
	"image"
	"image/color"

	"github.com/Dasongzi1366/AutoGo/opencv"
)

// InventoryItem 仓库物品信息
type InventoryItem struct {
	X, Y             int    // 左上角坐标
	W, H             int    // 宽度和高度
	CenterX, CenterY int    // 中心点坐标
	Size             int    // 占用格数
	Color            string // 颜色标签
}

// String 实现 fmt.Stringer 接口
func (ii InventoryItem) String() string {
	return fmt.Sprintf("InventoryItem(x:%d, y:%d, w:%d, h:%d, centerX:%d, centerY:%d, size:%d, color:%s)",
		ii.X, ii.Y, ii.W, ii.H, ii.CenterX, ii.CenterY, ii.Size, ii.Color)
}

// ColorRange 颜色检测范围配置
type ColorRange struct {
	Name  string   // 颜色名称标签（如"黄色"）
	HSV   HSVRange // HSV 范围
	Color string   // 颜色英文标识（如"yellow"）
}

// CreateCellArray 根据 cellSize 生成面积区间（默认容差 5%，长度 9）
// 参数：
//   - cellSize: 单元格大小（如 68 或 110）
//   - length: 生成区间数量（通常为 9，表示 1-9 格）
//   - tolerance: 容差百分比（如 0.05 表示 5%）
func CreateCellArray(cellSize, length int, tolerance float64) [][2]float64 {
	base := float64(cellSize * cellSize) // 例如: 110*110 = 12100
	res := make([][2]float64, 0, length)
	for i := 1; i <= length; i++ {
		current := base * float64(i)
		lower := current * (1 - tolerance)
		upper := current * (1 + tolerance)
		res = append(res, [2]float64{lower, upper})
	}
	return res
}

// CalculateOccupiedCells 依据矩形面积匹配"占用格数"
// 参数：
//   - x, y: 矩形左上角坐标（当前未使用，预留用于未来扩展）
//   - w, h: 矩形宽度和高度
//   - cellSize: 单元格大小
//
// 返回：
//   - 占用格数（1-9），未命中区间返回 0
func CalculateOccupiedCells(x, y, w, h int, cellSize int) int {
	area := float64(w * h)
	ranges := CreateCellArray(cellSize, 9, 0.05)
	for i, r := range ranges { // i 从 0 开始，需 +1
		if area >= r[0] && area <= r[1] {
			return i + 1
		}
	}
	return 0
}

// FindByColor 按 HSV 范围筛选轮廓并返回符合条件的 InventoryItem
// 参数：
//   - img: 输入图像（*image.NRGBA）
//   - colorRanges: 要检测的颜色范围列表
//   - minSize: 最小宽高阈值（像素），过滤小于此值的轮廓
//   - cellSize: 单元格大小（用于计算占用格数）
//
// 返回：
//   - 检测到的仓库物品列表
//   - error: 错误信息（当前实现始终返回 nil）
func FindByColor(img *image.NRGBA, colorRanges []ColorRange, minSize, cellSize int) ([]InventoryItem, error) {
	Matimg := ImageToMat(img)
	defer Matimg.Close()

	// BGR -> HSV
	hsv := opencv.NewMat()
	defer hsv.Close()
	opencv.CvtColor(Matimg, &hsv, opencv.ColorBGRToHSV)

	// 存储所有颜色的结果
	out := make([]InventoryItem, 0)

	// 遍历每种颜色
	for _, cr := range colorRanges {
		lower, upper := CreateHSVScalars(cr.HSV)

		// 为当前颜色生成掩膜
		mask := opencv.NewMat()
		opencv.InRangeWithScalar(hsv, lower, upper, &mask)

		// 可选：保存掩膜用于调试
		// mask_img := MaskToImage(mask)
		// images.Save(mask_img, fmt.Sprintf("/sdcard/mask_%s.png", cr.Color), 100)

		// 查找外部轮廓
		contours := opencv.FindContours(mask, opencv.RetrievalExternal, opencv.ChainApproxNone)

		// 处理当前颜色的轮廓
		for i := 0; i < contours.Size(); i++ {
			c := contours.At(i)
			rect := opencv.BoundingRect(c)
			w, h := rect.Dx(), rect.Dy()

			// 过滤条件：宽高都大于最小阈值
			if w > minSize && h > minSize {
				size := CalculateOccupiedCells(rect.Min.X, rect.Min.Y, w, h, cellSize)
				if size > 0 {
					out = append(out, InventoryItem{
						X:       rect.Min.X,
						Y:       rect.Min.Y,
						W:       w,
						H:       h,
						CenterX: rect.Min.X + w/2,
						CenterY: rect.Min.Y + h/2,
						Size:    size,
						Color:   cr.Name,
					})
				}
			}
		}

		// 释放掩膜资源
		mask.Close()
	}

	return out, nil
}

// VisualizeResults 在图上绘制矩形与编号，返回绘制后的图像
// 参数：
//   - img: 输入图像（*image.NRGBA）
//   - items: 要绘制的仓库物品列表
//
// 返回：
//   - opencv.Mat: 绘制后的图像（调用方需要负责 Close()）
//   - error: 错误信息（当前实现始终返回 nil）
//
// 注意：返回的 Mat 需要调用方手动调用 Close() 释放资源
func VisualizeResults(img *image.NRGBA, items []InventoryItem) (opencv.Mat, error) {
	Matimg := ImageToMat(img)
	// 注意：这里不 defer Close()，因为返回值需要在外部使用

	green := color.RGBA{0, 255, 0, 255}
	black := color.RGBA{0, 0, 0, 255}

	for i, it := range items {
		rect := image.Rect(it.X, it.Y, it.X+it.W, it.Y+it.H)
		opencv.Rectangle(&Matimg, rect, green, 2)

		// 文本居中
		text := fmt.Sprintf("%d", i+1)
		size := opencv.GetTextSize(text, opencv.FontHersheySimplex, 0.6, 2)
		cx := it.X + it.W/2
		cy := it.Y + it.H/2
		tx := cx - size.X/2
		ty := cy + size.Y/2

		// 文字（不绘制白底，简化实现）
		opencv.PutText(&Matimg, text, image.Pt(tx, ty), opencv.FontHersheySimplex, 0.6, black, 2)
	}

	return Matimg, nil
}
