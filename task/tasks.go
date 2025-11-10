package task

import (
	"app/core"
	"app/scene"
	"app/state"
	"app/util"
	"errors"
	"fmt"
	"time"

	"github.com/Dasongzi1366/AutoGo/storages"
)

// 全局变量存储账号ID
var 账号id string
var 游戏昵称 string

// 初始化函数，自动注册所有任务
func init() {
	// 注册所有任务
	// state.Register("买药", 买药, 1, 12*time.Hour)
	// 注册数据收集发送任务：初始化立即执行，然后每隔30分钟执行一次
	state.Register("数据收集发送", 数据收集发送, 0, 30*time.Minute)
}

func 打怪升级() error {
	return nil
}

func 礼包领取() error {
	var x, y = core.OpenCV.FindImage(211, 0, 408, 68, "img/main/礼包.png", false, 1, 0.7)
	if x > 0 && y > 0 {
		core.RandomClickInArea(x, y, x+5, y+5)
		core.RandomSleep(2000, 3000)
	}

	for i := 0; i < 5; i++ {
		x, y = core.Color.FindColor(352, 180, 421, 598, "ba0404", 0.9, 0)
		if x > 0 && y > 0 {
			core.RandomClickInArea(x-100, y, x-20, y+20)
			core.RandomSleep(1000, 1300)

			// 点击领取按钮
			x, y = core.Color.FindColor(413, 186, 1082, 600, "a60304", 0.9, 3)
			if x > 0 && y > 0 {
				core.RandomClickInArea(x-20, y, x, y+20)
				core.RandomSleep(1000, 1300)
				// 随机点击空白处
				core.RandomClickInArea(591, 132, 699, 166)
				core.RandomSleep(1000, 1300)
			}
		} else {
			continue
		}
	}

	core.CloseAllWindows()
	return nil
}

func 买药() error {
	type ItemConfig struct {
		X, Y   int // 检测坐标
		Amount int // 增加数量 (100, 200, 1000)
	}

	回城()

	var SceneName = scene.Identify("scene_map")
	core.Log("当前场景: " + SceneName)
	if SceneName != "阿斯加德城" {
		//TODO 去阿斯加德城
	}

	// 点击主城导航
	var x1, y1 = core.OpenCV.FindImage(0, 3, 442, 89, "img/main/主城导航.png", false, 1, 0.7)
	if x1 > 0 && y1 > 0 {
		// 点击药水商人
		core.RandomClickInArea(x1, y1, x1+5, y1+5)
		core.RandomSleep(2000, 3000)
	}

	// 配置多个道具的坐标和增加数量
	itemConfigs := []ItemConfig{
		{X: 873, Y: 237, Amount: 1000}, // 中药水
		{X: 583, Y: 393, Amount: 100},  // 回程卷轴
		{X: 854, Y: 383, Amount: 100},  // 钥匙
	}

	// 点击药水商人
	x1, y1 = core.OpenCV.FindImage(7, 3, 256, 340, "img/main/药水商人.png", false, 1, 0.7)
	if x1 > 0 && y1 > 0 {
		core.RandomClickInArea(x1, y1, x1+100, y1+20)
	} else {
		return errors.New("未找到药水商人图标")
	}

	// 等待药水商人出现
	if !core.OCR.WaitFor("잡화상인", 28, 8, 166, 55, 0.8, 2*time.Second, 120, "等待杂货商人出现") {
		return errors.New("等待杂货商人出现失败")
	}

	// 点击 批量购买设置
	found := core.OCR.ClickIfTextExists("일괄구매설정", 85, 662, 202, 692, 0.6)
	if !found {
		return errors.New("未找到 批量购买设置")
	}
	core.Sleep(2000)

	// 遍历所有道具配置
	for _, config := range itemConfigs {
		// 检测道具是否需要批量购买配置
		var flag = core.Color.CmpColor(config.X, config.Y, "1d2625", 0.98)
		if flag {
			// 点击道具
			core.Click(config.X, config.Y)
			core.Sleep(1000)

			// 点击重置
			core.RandomClickInArea(505, 479, 550, 520)
			core.Sleep(1000)

			// 根据配置的数量进行相应的点击
			switch config.Amount {
			case 1000:
				// 增加数量 +1000
				core.RandomClickInArea(726, 387, 781, 442)
			case 100:
				// 增加数量 +100
				core.RandomClickInArea(726, 313, 783, 364)
			default:
				// 对于其他数量，通过循环点击 +100 按钮
				if config.Amount > 100 && config.Amount%100 == 0 {
					clickTimes := config.Amount / 100
					for i := 0; i < clickTimes; i++ {
						core.RandomClickInArea(726, 313, 783, 364) // +100
						if i < clickTimes-1 {                      // 最后一次不需要等待
							core.Sleep(500)
						}
					}
				}
			}

			core.Sleep(1000)

			// 点击确认
			core.RandomClickInArea(677, 573, 769, 605)
			core.Sleep(1000)
		}
	}

	// 最终确认
	if core.OCR.ClickIfTextExists("설정저장", 709, 583, 810, 614, 0.6) {
		core.Sleep(1000)
	}

	// 批量购买
	if core.OCR.ClickIfTextExists("일괄구매", 264, 663, 358, 695, 0.6) {
		core.Sleep(1400)
		core.OCR.ClickIfTextExists("구매", 687, 630, 760, 667, 0.5)
	}

	// 关闭
	core.CloseAllWindows()
	return nil
}

