package scene

import (
	"app/assets"
	"app/core"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// 预定义的区域类型映射
var regionTypes = map[string]string{
	"map_name": "57,14,188,40", // 地图名区域
}

// SceneConfig 场景配置结构
type SceneConfig struct {
	Region     string  `json:"region"`     // 检测区域 "x1,y1,x2,y2"
	RegionType string  `json:"regionType"` // 区域类型，优先于region使用
	Images     string  `json:"images"`     // 图片名称，支持 | (或) 和 & (且) 操作符
	Text       string  `json:"text"`       // OCR文本识别内容
	IsGray     bool    `json:"isGray"`     // 是否灰度匹配
	Sim        float32 `json:"sim"`        // 相似度阈值
	Scale      float32 `json:"scale"`      // 缩放比例
}

// OCRCache OCR结果缓存结构
type OCRCache struct {
	Region string // 区域坐标
	Text   string // OCR识别结果
}

// SceneManager 场景管理器
type SceneManager struct {
	scenes   map[string]SceneConfig // 默认场景配置
	sceneMap map[string]SceneConfig // 地图场景配置
	opencv   *core.OpenCVHandler
	ocr      *core.OCRHandler
	mutex    sync.Mutex        // 添加互斥锁防止并发访问
	ocrCache map[string]string // OCR结果缓存 region -> text
}

var manager *SceneManager

// init 包初始化
func init() {
	manager = &SceneManager{
		scenes:   make(map[string]SceneConfig),
		sceneMap: make(map[string]SceneConfig),
		opencv:   core.NewOpenCVHandler(),
		ocr:      core.NewOCRHandler(),
		ocrCache: make(map[string]string),
	}
	loadSceneConfig()
}

// loadSceneConfig 加载所有场景配置
func loadSceneConfig() {
	// 加载默认场景配置
	loadConfigFile("config/scene.json", manager.scenes, "默认场景")

	// 加载地图场景配置
	// loadConfigFile("config/scene_map.json", manager.sceneMap, "地图场景")
}

// loadConfigFile 加载指定的配置文件
func loadConfigFile(filePath string, targetMap map[string]SceneConfig, configName string) {
	configData, err := assets.ConfigFile.ReadFile(filePath)
	if err != nil {
		fmt.Printf("加载%s配置失败: %v\n", configName, err)
		return
	}

	var tempConfig map[string]SceneConfig
	err = json.Unmarshal(configData, &tempConfig)
	if err != nil {
		fmt.Printf("解析%s配置失败: %v\n", configName, err)
		return
	}

	// 为所有场景设置默认值并添加到目标map
	for sceneName, config := range tempConfig {
		targetMap[sceneName] = normalizeSceneConfig(config)
	}

	fmt.Printf("成功加载 %d 个%s配置\n", len(tempConfig), configName)
}

// cleanMapName 清理地图名，去掉数字及后面的字符（如: 폴크방1층 -> 폴크방）
func cleanMapName(text string) string {
	// 遍历字符串，找到第一个数字的位置
	for i, r := range text {
		if r >= '0' && r <= '9' {
			return text[:i]
		}
	}
	return text
}

// calculateTextSimilarity 计算两个文本的相似度，使用编辑距离算法
func calculateTextSimilarity(text1, text2 string) float32 {
	// 清理地图名，去掉数字及后面的字符
	cleanedText1 := cleanMapName(text1)
	cleanedText2 := cleanMapName(text2)

	// 去除空格和转换为小写，提高匹配准确性
	s1 := strings.ReplaceAll(strings.ToLower(cleanedText1), " ", "")
	s2 := strings.ReplaceAll(strings.ToLower(cleanedText2), " ", "")

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

// normalizeSceneConfig 规范化场景配置，设置默认值
func normalizeSceneConfig(config SceneConfig) SceneConfig {
	// 优先使用RegionType，如果设置了regionType，则转换为具体坐标
	if config.RegionType != "" {
		if regionCoord, exists := regionTypes[config.RegionType]; exists {
			config.Region = regionCoord
		} else {
			fmt.Printf("警告: 未知的区域类型 '%s'，使用默认区域\n", config.RegionType)
		}
	}

	// 设置默认region为全屏
	if config.Region == "" {
		config.Region = "0,0,0,0"
	}

	// 设置默认相似度阈值
	if config.Sim == 0 {
		config.Sim = 0.8
	}

	// 设置默认缩放比例
	if config.Scale == 0 {
		config.Scale = 1.0
	}

	// 验证配置的有效性
	if config.Images == "" && config.Text == "" {
		fmt.Printf("警告: 场景配置必须包含 'images' 或 'text' 字段之一\n")
	}

	// OCR模式下，region是必填的，确保有有效区域
	if config.Text != "" && config.Region == "0,0,0,0" {
		fmt.Printf("警告: OCR模式建议指定具体的检测区域，而不是全屏检测\n")
	}

	// isGray 默认为 false，不需要特殊处理

	return config
}

// normalizeImagePath 规范化图片路径，自动添加img/前缀
func normalizeImagePath(imagePath string) string {
	imagePath = strings.TrimSpace(imagePath)
	if strings.HasPrefix(imagePath, "img/") {
		return imagePath
	}
	return "img/" + imagePath
}

// parseRegion 解析区域字符串 "x1,y1,x2,y2"
func parseRegion(region string) (int, int, int, int) {
	parts := strings.Split(region, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0
	}

	x1, _ := strconv.Atoi(parts[0])
	y1, _ := strconv.Atoi(parts[1])
	x2, _ := strconv.Atoi(parts[2])
	y2, _ := strconv.Atoi(parts[3])

	return x1, y1, x2, y2
}

// checkImageCondition 检查图片条件
func (sm *SceneManager) checkImageCondition(images string, x1, y1, x2, y2 int, isGray bool, scale, sim float32) bool {
	// 处理 & (且) 操作符
	if strings.Contains(images, "&") {
		imageParts := strings.Split(images, "&")
		for _, img := range imageParts {
			imgPath := normalizeImagePath(img)
			x, y := sm.opencv.FindImage(x1, y1, x2, y2, imgPath, isGray, scale, sim)
			if x == -1 && y == -1 {
				return false // 只要有一个图片没找到，就返回false
			}
		}
		return true // 所有图片都找到了
	}

	// 处理 | (或) 操作符
	if strings.Contains(images, "|") {
		imageParts := strings.Split(images, "|")
		for _, img := range imageParts {
			imgPath := normalizeImagePath(img)
			x, y := sm.opencv.FindImage(x1, y1, x2, y2, imgPath, isGray, scale, sim)
			if x != -1 && y != -1 {
				return true // 只要找到一个图片，就返回true
			}
		}
		return false // 所有图片都没找到
	}

	// 单个图片检查
	imgPath := normalizeImagePath(images)
	x, y := sm.opencv.FindImage(x1, y1, x2, y2, imgPath, isGray, scale, sim)
	return x != -1 && y != -1
}

// checkSceneCondition 检查场景条件，支持图像识别和OCR文本识别
func (sm *SceneManager) checkSceneCondition(config SceneConfig) bool {
	// 如果配置了文本识别，使用OCR模式
	if config.Text != "" {
		return sm.checkTextCondition(config.Text, config.Region, config.Sim)
	}

	// 否则使用图像识别模式
	if config.Images != "" {
		x1, y1, x2, y2 := parseRegion(config.Region)
		return sm.checkImageCondition(config.Images, x1, y1, x2, y2, config.IsGray, config.Scale, config.Sim)
	}

	return false
}

// getOrCacheOCRText 获取或缓存OCR识别结果
func (sm *SceneManager) getOrCacheOCRText(region string) string {
	// 检查缓存
	if cachedText, exists := sm.ocrCache[region]; exists {
		return cachedText
	}

	// 解析区域坐标
	x1, y1, x2, y2 := parseRegion(region)

	// 执行OCR识别
	detectedText := sm.ocr.DetectText(x1, y1, x2, y2)

	// 缓存结果
	sm.ocrCache[region] = detectedText

	return detectedText
}

// clearOCRCache 清理OCR缓存
func (sm *SceneManager) clearOCRCache() {
	sm.ocrCache = make(map[string]string)
}

// checkTextCondition OCR文本识别检查（使用缓存优化）
func (sm *SceneManager) checkTextCondition(expectedText string, region string, simThreshold float32) bool {
	// 使用缓存的OCR结果
	detectedText := sm.getOrCacheOCRText(region)
	if detectedText == "" {
		return false
	}

	// 计算文本相似度
	similarity := calculateTextSimilarity(expectedText, detectedText)

	return similarity >= simThreshold
}

// Identify 识别当前场景 (默认使用 scene.json)
func Identify(configType ...string) string {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if len(configType) > 0 {
		return manager.identifyWithConfig(configType[0])
	}
	return manager.identifyWithConfig("scene")
}

// identifyWithConfig 使用指定配置执行场景识别
func (sm *SceneManager) identifyWithConfig(configType string) string {
	// 清理OCR缓存，确保每次识别使用最新的屏幕内容
	sm.clearOCRCache()

	var targetScenes map[string]SceneConfig

	switch configType {
	case "scene_map":
		targetScenes = sm.sceneMap
	case "scene":
		fallthrough
	default:
		targetScenes = sm.scenes
	}

	for sceneName, config := range targetScenes {
		if sm.checkSceneCondition(config) {
			return sceneName
		}
	}

	return "未知场景"
}

// identify 执行场景识别 (保持向后兼容)
func (sm *SceneManager) identify() string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return sm.identifyWithConfig("scene")
}
