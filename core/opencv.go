package core

import (
	"app/assets"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/imgui"
	"github.com/Dasongzi1366/AutoGo/opencv"
)

// 用于保护缓存的互斥锁
var cacheMutex sync.RWMutex

// 用于保护OpenCV操作的全局互斥锁
var opencvMutex sync.Mutex

// OpenCVHandler 处理OpenCV相关操作
type OpenCVHandler struct{}

// NewOpenCVHandler 创建一个新的OpenCVHandler实例
func NewOpenCVHandler() *OpenCVHandler {
	return &OpenCVHandler{}
}

// ClearTemplateCache 清空模板缓存
func (h *OpenCVHandler) ClearTemplateCache() {
	opencvMutex.Lock()
	defer opencvMutex.Unlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 关闭所有缓存的Mat对象
	for _, mat := range assets.TemplateMap {
		if !mat.Empty() {
			mat.Close()
		}
	}
	for _, mat := range assets.MaskMap {
		if !mat.Empty() {
			mat.Close()
		}
	}

	// 清空缓存映射
	assets.TemplateMap = make(map[string]opencv.Mat)
	assets.MaskMap = make(map[string]opencv.Mat)

	fmt.Println("模板缓存已清空")
}

// FindImage 在指定区域内查找单个图像匹配
func (h *OpenCVHandler) FindImage(x1, y1, x2, y2 int, picName string, isGray bool, scale, sim float32) (int, int) {
	opencvMutex.Lock()
	defer opencvMutex.Unlock()

	if scale < 0.1 {
		scale = 1
	}
	template, _ := assets.ImageFile.ReadFile(picName)
	mat2, mat3 := byte2mat(&template, isGray, scale, picName)
	if mat2.Empty() {
		return -1, -1
	}

	img := images.CaptureScreen(x1, y1, x2, y2, 0)
	if img == nil {
		return -1, -1
	}

	bounds := img.Bounds()
	mat1, err := opencv.NewMatFromBytes(bounds.Dy(), bounds.Dx(), opencv.MatTypeCV8UC4, img.Pix)
	defer mat1.Close()
	if err != nil {
		return -1, -1
	}

	if isGray {
		mat1 = matGray(mat1)
	}

	result := opencv.NewMat()
	defer result.Close()

	opencv.MatchTemplate(mat1, mat2, &result, opencv.TmCcoeffNormed, mat3)

	_, maxVal, _, maxLoc := opencv.MinMaxLoc(result)
	if maxVal >= 0.5+sim*0.5 {
		return int(float32(maxLoc.X)/scale) + x1, int(float32(maxLoc.Y)/scale) + y1
	}

	return -1, -1
}

// FindImageAll 在指定区域内查找所有图像匹配
func (h *OpenCVHandler) FindImageAll(x1, y1, x2, y2 int, picName string, isGray bool, scale, sim float32) []map[string]interface{} {
	opencvMutex.Lock()
	defer opencvMutex.Unlock()

	if scale < 0.1 {
		scale = 1
	}

	resultT := make([]map[string]interface{}, 0)

	template, err := assets.ImageFile.ReadFile(picName)
	if err != nil {
		fmt.Println("读取模板图片失败:", err)
		return nil
	}

	mat2, mat3 := byte2mat(&template, isGray, scale, picName)
	if mat2.Empty() {
		fmt.Println("模板图片为空")
		return nil
	}

	img := images.CaptureScreen(x1, y1, x2, y2, 0)
	if img == nil {
		fmt.Println("截屏失败")
		return nil
	}

	bounds := img.Bounds()
	mat1, err := opencv.NewMatFromBytes(bounds.Dy(), bounds.Dx(), opencv.MatTypeCV8UC4, img.Pix)
	defer mat1.Close()
	if err != nil {
		fmt.Println("屏幕图片转换为 Mat 失败:", err)
		return nil
	}

	if isGray {
		mat1 = matGray(mat1)
	}

	result := opencv.NewMat()
	defer result.Close()

	if !mat2.Empty() {
		name := extractFileName(picName)
		for {
			opencv.MatchTemplate(mat1, mat2, &result, opencv.TmCcoeffNormed, mat3)
			_, maxVal, _, maxLoc := opencv.MinMaxLoc(result)

			if maxVal >= sim {
				x, y := int(float32(maxLoc.X)/scale), int(float32(maxLoc.Y)/scale)
				rect := image.Rectangle{
					Min: image.Point{X: x, Y: y},
					Max: image.Point{X: x + mat2.Cols(), Y: y + mat2.Rows()},
				}

				opencv.Rectangle(&mat1, rect, color.RGBA{R: 0, G: 255, B: 0, A: 255}, -1)
				resultT = append(resultT, map[string]interface{}{
					"name": name,
					"x1":   x,
					"y1":   y,
					"x2":   x + mat2.Cols(),
					"y2":   y + mat2.Rows(),
					"zx":   x + mat2.Cols()/2,
					"zy":   y + mat2.Rows()/2,
					"sim":  maxVal,
				})
			} else {
				break
			}
		}
	}

	if len(resultT) > 0 {
		return resultT
	}
	return nil
}

