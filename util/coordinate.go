package util

// Rect 定义矩形范围
type Rect struct {
	X1, Y1, X2, Y2 int
}

// ToGlobalCoords 将相对于范围的坐标转换为全局坐标
// rect: 范围 {x1, y1, x2, y2}
// x, y: 相对于范围的坐标（相对于x1, y1的偏移）
// 返回: 全局坐标 (globalX, globalY)
func ToGlobalCoords(rect Rect, x, y int) (int, int) {
	globalX := rect.X1 + x
	globalY := rect.Y1 + y
	return globalX, globalY
}

// ToGlobalCoordsWithValidation 带验证的转换函数
func ToGlobalCoordsWithValidation(rect Rect, x, y int) (globalX, globalY int, valid bool) {
	globalX = rect.X1 + x
	globalY = rect.Y1 + y

	// 验证转换后的坐标是否在范围内
	valid = globalX >= rect.X1 && globalX <= rect.X2 &&
		globalY >= rect.Y1 && globalY <= rect.Y2

	return globalX, globalY, valid
}

//func main() {
//	// 示例：范围从 (100, 200) 到 (500, 600)
//	rect := Rect{X1: 100, Y1: 200, X2: 500, Y2: 600}
//
//	// 相对坐标 (50, 80)
//	relativeX, relativeY := 50, 80
//
//	// 转换为全局坐标
//	globalX, globalY := ToGlobalCoords(rect, relativeX, relativeY)
//	fmt.Printf("相对坐标 (%d, %d) -> 全局坐标 (%d, %d)\n",
//		relativeX, relativeY, globalX, globalY)
//	// 输出: 相对坐标 (50, 80) -> 全局坐标 (150, 280)
//
//	// 使用带验证的函数
//	gx, gy, valid := ToGlobalCoordsWithValidation(rect, relativeX, relativeY)
//	fmt.Printf("全局坐标 (%d, %d), 是否在范围内: %v\n", gx, gy, valid)
//}
