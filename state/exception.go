package state

import (
	"fmt"
	"time"
)

// 异常处理器
type ExceptionHandler struct {
	name        string        // 异常名称
	detector    func() bool   // 检测函数
	handler     func()        // 处理函数
	interval    time.Duration // 检测间隔
	lastChecked time.Time     // 上次检测时间
}

// 异常守护结构
type ExceptionGuard struct {
	stateMachine *StateMachine       // 状态机引用
	handlers     []*ExceptionHandler // 已注册的异常处理器
}

// 创建异常守护
func NewExceptionGuard(sm *StateMachine) *ExceptionGuard {
	eg := &ExceptionGuard{
		stateMachine: sm,
		handlers:     make([]*ExceptionHandler, 0),
	}
	return eg
}

// 注册异常处理器 - 支持可选的检测间隔（秒）
func (eg *ExceptionGuard) RegisterException(name string, detector func() bool, handler func(), intervalSeconds ...int) {
	// 默认检测间隔为10秒
	interval := 10 * time.Second

	// 如果提供了间隔参数，使用第一个参数
	if len(intervalSeconds) > 0 && intervalSeconds[0] > 0 {
		interval = time.Duration(intervalSeconds[0]) * time.Second
	}

	eg.handlers = append(eg.handlers, &ExceptionHandler{
		name:        name,
		detector:    detector,
		handler:     handler,
		interval:    interval,
		lastChecked: time.Time{}, // 初始化为零值，确保第一次检测会执行
	})

	fmt.Printf("注册异常处理器: %s (检测间隔: %v)\n", name, interval)
}

// 启动异常守护
func (eg *ExceptionGuard) Start() {
	// 启动检测线程
	go eg.runDetection()
}

// 运行检测循环
func (eg *ExceptionGuard) runDetection() {
	ticker := time.NewTicker(1 * time.Second) // 每秒检查一次是否需要执行检测
	defer ticker.Stop()

	fmt.Println("异常检测开始...")

	// 永久循环检测异常
	for range ticker.C {
		eg.checkExceptions()
	}
}

// 检查各种异常
func (eg *ExceptionGuard) checkExceptions() {
	now := time.Now()

	// 遍历所有已注册的异常处理器
	for _, handler := range eg.handlers {
		// 检查是否到了检测时间
		if now.Sub(handler.lastChecked) >= handler.interval {
			handler.lastChecked = now

			// 执行检测
			if handler.detector() {
				eg.handleException(handler.name, handler.handler)
				return // 一次只处理一个异常
			}
		}
	}
}

// 通用异常处理方法
func (eg *ExceptionGuard) handleException(exceptionName string, handlerFunc func()) {
	fmt.Printf("发现异常: %s\n", exceptionName)
	eg.stateMachine.StartExceptionHandling()

	handlerFunc() // 执行具体的处理逻辑

	eg.stateMachine.NotifyExceptionHandled()
}