// byte2mat 将字节数据转换为OpenCV矩阵
func byte2mat(pngData *[]byte, isGray bool, scale float32, picName string) (opencv.Mat, opencv.Mat) {
	// 使用文件路径作为缓存键，而不是指针地址
	sign := fmt.Sprintf("%s-%t-%.2f", picName, isGray, scale)

	// 先尝试读取缓存
	cacheMutex.RLock()
	if cachedMat, ok := assets.TemplateMap[sign]; ok {
		maskMat := assets.MaskMap[sign]
		cacheMutex.RUnlock()
		return cachedMat, maskMat
	}
	cacheMutex.RUnlock()

	// 获取写锁来创建新的缓存条目
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 双重检查，避免并发创建同一个缓存项
	if cachedMat, ok := assets.TemplateMap[sign]; ok {
		return cachedMat, assets.MaskMap[sign]
	}

	img, _, err := image.Decode(bytes.NewReader(*pngData))
	if err != nil {
		fmt.Println("图像解码失败")
		return opencv.NewMat(), opencv.NewMat()
	}
	imgNrgba := ImageToNRGBA(img)

	bounds := imgNrgba.Bounds()
	templateMat, _ := opencv.NewMatFromBytes(bounds.Dy(), bounds.Dx(), opencv.MatTypeCV8UC4, imgNrgba.Pix)

	isTransparent := checkTransparent(imgNrgba)

	if isGray {
		templateMat = matGray(templateMat)
	}
	templateMat = matScale(templateMat, scale)

	var maskMat opencv.Mat
	if isTransparent {
		maskMat = createMask(imgNrgba)
	} else {
		maskMat = opencv.NewMat()
	}

	assets.TemplateMap[sign] = templateMat
	assets.MaskMap[sign] = maskMat

	return assets.TemplateMap[sign], assets.MaskMap[sign]
}

// matGray 将图像转换为灰度图
func matGray(mat opencv.Mat) opencv.Mat {
	grayMat := opencv.NewMat()
	opencv.CvtColor(mat, &grayMat, opencv.ColorBGRToGray)
	_ = mat.Close()
	return grayMat
}

// matScale 缩放图像
func matScale(mat opencv.Mat, scale float32) opencv.Mat {
	const epsilon = 1e-6
	if math.Abs(float64(scale-1)) < epsilon {
		return mat
	}
	fmt.Println("缩放")
	scaledMat := opencv.NewMat()
	opencv.Resize(mat, &scaledMat, image.Point{X: int(float32(mat.Cols()) * scale), Y: int(float32(mat.Rows()) * scale)}, 0, 0, opencv.InterpolationLinear)
	_ = mat.Close()
	return scaledMat
}

