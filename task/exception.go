package task

import (
	"app/core"
	"app/scene"
	"app/util"
	"fmt"

	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/imgui"
	"github.com/Dasongzi1366/AutoGo/utils"
)

type ExceptionGuardInterface interface {
	RegisterException(name string, detector func() bool, handler func(), intervalSeconds ...int)
}

func RegisterDefaultExceptions(eg ExceptionGuardInterface) {
	// 使用默认的10秒检测间隔
	// eg.RegisterException("检测游戏进程", 检测_游戏进程是否存在, 游戏进程开启)

	// 使用自定义的20秒检测间隔
	//eg.RegisterException("检测界面异常", 检测界面异常, 游戏进程开启, 20)

}

func 检测_游戏进程是否存在() bool {
	// 模拟检查服务器弹窗的逻辑
	packageName := app.CurrentPackage()
	gamePackage := core.Config.GetString("app_packages.ymir")
	if packageName != gamePackage {
		// 消抖
		core.Sleep(1000)
		// 再次检查
		packageName = app.CurrentPackage()
		if packageName != gamePackage {
			return true
		}

	}
	return false
}

func 检测界面异常() bool {
	var sceneName = scene.Identify()
	validScenes := []string{"登录_游戏资源更新弹框", "登录_游戏更新中", "登录_通知弹框", "登录_点击开始", "登录_谷歌登录弹窗", "登录_选择服务器", "登录_确认角色", ""}
	// 如果当前场景是登录相关场景，说明界面异常（游戏回到了登录界面）
	if util.StringInSlice(sceneName, validScenes) {
		return true
	}
	return false
}

func 游戏进程开启() {
	gamePackage := core.Config.GetString("app_packages.ymir")
	kitPackage := core.Config.GetString("app_packages.kit")

	// 杀掉后台进程 有时候游戏进程在但只是不在前台
	app.ForceStop(gamePackage)
	utils.Sleep(2000)

	// 启动KIT并且点击启动
	success := app.Launch(kitPackage, 0)
	utils.Sleep(2000)
	if success {
		var x, y = core.OpenCV.FindImage(0, 0, 720, 1280, "img/sys/kit_start.png", false, 1.0, 0.8)
		core.Click(x, y)
	}
	utils.Sleep(3000)
	// 启动游戏
	success = app.Launch(gamePackage, 0)
	if !success {
		fmt.Println("游戏启动失败，请检查应用包名配置")
	}

	var isDownload = false
	// for 直到游戏成功登录到游戏界面为止
	for {
		// 游戏进程开启 查看当前场景
		var sceneName = scene.Identify()
		fmt.Printf("Scene name: %s\n", sceneName)
		if sceneName == "登录_游戏资源更新弹框" {
			fmt.Println("游戏资源更新弹框已出现，开始处理...")
			utils.Sleep(2000)
			app.ForceStop(kitPackage)
			fmt.Println("关闭KIT弹框 使用本地网络更新")
			utils.Sleep(2000)
			core.RandomClickInArea(600, 440, 680, 464)
			utils.Sleep(3000)
		}

		if sceneName == "登录_游戏更新中" {
			// 从scene.json配置中获取OCR识别区域：507,644,769,682
			ocrText := core.OCR.DetectText(507, 644, 769, 682)
			imgui.Toast("游戏更新中: " + ocrText)
			utils.Sleep(3000)
			isDownload = true
		}

		if sceneName == "登录_通知弹框" {
			core.Click(128, 599)
		}

		if sceneName == "登录_点击开始" {
			core.RandomClickInArea(554, 483, 729, 518)
		}

		if sceneName == "登录_谷歌登录弹窗" {
			core.RandomClickInArea(545, 356, 829, 403)
			if isDownload {
				// 说明要登录谷歌号了 那么就需要开VPN 先吧游戏退出去
				app.ForceStop(gamePackage)
				fmt.Println("退出循环1")
				return
			}
		}

		if sceneName == "登录_选择服务器" {
			core.RandomClickInArea(603, 403, 688, 435)
		}

		if sceneName == "登录_确认角色" {
			core.RandomClickInArea(1133, 646, 1248, 680)
		}

		if sceneName == "主界面" {
			fmt.Println("退出循环2")
			return
		}

		utils.Sleep(2000)
	}
}
