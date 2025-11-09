package state

import (
	"app/core"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
)

// 任务函数类型定义
type TaskFunc func() error

// 任务信息
type TaskInfo struct {
	Name         string        // 任务名称
	Func         TaskFunc      // 任务函数
	Order        int           // 执行顺序
	Interval     time.Duration // 执行间隔，0表示每轮都执行
	LastExecuted time.Time     // 上次执行时间
	Description  string        // 任务描述
}

// 任务注册表
type Registry struct {
	tasks map[string]*TaskInfo
}

// 全局注册表实例
var GlobalRegistry = NewRegistry()

// 创建新的注册表
func NewRegistry() *Registry {
	return &Registry{
		tasks: make(map[string]*TaskInfo),
	}
}

// 注册任务
func (r *Registry) Register(name string, taskFunc TaskFunc, order int, interval time.Duration) {
	r.tasks[name] = &TaskInfo{
		Name:     name,
		Func:     taskFunc,
		Order:    order,
		Interval: interval,
	}
	fmt.Printf("注册任务: %s (顺序: %d)\n", name, order)
}

// 获取所有任务，按顺序排序
func (r *Registry) GetOrderedTasks() []*TaskInfo {
	var tasks []*TaskInfo
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	// 按顺序排序
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Order < tasks[j].Order
	})

	return tasks
}

// 获取任务
func (r *Registry) GetTask(name string) (*TaskInfo, bool) {
	task, exists := r.tasks[name]
	return task, exists
}

// 检查任务是否应该执行
func (task *TaskInfo) ShouldExecute() bool {
	if task.Interval == 0 {
		return true // 每轮都执行
	}
	return time.Since(task.LastExecuted) >= task.Interval
}

// 执行任务
func (task *TaskInfo) Execute() error {
	fmt.Printf("执行任务: %s\n", task.Name)
	err := task.Func()
	if err == nil {
		task.LastExecuted = time.Now()
		core.Log("任务完成: " + task.Name)
	} else {
		wrappedErr := errors.WithStack(err)
		core.Log("任务执行失败: " + task.Name + ", 错误: " + fmt.Sprintf("%+v", wrappedErr))
	}
	return err
}

// 便捷的注册函数
func Register(name string, taskFunc TaskFunc, order int, interval time.Duration) {
	GlobalRegistry.Register(name, taskFunc, order, interval)
}
