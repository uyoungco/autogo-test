package core

import (
	"fmt"
	"sync"
)

// ColorHandler 处理颜色识别相关操作
type ApiHandler struct {
	mu sync.Mutex
}

// NewColorHandler 创建一个新的ColorHandler实例
func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

func (ctx *ApiHandler) LoginAndSetup(username, password, windowId string) bool {

	fmt.Printf("username: %v, password: %v, windowId: %v", username, password, windowId)

	if windowId == "" {
		return true
	}
	return false
}
