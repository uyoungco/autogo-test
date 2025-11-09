package state

import (
	"fmt"
	"sync"
	"time"
)

// 全局状态机实例
var globalStateMachine *StateMachine

// 状态机结构
type StateMachine struct {
	currentTaskIndex    int           // 当前任务索引
	tasks               []*TaskInfo   // 任务列表
	nextTaskChan        chan bool     // 下一个任务信号
	exceptionChan       chan bool     // 异常处理信号
	isHandlingException bool          // 是否正在处理异常
	mutex               sync.RWMutex  // 读写锁
	taskInterval        time.Duration // 任务执行间隔
}

// 创建状态机
func NewStateMachine() *StateMachine {
	sm := &StateMachine{
		tasks:         GlobalRegistry.GetOrderedTasks(),
		nextTaskChan:  make(chan bool, 1),
		exceptionChan: make(chan bool, 1),
		taskInterval:  2 * time.Second, // 默认2秒间隔
	}

	// 设置全局实例
	globalStateMachine = sm

	fmt.Printf("加载了 %d 个任务\n", len(sm.tasks))
	for _, t := range sm.tasks {
		fmt.Printf("- %s (%s)\n", t.Name, t.Description)
	}

	// 启动第一个任务
	go func() {
		time.Sleep(1 * time.Second)
		sm.nextTaskChan <- true
	}()

	return sm
}

// 运行状态机
func (sm *StateMachine) Run() {
	fmt.Println("开始执行任务...")

	// 主循环 - 只处理两个信号
	for {
		select {
		case <-sm.nextTaskChan:
			// 执行下一个任务
			sm.executeNextTask()

		case <-sm.exceptionChan:
			// 异常处理完成，恢复执行
			sm.StopExceptionHandling()
			sm.scheduleNextTask()
		}
	}
}

// 执行下一个任务
func (sm *StateMachine) executeNextTask() {
	// 检查是否正在处理异常
	sm.mutex.RLock()
	if sm.isHandlingException {
		sm.mutex.RUnlock()
		// 异常处理期间，延迟执行常规任务
		sm.scheduleNextTask()
		return
	}
	sm.mutex.RUnlock()

	// 获取当前任务
	currentTask := sm.tasks[sm.currentTaskIndex]

	// 检查是否应该执行这个任务
	if !currentTask.ShouldExecute() {
		sm.moveToNextTask()
		sm.scheduleNextTask()
		return
	}

	// 执行任务
	currentTask.Execute()

	sm.moveToNextTask()
	sm.scheduleNextTask()
}

// 移动到下一个任务
func (sm *StateMachine) moveToNextTask() {
	sm.currentTaskIndex++
	if sm.currentTaskIndex >= len(sm.tasks) {
		sm.currentTaskIndex = 0 // 重新开始
		fmt.Println("任务循环，重新开始")
	}
}

// 调度下一个任务
func (sm *StateMachine) scheduleNextTask() {
	go func() {
		time.Sleep(sm.taskInterval)
		select {
		case sm.nextTaskChan <- true:
		default:
		}
	}()
}

// 开始异常处理
func (sm *StateMachine) StartExceptionHandling() {
	sm.mutex.Lock()
	sm.isHandlingException = true
	sm.mutex.Unlock()
	fmt.Println("暂停常规任务执行，开始异常处理...")
}

// 结束异常处理
func (sm *StateMachine) StopExceptionHandling() {
	sm.mutex.Lock()
	sm.isHandlingException = false
	sm.mutex.Unlock()
	fmt.Println("异常处理完成，恢复常规任务执行...")
}

// 异常处理完成，发送恢复信号
func (sm *StateMachine) NotifyExceptionHandled() {
	select {
	case sm.exceptionChan <- true:
	default:
	}
}

// Wait 等待异常处理完成，供任务内部调用
// 如果当前有异常任务在执行，则每秒检测一次，直到异常处理完成
func (sm *StateMachine) Wait() {
	for {
		sm.mutex.RLock()
		isHandling := sm.isHandlingException
		sm.mutex.RUnlock()

		if !isHandling {
			// 没有异常处理，可以继续执行
			return
		}

		// 有异常处理在进行，等待1秒后再检测
		time.Sleep(1 * time.Second)
	}
}

// Wait 全局等待函数，供任务内部调用
// 如果当前有异常任务在执行，则每秒检测一次，直到异常处理完成
func Wait() {
	if globalStateMachine != nil {
		globalStateMachine.Wait()
	}
}

// IsHandlingException 检查是否正在处理异常（只读方法）
func (sm *StateMachine) IsHandlingException() bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.isHandlingException
}

// IsHandlingException 全局检查函数，检查是否正在处理异常
func IsHandlingException() bool {
	if globalStateMachine != nil {
		return globalStateMachine.IsHandlingException()
	}
	return false
}
