package vdisplay

// Create 创建一个虚拟显示设备，该方法要求安卓10及以上版本。
// 参数:
//
//	width: 显示设备宽度。
//	height: 显示设备高度。
//	dpi: 显示设备的分辨率（每英寸点数）。
//
// 返回:
//
//	int: 创建成功后的虚拟显示设备ID。如果创建失败，返回错误码。
func Create(width, height, dpi int) int {
	return -1
}

// Destroy 销毁指定的虚拟显示设备。
// 参数:
//
//	displayId: 要销毁的虚拟显示设备ID。
//
// 功能:
//
//	释放与指定显示设备相关的所有资源。
func Destroy(displayId int) {

}