// checkTransparent 检查图像是否为透明图
func checkTransparent(img image.Image) bool {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width < 2 || height < 2 {
		return false
	}

	c0 := getRGB(img.At(0, 0))
	c1 := getRGB(img.At(width-1, 0))
	c2 := getRGB(img.At(0, height-1))
	c3 := getRGB(img.At(width-1, height-1))

	if c0 != c1 || c0 != c2 || c0 != c3 {
		return false
	}

	transparentCount := 0
	totalPixels := width * height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if getRGB(img.At(x, y)) == c0 {
				transparentCount++
			}
		}
	}

	if transparentCount >= int(float32(totalPixels)*0.3) && transparentCount < totalPixels {
		return true
	}

	return false
}

// createMask 创建透明图遮罩
func createMask(img image.Image) opencv.Mat {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	c0 := getRGB(img.At(0, 0))

	mask := opencv.NewMatWithSize(height, width, opencv.MatTypeCV8U)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if getRGB(img.At(x, y)) == c0 {
				mask.SetUCharAt(y, x, 1)
			} else {
				mask.SetUCharAt(y, x, 0)
			}
		}
	}

	return mask
}

// getRGB 获取RGB颜色值
func getRGB(c color.Color) color.RGBA {
	r, g, b, _ := c.RGBA() // 忽略 Alpha 通道
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
}

// ImageToNRGBA 将图像转换为NRGBA格式
func ImageToNRGBA(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	nrgbaImg := image.NewNRGBA(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				srcColor := img.At(x, y)
				r, g, b, a := srcColor.RGBA()
				i := (y-bounds.Min.Y)*nrgbaImg.Stride + (x-bounds.Min.X)*4
				nrgbaImg.Pix[i] = uint8(r >> 8)
				nrgbaImg.Pix[i+1] = uint8(g >> 8)
				nrgbaImg.Pix[i+2] = uint8(b >> 8)
				nrgbaImg.Pix[i+3] = uint8(a >> 8)
			}
		}(y)
	}
	wg.Wait()
	return nrgbaImg
}

// extractFileName 从路径中提取文件名
func extractFileName(path string) string {
	parts := strings.Split(path, "/")
	fileName := parts[len(parts)-1]
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// WaitFor 等待在指定区域检测到目标图像，返回检测到的坐标和是否成功
// x1, y1, x2, y2: 检测区域坐标
// picName: 图片名称
// isGray: 是否灰度匹配
// scale: 缩放比例
// sim: 相似度阈值
// interval: 检测间隔时间，默认1秒
// maxAttempts: 最大检测次数，默认60次
func (h *OpenCVHandler) WaitFor(x1, y1, x2, y2 int, picName string, isGray bool, scale, sim float32, interval time.Duration, maxAttempts int, context string) (int, int, bool) {
	if interval <= 0 {
		interval = time.Second
	}
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	for i := 0; i < maxAttempts; i++ {
		x, y := h.FindImage(x1, y1, x2, y2, picName, isGray, scale, sim)
		if x != -1 && y != -1 {
			return x, y, true
		}
		time.Sleep(interval)
		imgui.Toast(context)
	}
	return -1, -1, false
}

// ClickWhileExists 在指定区域内查找并点击目标图像，直到图像消失或达到最大点击次数
// x1, y1, x2, y2: 检测区域坐标
// picName: 图片名称
// isGray: 是否灰度匹配
// scale: 缩放比例
// sim: 相似度阈值
// interval: 检测间隔时间，默认1秒
// maxAttempts: 最大点击次数，默认60次
func (h *OpenCVHandler) ClickWhileExists(x1, y1, x2, y2 int, picName string, isGray bool, scale, sim float32, interval time.Duration, maxAttempts int) bool {
	if interval <= 0 {
		interval = time.Second
	}
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	for i := 0; i < maxAttempts; i++ {
		x, y := h.FindImage(x1, y1, x2, y2, picName, isGray, scale, sim)
		if x != -1 && y != -1 {
			// 找到图像，点击该位置
			Click(x, y)
			time.Sleep(interval)
		} else {
			// 没有找到图像，停止点击
			return true
		}
	}
	return false
}
