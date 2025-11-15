package util

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// StringToIntMap 使用strings.Map过滤字符
func StringToIntMap(s string) (int, error) {
	// 定义过滤函数：保留数字和负号，去除空格、逗号等
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '-' || r == '+' {
			return r
		}
		return -1 // 返回-1表示删除该字符
	}, s)

	if cleaned == "" {
		return 0, fmt.Errorf("字符串不包含有效数字")
	}

	result, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0, fmt.Errorf("转换失败: %v", err)
	}

	return result, nil
}
