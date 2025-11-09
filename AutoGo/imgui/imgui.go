package imgui

import "C"
import "image"

// TextItem 表示颜色和文本的组合。
type TextItem struct {
	TextColor string //文字颜色 格式如 #FFFFFF
	Text      string
}

type Hud struct {
}

type Rect struct {
}

type Line struct {
}

type Image struct {
}

// Init 初始化imgui。
// 参数:
//
//	noCaptureMode: 是否设置imgui画面禁止录屏,设置后imgui窗口不会影响脚本图色功能但是会导致投屏软件也无法显示imgui窗口,模拟器运行是此参数强制为true。
//
// 返回:
//
//	无返回值。
func Init(noCaptureMode bool) {}

// HudCreate 创建状态条。
// 参数:
//
//	x1, y1: 状态条的左上角坐标。
//	x2, y2: 状态条的右下角坐标。
//	bgColor: 状态条的背景颜色，格式如 #646464。
//	textSize: 状态条上的文字大小，如果小于等于 0，则默认使用 45。
//
// 返回:
//
//	Hud对象。
func HudCreate(x1, y1, x2, y2 int, bgColor string, fontSize int) *Hud {
	return nil
}

// SetText 设置状态条的文本。
// 参数:
//
//	items: 一个包含多个 TextItem 的切片，每个 TextItem 包含文字颜色和文本内容。
//
// 返回:
//
//	无返回值。
func (h *Hud) SetText(items []TextItem) {

}

// Destroy 销毁状态条。
// 返回:
//
//	无返回值。
func (h *Hud) Destroy() {

}

// RectCreate 创建矩形。
// 参数:
//
//	x1, y1: 矩形的左上角坐标。
//	x2, y2: 矩形的右下角坐标。
//	color: 矩形的边框颜色，格式如 #FFFFFF。
//	thickness: 矩形的边框宽度。
//
// 返回:
//
//	矩形对象。
func RectCreate(x1, y1, x2, y2 int, color string, thickness float32) *Rect {
	return nil
}

// Destroy 销毁矩形。
// 返回:
//
//	无返回值。
func (r *Rect) Destroy() {

}

// LineCreate 创建直线。
// 参数:
//
//	x1, y1: 起点坐标。
//	x2, y2: 终点坐标。
//	color: 线条颜色，格式如 #FFFFFF。
//	thickness: 线条宽度。
//
// 返回:
//
//	直线对象。
func LineCreate(x1, y1, x2, y2 int, color string, thickness float32) *Line {
	return nil
}

// Destroy 销毁直线。
// 返回:
//
//	无返回值。
func (l *Line) Destroy() {

}

// ImageCreate 创建图像。
// 参数:
//
//	x1, y1: 区域的左上角坐标。
//	x2, y2: 区域的右下角坐标。
//
// 返回:
//
//	图像对象。
func ImageCreate(x1, y1, x2, y2 int) *Image {
	return nil
}

// SetImage 设置图像。
// 参数:
//
//	img: 图片对象。
//	rotation: 图片旋转角度 0=不旋转，1=旋转90度，2=旋转180度，3=旋转270度。
//
// 返回:
//
//	无返回值。
func (i *Image) SetImage(img *image.NRGBA, rotation int) {

}

// Destroy 销毁图像。
// 返回:
//
//	无返回值。
func (i *Image) Destroy() {

}

// Toast 显示 Toast 提示信息。
// 参数:
//
//	message: 要显示的提示信息。
//	fontSize: 字体大小，可选参数。
//
// 返回:
//
//	无返回值。
func Toast(message string, fontSize ...int) {

}

// Alert 弹出提示信息，仅APP运行模式有效。
// 参数:
//
//	title: 标题。
//	message: 内容。
//
// 返回:
//
//	无返回值。
func Alert(title, message string) {

}