func 回城() error {
	// 判断是否有退出地图按钮 (角色可能在副本内)
	var x, y = core.OpenCV.FindImage(184, 2, 311, 75, "img/main/退出地图.png", true, 1, 0.7)
	if x > 0 && y > 0 {
		core.RandomClickInArea(x, y, x+5, y+5)
		core.Sleep(1500)
		if core.OCR.ClickIfTextExists("이동", 702, 427, 759, 459, 0.6) {
			core.Sleep(4000)
		}
		println("1")
		core.OpenCV.WaitFor(1085, 455, 1210, 580, "img/main/跳跃.png", true, 1, 0.9, 2*time.Second, 60, "等待进入阿斯加德城")
		println("2")
	}

	// 看是否在阿斯加德城
	var SceneName = scene.Identify("scene_map")
	core.Toast("当前场景: " + SceneName)
	if SceneName == "阿斯加德城" {
		// 如果在阿斯加德城，直接返回
		return nil
	}

	// 判断是否有回程卷轴
	x, y = core.OpenCV.FindImage(817, 581, 1002, 704, "img/main/回程卷轴.png", false, 1, 0.7)
	if x > 0 && y > 0 {
		// 点击回程卷轴
		core.RandomClickInArea(x, y, x+5, y+5)
		core.Sleep(3000)
		core.OpenCV.WaitFor(1085, 455, 1210, 580, "img/main/跳跃.png", true, 1, 0.9, 2*time.Second, 60, "等待进入阿斯加德城")
	}

	// 如果还是不在阿斯加德城，尝试使用地图强制回程
	SceneName = scene.Identify("scene_map")
	core.Toast("当前场景: " + SceneName)

	if SceneName == "阿斯加德城" {
		// 如果在阿斯加德城，直接返回
		return nil
	}

	if SceneName == "未知场景" {
		return errors.New("无法识别当前场景，请检查游戏是否正常运行")
	}

	core.ClickMap()
	time.Sleep(3000)
	if core.OCR.ClickIfTextExists("아스가르드", 192, 9, 302, 43, 0.6) {
		core.RandomSleep(1000, 1500)
		core.RandomClickInArea(604, 319, 643, 340)
		core.RandomSleep(1000, 1500)
	}
	if core.OCR.ClickIfTextExists("잡화상인리바라", 1014, 256, 1136, 287, 0.6) {
		core.RandomSleep(1000, 1500)
		core.RandomClickInArea(1062, 656, 1138, 689)
		core.RandomSleep(1000, 1500)
	}

	return nil
}

