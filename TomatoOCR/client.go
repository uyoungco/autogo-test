// file: ppocr/client.go
package TomatoOCR

import (
	"encoding/json"
	"fmt"
)

// Client 是一个封装了 TomatoOCR 功能的客户端
type Client struct {
	instance *TmoOcr
}

// Config 用于初始化 Client 的配置
type Config struct {
	LicenseKey string
	Remark     string
	// 可以添加更多初始化时需要的配置
}

// NewClient 创建并初始化一个新的 OCR 客户端
func NewClient(config Config) (*Client, error) {
	instance := New()
	if instance == nil {
		return nil, fmt.Errorf("failed to create TomatoOCR instance")
	}

	// 设置 License
	// 假设 SetLicense 在成功时返回非空字符串，失败时返回空字符串
	if flag := instance.SetLicense(config.LicenseKey, config.Remark); flag == "" {
		// 可以在这里选择释放 instance，或者返回错误让调用者处理
		return nil, fmt.Errorf("failed to set license, license key might be invalid")
	}

	return &Client{instance: instance}, nil
}

// DetectOptions 定义了执行检测/识别时的各种参数
type DetectOptions struct {
	RecType           string  // "ch-3.0", "cht", "japan", "korean"
	DetBoxType        string  // "rect", "quad"
	DetUnclipRatio    float64 // 1.6-2.5
	RecScoreThreshold float64 // 0.1-0.9
	ReturnType        string  // "json", "text", "num"
	BinaryThresh      int     // 0-255
	RunMode           string  // "slow", "fast"
	DetectType        int     // 2 (recognize), 3 (detect+recognize)
}

// DefaultDetectOptions 返回一组推荐的默认配置
func DefaultDetectOptions() DetectOptions {
	return DetectOptions{
		RecType:           "ch-3.0",
		DetBoxType:        "rect",
		DetUnclipRatio:    1.9,
		RecScoreThreshold: 0.3,
		ReturnType:        "json",
		BinaryThresh:      0,
		RunMode:           "slow",
		DetectType:        3,
	}
}

// applyOptions 将配置应用到 OCR 实例
func (c *Client) applyOptions(opts DetectOptions) {
	c.instance.SetRecType(opts.RecType)
	c.instance.SetDetBoxType(opts.DetBoxType)
	c.instance.SetDetUnclipRatio(float32(opts.DetUnclipRatio))
	c.instance.SetRecScoreThreshold(float32(opts.RecScoreThreshold))
	c.instance.SetReturnType(opts.ReturnType)
	c.instance.SetBinaryThresh(opts.BinaryThresh)
	c.instance.SetRunMode(opts.RunMode)
}

// DetectInArea 在指定区域执行检测和识别
// 返回解析后的结果列表和可能的错误
func (c *Client) DetectInArea(x1, y1, x2, y2 int, opts DetectOptions) ([]DetectResult, error) {
	if opts.ReturnType != "json" {
		return nil, fmt.Errorf("this function requires ReturnType to be 'json'")
	}
	c.applyOptions(opts)

	jsonStr := c.instance.Detect(x1, y1, x2, y2, opts.DetectType)
	if jsonStr == "" {
		return nil, nil // 没有检测到结果，不是错误
	}

	var results []DetectResult
	err := json.Unmarshal([]byte(jsonStr), &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal detect results: %w", err)
	}
	return results, nil
}

// FindSingleTapPoint 在上一次的识别结果中查找单个文字的中心点
// 返回 [x, y] 坐标和可能的错误
func (c *Client) FindSingleTapPoint(word string) ([2]int, error) {
	pointStr := c.instance.FindTapPoint(word)
	if pointStr == "" {
		return [2]int{}, fmt.Errorf("word '%s' not found", word)
	}

	var point []int // 返回的是一个数组 [x, y]
	err := json.Unmarshal([]byte(pointStr), &point)
	if err != nil || len(point) != 2 {
		return [2]int{}, fmt.Errorf("failed to unmarshal tap point: %w", err)
	}

	return [2]int{point[0], point[1]}, nil
}

// FindAllTapPoints 在上一次的识别结果中查找所有匹配文字的中心点
func (c *Client) FindAllTapPoints(word string) ([]TapPoint, error) {
	pointsStr := c.instance.FindTapPoints(word)
	if pointsStr == "" {
		return nil, nil // 没有找到，不是错误
	}

	var results []TapPoint
	err := json.Unmarshal([]byte(pointsStr), &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tap points: %w", err)
	}
	return results, nil
}

// Close 释放 OCR 实例资源，如果库提供了释放函数
func (c *Client) Close() {
	// 假设 TomatoOCR 库提供了一个 Release 或 Close 方法
	// c.instance.Release()
	// 如果没有，可以留空，或者设为 nil 帮助GC
	c.instance = nil
}
