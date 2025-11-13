package util

import (
	"image"

	"github.com/Dasongzi1366/AutoGo/opencv"
)

// ImageToMat 将 *image.NRGBA 转换为 opencv.Mat (BGR格式)
func ImageToMat(img *image.NRGBA) opencv.Mat {
	if img == nil {
		return opencv.NewMat()
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建BGR格式的数据
	bgrData := make([]byte, width*height*3)

	// 将RGBA转换为BGR
	for i := 0; i < width*height; i++ {
		srcIndex := i * 4
		dstIndex := i * 3

		// RGBA -> BGR (跳过Alpha通道)
		bgrData[dstIndex] = img.Pix[srcIndex+2]   // B
		bgrData[dstIndex+1] = img.Pix[srcIndex+1] // G
		bgrData[dstIndex+2] = img.Pix[srcIndex]   // R
	}

	// 创建BGR格式的Mat (CV8UC3)
	mat, _ := opencv.NewMatFromBytes(height, width, opencv.MatTypeCV8UC3, bgrData)
	return mat
}

// MatToImage 将 opencv.Mat (BGR格式) 转换为 *image.NRGBA
func MatToImage(mat opencv.Mat) *image.NRGBA {
	if mat.Empty() {
		return nil
	}

	width := mat.Cols()
	height := mat.Rows()
	matData := mat.ToBytes()

	// 创建 NRGBA 图像
	bounds := image.Rect(0, 0, width, height)
	img := &image.NRGBA{
		Pix:    make([]uint8, width*height*4),
		Stride: width * 4,
		Rect:   bounds,
	}

	// 将BGR转换为RGBA格式
	for i := 0; i < width*height; i++ {
		srcIndex := i * 3
		dstIndex := i * 4

		// BGR -> RGBA
		img.Pix[dstIndex] = matData[srcIndex+2]   // R (BGR中的B)
		img.Pix[dstIndex+1] = matData[srcIndex+1] // G (BGR中的G)
		img.Pix[dstIndex+2] = matData[srcIndex]   // B (BGR中的R)
		img.Pix[dstIndex+3] = 255                 // A (完全不透明)
	}

	return img
}

// MaskToImage 将单通道掩膜 opencv.Mat (CV8UC1) 转换为 *image.NRGBA
func MaskToImage(mask opencv.Mat) *image.NRGBA {
	if mask.Empty() {
		return nil
	}

	// 获取掩膜数据
	maskData := mask.ToBytes()
	width := mask.Cols()
	height := mask.Rows()

	// 创建 NRGBA 图像
	bounds := image.Rect(0, 0, width, height)
	img := &image.NRGBA{
		Pix:    make([]uint8, width*height*4),
		Stride: width * 4,
		Rect:   bounds,
	}

	// 将单通道掩膜转换为 RGBA 格式
	for i := 0; i < len(maskData); i++ {
		pixelValue := maskData[i]
		baseIndex := i * 4

		// 白色像素 (255) -> 白色 RGBA，黑色像素 (0) -> 黑色 RGBA
		img.Pix[baseIndex] = pixelValue   // R
		img.Pix[baseIndex+1] = pixelValue // G
		img.Pix[baseIndex+2] = pixelValue // B
		img.Pix[baseIndex+3] = 255        // A (完全不透明)
	}

	return img
}

// Mat3ToImage 将三通道 opencv.Mat (CV8UC3) 转换为 *image.NRGBA
func Mat3ToImage(mat opencv.Mat) *image.NRGBA {
	if mat.Empty() {
		return nil
	}

	// 获取图像数据
	matData := mat.ToBytes()
	width := mat.Cols()
	height := mat.Rows()

	// 创建 NRGBA 图像
	bounds := image.Rect(0, 0, width, height)
	img := &image.NRGBA{
		Pix:    make([]uint8, width*height*4),
		Stride: width * 4,
		Rect:   bounds,
	}

	// 将三通道BGR/HSV转换为 RGBA 格式
	for i := 0; i < width*height; i++ {
		srcIndex := i * 3
		dstIndex := i * 4

		// 复制三个通道的数据
		img.Pix[dstIndex] = matData[srcIndex+2]   // R (BGR中的B或HSV中的V)
		img.Pix[dstIndex+1] = matData[srcIndex+1] // G (BGR中的G或HSV中的S)
		img.Pix[dstIndex+2] = matData[srcIndex]   // B (BGR中的R或HSV中的H)
		img.Pix[dstIndex+3] = 255                 // A (完全不透明)
	}

	return img
}

// HSVRange 定义 HSV 范围
type HSVRange struct {
	LowerH, LowerS, LowerV int
	UpperH, UpperS, UpperV int
}

// CreateHSVScalars 创建 HSV 上下限 Scalar（用于 InRange）
func CreateHSVScalars(hsvRange HSVRange) (opencv.Scalar, opencv.Scalar) {
	// HSV 在 OpenCV 中的范围：
	// H: 0-179, S: 0-255, V: 0-255
	lower := opencv.NewScalar(
		float64(hsvRange.LowerH),
		float64(hsvRange.LowerS),
		float64(hsvRange.LowerV),
		0.0,
	)

	upper := opencv.NewScalar(
		float64(hsvRange.UpperH),
		float64(hsvRange.UpperS),
		float64(hsvRange.UpperV),
		0.0,
	)

	return lower, upper
}
