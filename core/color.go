package core

import (
	"log"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
)

// ColorHandler 处理颜色识别相关操作
type ColorHandler struct {
	mu sync.Mutex
}

// NewColorHandler 创建一个新的ColorHandler实例
func NewColorHandler() *ColorHandler {
	return &ColorHandler{}
}

// Pixel 获取指定坐标点的颜色值
// x, y: 坐标点的位置
// 返回值: 表示颜色值的 "RRGGBB" 格式字符串
func (h *ColorHandler) Pixel(x, y int) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	color := images.Pixel(x, y, 0)
	if color == "" {
		log.Printf("获取坐标(%d, %d)的颜色失败", x, y)
		return ""
	}

	return color
}

// CmpColor 比较指定坐标点的颜色
// x, y: 坐标点的位置
// colorStr: 颜色字符串，格式如 "FFFFFF|CCCCCC-101010"
// sim: 相似度，范围 0.1 - 1.0
// 返回值: true 表示颜色匹配，false 表示颜色不匹配
func (h *ColorHandler) CmpColor(x, y int, colorStr string, sim float32) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sim < 0.1 || sim > 1.0 {
		log.Printf("相似度参数无效: %f，应该在0.1-1.0范围内", sim)
		return false
	}

	matched := images.CmpColor(x, y, colorStr, sim, 0)
	return matched
}

// FindColor 在指定区域内查找目标颜色
// x1, y1: 区域左上角的坐标
// x2, y2: 区域右下角的坐标，当 x2 或 y2 为 0 时，表示使用屏幕的最大宽度或高度
// colorStr: 颜色格式字符串，例如 "FFFFFF|CCCCCC-101010"
// sim: 相似度，范围 0.1 - 1.0
// dir: 查找方向，0-从左到右/从上到下，1-从右到左/从上到下，2-从左到右/从下到上，3-从右到左/从下到上
// 返回值: 找到颜色的坐标 (x, y)，找不到返回 (-1, -1)
func (h *ColorHandler) FindColor(x1, y1, x2, y2 int, colorStr string, sim float32, dir int) (int, int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sim < 0.1 || sim > 1.0 {
		log.Printf("相似度参数无效: %f，应该在0.1-1.0范围内", sim)
		return -1, -1
	}

	if dir < 0 || dir > 3 {
		log.Printf("查找方向参数无效: %d，应该在0-3范围内", dir)
		return -1, -1
	}

	x, y := images.FindColor(x1, y1, x2, y2, colorStr, sim, dir, 0)
	return x, y
}

// GetColorCountInRegion 计算指定区域内符合颜色条件的像素数量
// x1, y1: 区域左上角的坐标
// x2, y2: 区域右下角的坐标，当 x2 或 y2 为 0 时，表示使用屏幕的最大宽度或高度
// colorStr: 要查找的颜色字符串，格式为 "FFFFFF|CCCCCC-101010"
// sim: 相似度，范围 0.1 - 1.0
// 返回值: 符合条件的颜色像素数量，如果未找到符合条件的像素，则返回 0
func (h *ColorHandler) GetColorCountInRegion(x1, y1, x2, y2 int, colorStr string, sim float32) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sim < 0.1 || sim > 1.0 {
		log.Printf("相似度参数无效: %f，应该在0.1-1.0范围内", sim)
		return 0
	}

	count := images.GetColorCountInRegion(x1, y1, x2, y2, colorStr, sim, 0)
	return count
}

// DetectsMultiColors 根据指定的颜色串信息在屏幕进行多点颜色比对
// colors: 颜色模板字符串，例如 "369,1220,ffab2d-101010,370,1221,24b1ff-101010,380,390,907efd-101010"
// sim: 相似度，范围 0.1 - 1.0
// 返回值: true 表示比对成功，false 表示比对失败
func (h *ColorHandler) DetectsMultiColors(colors string, sim float32) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sim < 0.1 || sim > 1.0 {
		log.Printf("相似度参数无效: %f，应该在0.1-1.0范围内", sim)
		return false
	}

	result := images.DetectsMultiColors(colors, sim, 0)
	return result
}

// FindMultiColors 在指定区域内查找匹配的多点颜色序列
// x1, y1: 区域左上角的坐标
// x2, y2: 区域右下角的坐标，当 x2 或 y2 为 0 时，表示使用屏幕的最大宽度或高度
// colors: 颜色模板字符串，例如 "ffccff-151515,635,978,ffab2d-101010,6,29,24b1ff-101010,68,35,907efd-101010"
// sim: 相似度，范围 0.1 - 1.0
// dir: 查找方向，0-从左到右/从上到下，1-从右到左/从上到下，2-从左到右/从下到上，3-从右到左/从下到上
// 返回值: 匹配的首个颜色点的屏幕坐标 (x, y)，如果未找到则返回 (-1, -1)
func (h *ColorHandler) FindMultiColors(x1, y1, x2, y2 int, colors string, sim float32, dir int) (int, int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sim < 0.1 || sim > 1.0 {
		log.Printf("相似度参数无效: %f，应该在0.1-1.0范围内", sim)
		return -1, -1
	}

	if dir < 0 || dir > 3 {
		log.Printf("查找方向参数无效: %d，应该在0-3范围内", dir)
		return -1, -1
	}

	x, y := images.FindMultiColors(x1, y1, x2, y2, colors, sim, dir, 0)
	return x, y
}

// WaitForColor 等待在指定区域内检测到目标颜色
// x1, y1, x2, y2: 检测区域坐标
// colorStr: 目标颜色字符串
// sim: 相似度阈值 (0.1-1.0)
// maxAttempts: 最大检测次数，默认60次
// 返回值: 找到颜色的坐标 (x, y)，找不到返回 (-1, -1)
func (h *ColorHandler) WaitForColor(x1, y1, x2, y2 int, colorStr string, sim float32, maxAttempts int) (int, int) {
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	for i := 0; i < maxAttempts; i++ {
		x, y := h.FindColor(x1, y1, x2, y2, colorStr, sim, 0)
		if x != -1 && y != -1 {
			return x, y
		}
		time.Sleep(time.Second)
	}

	return -1, -1
}
