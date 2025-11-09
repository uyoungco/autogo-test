package main

import (
	"app/core"
	"app/state"
	"app/task"
	"fmt"
	"os"

	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/imgui"
	"github.com/Dasongzi1366/AutoGo/storages"
	"github.com/gin-gonic/gin"
)

// /mnt/shared/Pictures/

func uiCallback(c *gin.Context) bool {
	urlPath := c.Request.URL.Path

	if urlPath == "/login" {
		// 获取所有表单数据
		username := c.Query("username")
		password := c.Query("password")
		windowId := c.Query("windowId")
		isAppleStr := c.Query("isApple")

		storages.Put("data", "username", username)
		storages.Put("data", "password", password)
		storages.Put("data", "windowId", windowId)
		storages.Put("data", "isApple", isAppleStr)

		return true
	}

	if urlPath == "/getFormData" {
		username := storages.Get("data", "username")
		password := storages.Get("data", "password")
		windowId := storages.Get("data", "windowId")
		isApple := storages.Get("data", "isApple")

		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
			"windowId": windowId,
			"isApple":  isApple,
		})
	}

	if urlPath == "/close" {
		imgui.Toast("取消登录")
		os.Exit(0)
	}

	return false
}

func main() {
	username := storages.Get("data", "username")
	password := storages.Get("data", "password")
	windowId := storages.Get("data", "windowId")

	err := core.API.LoginAndSetup(username, password, windowId)

	if err {
		fmt.Print("登录失败")
		// system.RestartSelf()
		// return
	} else {
		fmt.Print("登录成功")
		packageName := app.CurrentPackage()
		println(packageName)
		// 创建状态机
		machine := state.NewStateMachine()

		// 创建异常守护
		exceptionGuard := state.NewExceptionGuard(machine)

		// 注册默认异常处理器
		task.RegisterDefaultExceptions(exceptionGuard)

		// 启动异常守护
		exceptionGuard.Start()

		// 启动状态机
		machine.Run()
	}
}