// 获取用户ID函数
func 获取用户id() error {
	// 判断场景值获取用户ID
	for {
		sceneName := scene.Identify()
		core.Log("当前场景: " + sceneName)
		if sceneName == "主界面" || sceneName == "主界面2" {
			break
		}
		core.Sleep(3000)
	}

	core.RandomClickInArea(18, 22, 78, 81)
	// core.RandomSleep(2000, 3000)

	// 获取昵称编辑图标位置
	x1, _, _ := core.OpenCV.WaitFor(956, 585, 1258, 611, "img/editicon.png", true, 1, 0.8, 2*time.Second, 60, "等待进入个人信息")
	// 获取游戏昵称
	// fmt.Print("游戏昵称位置 x1:", x1, "\n")
	ganmeName := core.OCR.DetectAllText(958, 587, x1, 610)
	if len(ganmeName) != 0 {
		fmt.Print("游戏昵称: ", ganmeName, "\n")
		游戏昵称 = ganmeName[0]
	}

	// 获取游戏id
	gnameUserID := core.OCR.DetectAllText(990, 673, 1180, 692)
	if len(gnameUserID) != 0 {
		fmt.Print("获取用户id: ", gnameUserID, "\n")
		账号id = gnameUserID[0]
	}

	if len(游戏昵称) == 0 || len(账号id) == 0 {
		fmt.Print(游戏昵称, 账号id)
		return errors.New("游戏昵称或账号id获取失败")
	}
	return nil
}

// 获取仓库物品价值函数
func 获取仓库物品价值() (string, error) {
	// 判断场景值执行获取仓库物品价值
	for {
		sceneName := scene.Identify()
		core.Log("当前场景: " + sceneName)
		if sceneName == "界面仓库" {
			break
		}
		core.Sleep(5000)
	}

	x1, y1, _ := core.OpenCV.WaitFor(947, 4, 1279, 44, "img/仓库-哈夫币-icon.png", false, 1, 0.8, 2*time.Second, 60, "等待获取哈夫币图标")
	x2, _ := core.OpenCV.FindImage(947, 4, 1279, 44, "img/仓库-三角币-icon.png", false, 1, 0.8)

	x3 := x1 + 22
	y3 := y1

	x4 := x2
	y4 := 42

	hafubuNum := core.OCR.DetectAllText(x3, y3, x4, y4)
	if hafubuNum == nil {
		return "", errors.New(" OCR获取仓库价值获取失败")
	}

	fmt.Print("仓库价值: ", hafubuNum, "\n")
	return hafubuNum[0], nil
}

// HTTP数据发送函数
func 发送数据(userID string, gameName string, warehouseValue string) error {
	// 设备编号
	windowId := storages.Get("data", "windowId")
	data := map[string]interface{}{
		"deviceCode":      windowId,
		"accountId":       userID,
		"gameName":        gameName,
		"havalCoinAmount": warehouseValue,
	}

	_, _, err := util.HttpRequest.PostJSON(
		"http://175.24.153.109:4000/device/update",
		data,
		nil,
		1000*30,
	)

	if err != nil {
		fmt.Print(err)
		return err
	}

	fmt.Print(data)
	core.Log("数据发送成功")
	return nil
}

// 主任务函数：数据收集和发送
func 数据收集发送() error {
	// 等待异常处理完成
	state.Wait()

	// 判断账号ID是否存在
	if 账号id == "" {
		// 获取用户ID
		err := 获取用户id()
		if err != nil {
			return err
		}
		// 点击返回
		core.RandomClickInArea(16, 20, 45, 44)
		// 随机延迟
		core.RandomSleep(2000, 3000)
		// 点击仓库
		core.RandomClickInArea(23, 647, 133, 682)
		// 获取仓库物品价值
		warehouseValue, err := 获取仓库物品价值()
		if err != nil {
			return err
		}
		// 点击返回
		core.RandomClickInArea(16, 20, 45, 44)

		// 发送HTTP请求
		return 发送数据(账号id, 游戏昵称, warehouseValue)
	} else {
		// 账号ID已存在，直接获取仓库物品价值
		warehouseValue, err := 获取仓库物品价值()
		if err != nil {
			return err
		}

		// 发送HTTP请求
		return 发送数据(账号id, 游戏昵称, warehouseValue)
	}
}
