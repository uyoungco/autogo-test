package core

import (
	"app/TomatoOCR"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/imgui"
)

// 用于保护OCR操作的全局互斥锁
var ocrMutex sync.Mutex

// OCRHandler 处理OCR相关操作
type OCRHandler struct {
	client *TomatoOCR.Client
	mu     sync.Mutex
	inited bool
}

// NewOCRHandler 创建一个新的OCRHandler实例
func NewOCRHandler() *OCRHandler {
	return &OCRHandler{}
}

// initClient 初始化OCR客户端（懒加载）
func (h *OCRHandler) initClient() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.inited {
		return nil
	}

	config := TomatoOCR.Config{
		LicenseKey: Config.GetString("ocr.license_key"),
		Remark:     "测试",
	}

	client, err := TomatoOCR.NewClient(config)
	if err != nil {
		return err
	}

	h.client = client
	h.inited = true
	return nil
}

// DetectText 在指定区域识别文字，返回识别到的第一个文字内容，识别不到返回空字符串
func (h *OCRHandler) DetectText(x1, y1, x2, y2 int) string {
	ocrMutex.Lock()
	defer ocrMutex.Unlock()

	if err := h.initClient(); err != nil {
		log.Printf("OCR客户端初始化失败: %v", err)
		return ""
	}

	opts := TomatoOCR.DefaultDetectOptions()
	results, err := h.client.DetectInArea(x1, y1, x2, y2, opts)
	if err != nil {
		log.Printf("OCR识别出错: %v", err)
		return ""
	}

	if len(results) == 0 {
		return ""
	}

	// 返回第一个识别到的文字
	return results[0].Words
}

// DetectAllText 在指定区域识别所有文字，返回所有识别到的文字内容切片，识别不到返回 nil
func (h *OCRHandler) DetectAllText(x1, y1, x2, y2 int) []string {
	ocrMutex.Lock()
	defer ocrMutex.Unlock()

	if err := h.initClient(); err != nil {
		log.Printf("OCR客户端初始化失败: %v", err)
		return nil
	}

	opts := TomatoOCR.DefaultDetectOptions()
	results, err := h.client.DetectInArea(x1, y1, x2, y2, opts)
	if err != nil {
		log.Printf("OCR识别出错: %v", err)
		return nil
	}

	if len(results) == 0 {
		return nil
	}

	var texts []string
	for _, result := range results {
		texts = append(texts, result.Words)
	}

	return texts
}

// FindText 查找指定文字的坐标，返回 x, y 坐标，找不到返回 -1, -1
func (h *OCRHandler) FindText(text string) (int, int) {
	ocrMutex.Lock()
	defer ocrMutex.Unlock()

	if err := h.initClient(); err != nil {
		log.Printf("OCR客户端初始化失败: %v", err)
		return -1, -1
	}

	point, err := h.client.FindSingleTapPoint(text)
	if err != nil {
		return -1, -1
	}

	return point[0], point[1]
}

// WaitFor 等待在指定区域检测到目标文字，返回是否检测到
// targetText: 目标文字
// x1, y1, x2, y2: 检测区域坐标
// similarity: 文本相似度阈值 (0.0-1.0)
// interval: 检测间隔时间，默认1秒
// maxAttempts: 最大检测次数，默认60次
func (h *OCRHandler) WaitFor(targetText string, x1, y1, x2, y2 int, similarity float32, interval time.Duration, maxAttempts int, context string) bool {
	if interval <= 0 {
		interval = time.Second
	}
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	for i := 0; i < maxAttempts; i++ {
		detectedText := h.DetectText(x1, y1, x2, y2)
		if detectedText != "" {
			textSimilarity := calculateTextSimilarity(targetText, detectedText)
			if textSimilarity >= similarity {
				return true
			}
		}
		imgui.Toast(context)
		time.Sleep(interval)
	}
	return false
}

// ClickWhileExists 在指定区域内查找并点击目标文字，直到文字消失或达到最大点击次数
// targetText: 目标文字
// x1, y1, x2, y2: 检测区域坐标
// similarity: 文本相似度阈值 (0.0-1.0)
// interval: 检测间隔时间，默认1秒
// maxAttempts: 最大点击次数，默认60次
func (h *OCRHandler) ClickWhileExists(targetText string, x1, y1, x2, y2 int, similarity float32, interval time.Duration, maxAttempts int) bool {
	if interval <= 0 {
		interval = time.Second
	}
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	for i := 0; i < maxAttempts; i++ {
		detectedText := h.DetectText(x1, y1, x2, y2)
		if detectedText != "" {
			textSimilarity := calculateTextSimilarity(targetText, detectedText)
			if textSimilarity >= similarity {
				// 找到文字，点击区域内随机坐标
				RandomClickInArea(x1, y1, x2, y2)
				time.Sleep(interval)
			} else {
				// 没有找到文字，停止点击
				return true
			}
		} else {
			// 没有检测到任何文字，停止点击
			return true
		}
	}
	return false
}

// ClickIfTextExists 判断指定区域内是否存在目标文字，如果存在则点击并返回true，否则返回false
// targetText: 目标文字
// x1, y1, x2, y2: 检测区域坐标
// similarity: 文本相似度阈值 (0.0-1.0)，默认0.8
func (h *OCRHandler) ClickIfTextExists(targetText string, x1, y1, x2, y2 int, similarity float32) bool {
	ocrMutex.Lock()
	defer ocrMutex.Unlock()

	if similarity <= 0 {
		similarity = 0.8
	}

	if err := h.initClient(); err != nil {
		log.Printf("OCR客户端初始化失败: %v", err)
		return false
	}

	// 直接执行OCR检测，避免重复加锁
	opts := TomatoOCR.DefaultDetectOptions()
	results, err := h.client.DetectInArea(x1, y1, x2, y2, opts)
	if err != nil {
		log.Printf("OCR识别出错: %v", err)
		return false
	}

	if len(results) == 0 {
		return false
	}

	detectedText := results[0].Words
	if detectedText != "" {
		textSimilarity := calculateTextSimilarity(targetText, detectedText)
		if textSimilarity >= similarity {
			// 找到匹配的文字，点击区域内随机坐标
			RandomClickInArea(x1, y1, x2, y2)
			return true
		}
	}

	return false
}

// calculateTextSimilarity 计算两个文本的相似度，使用编辑距离算法
func calculateTextSimilarity(text1, text2 string) float32 {
	// 去除空格和转换为小写，提高匹配准确性
	s1 := strings.ReplaceAll(strings.ToLower(text1), " ", "")
	s2 := strings.ReplaceAll(strings.ToLower(text2), " ", "")

	if s1 == s2 {
		return 1.0
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// 使用编辑距离算法计算相似度
	distance := levenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))

	similarity := 1.0 - float32(distance)/float32(maxLen)
	if similarity < 0 {
		similarity = 0
	}

	return similarity
}

// levenshteinDistance 计算两个字符串的编辑距离
func levenshteinDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)
	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// 创建距离矩阵
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// 初始化第一行和第一列
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// 填充矩阵
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				min(matrix[i-1][j]+1, matrix[i][j-1]+1), // 删除和插入的最小值
				matrix[i-1][j-1]+cost,                   // 替换
			)
		}
	}

	return matrix[len1][len2]
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
