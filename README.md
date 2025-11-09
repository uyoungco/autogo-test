# AutoRun 游戏自动化框架

基于Go语言开发的Night Crows游戏自动化框架，采用状态机模式管理任务执行，结合图像识别、OCR文本识别和场景识别实现智能化游戏操作。

## 🚀 快速开始

### 1. 环境准备

#### AutoGo 自动化框架
- 参考官网：[AutoGo](https://autogo.cc/#/)
- 下载并配置 AutoGo 开发环境
- 确保设备连接正常

#### OCR 服务配置
- OCR 服务提供商：[番茄OCR](https://www.52tomato.com/console/buygoods)
- 需要购买卡密（测试卡几块钱即可）
- 在 `core/ocr.go` 中配置您的授权信息：

```go
config := TomatoOCR.Config{
    LicenseKey: "您的授权密钥",
    Remark:     "测试",
}
```

#### Screenshot 图色工具
- 自研的截图和图色识别工具
- **Mac电脑兼容**: 专门针对Mac系统优化，解决Mac平台兼容性问题
- 提供Web界面的图色调试功能
- 编译运行screenshot工具：

```bash
cd screenshot
go build -o screenshot-tool web_main.go
./screenshot-tool
```

### 2. 编译运行

```bash
# 编译运行
go build -o nightcrows main.go
./nightcrows
```

## 📁 项目结构

### 程序入口
- **`main.go`** - 程序入口，创建状态机和异常守护，启动自动化流程

### 核心模块 (core/)
- **`global.go`** - 全局实例定义
    - `OCR` - 全局OCR处理器实例
    - `OpenCV` - 全局图像识别处理器实例
    - `Color` - 全局颜色处理器实例

- **`config.go`** - 配置管理系统
    - `ConfigManager` - 配置管理器，支持多层路径配置访问
    - `Get()` - 获取配置值
    - `GetString()` - 获取字符串配置
    - `GetInt()` - 获取整数配置
    - `GetBool()` - 获取布尔配置
    - `GetMap()` - 获取映射配置
    - `DataStore` - 全局数据存储结构
    - `GetAppPackage()` - 获取应用包名

- **`ocr.go`** - OCR文本识别处理
    - `DetectText()` - 识别指定区域文字
    - `DetectAllText()` - 识别区域内所有文字
    - `FindText()` - 查找指定文字坐标
    - `WaitFor()` - 等待目标文字出现
    - `ClickWhileExists()` - 持续点击直到文字消失
    - `ClickIfTextExists()` - 检测文字存在则点击

- **`opencv.go`** - OpenCV图像识别处理
    - `FindImage()` - 查找单个图像匹配
    - `FindImageAll()` - 查找所有图像匹配
    - `WaitFor()` - 等待图像出现
    - `ClickWhileExists()` - 持续点击直到图像消失
    - `ClearTemplateCache()` - 清空模板缓存

- **`motion.go`** - 动作控制和交互
    - `Click()` - 点击指定坐标
    - `RandomClickInArea()` - 区域内随机点击
    - `Swipe()` - 滑动操作
    - `Sleep()` - 延时等待
    - `RandomSleep()` - 随机延时
    - `CloseAllWindows()` - 关闭所有弹窗

- **`color.go`** - 颜色识别和检测
    - `Pixel()` - 获取指定坐标颜色值
    - `CmpColor()` - 比较指定坐标颜色
    - `FindColor()` - 查找目标颜色
    - `GetColorCountInRegion()` - 统计区域内符合条件的像素数量
    - `DetectsMultiColors()` - 多点颜色比对
    - `FindMultiColors()` - 查找匹配的多点颜色序列
    - `WaitForColor()` - 等待目标颜色出现