package console

// Init 初始化控制台。
// 参数:
//
//	noCaptureMode: 是否设置控制台画面禁止录屏,设置后控制台窗口不会影响脚本图色功能但是会导致投屏软件也无法显示控制台窗口,模拟器运行是此参数强制为true。
//
// 返回:
//
//	无返回值。
func Init(noCaptureMode bool) {}

// SetWindowSize 设置控制台窗口的宽高。
// 参数:
//
//	width: 控制台窗口的宽度。
//	height: 控制台窗口的高度。
//
// 返回:
//
//	无返回值。
func SetWindowSize(width, height int) {

}

// SetWindowPosition 设置控制台窗口的屏幕位置。
// 参数:
//
//	x: 控制台窗口的左上角横坐标。
//	y: 控制台窗口的左上角纵坐标。
//
// 返回:
//
//	无返回值。
func SetWindowPosition(x, y int) {

}

// SetWindowColor 设置控制台窗口的背景颜色。
// 参数:
//
//	color: 背景颜色的十六进制字符串，格式如 "#1E1F22"。
//
// 返回:
//
//	无返回值。
func SetWindowColor(color string) {

}

// SetTextColor 设置控制台文字颜色。
// 参数:
//
//	color: 文字颜色的十六进制字符串，格式如 "#FFFFFF"。
//
// 返回:
//
//	无返回值。
func SetTextColor(color string) {

}

// SetTextSize 设置控制台文字大小。
// 参数:
//
//	size: 文字大小。
//
// 返回:
//
//	无返回值。
func SetTextSize(size int) {

}

// Println 向控制台打印一行文本。
// 参数:
//
//	a: 要打印的任意类型的参数，行为类似 fmt.Println。
//
// 返回:
//
//	无返回值。
func Println(a ...any) {

}

// Clear 清空控制台中的所有文本内容。
// 返回:
//
//	无返回值。
func Clear() {

}

// Show 显示控制台窗口。
// 返回:
//
//	无返回值。
func Show() {

}

// Hide 隐藏控制台窗口。
// 返回:
//
//	无返回值。
func Hide() {

}
