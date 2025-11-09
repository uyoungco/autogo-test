# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个基于 Go 语言的自动化测试框架，使用 AutoGo 库实现图像识别、OCR 文字识别、颜色检测等自动化操作功能。项目主要用于 Android 应用的自动化测试和操作。

## 核心架构

### 1. 核心模块 (core/)

项目采用全局单例模式，在 `core/global.go` 中暴露了四个核心处理器：

- **OCR**: OCR 文字识别处理器 (`core/ocr.go`)
- **OpenCV**: 图像识别处理器 (`core/opencv.go`)
- **Color**: 颜色检测处理器 (`core/color.go`)
- **API**: API 接口处理器 (`core/api.go`)

所有核心处理器都使用互斥锁保证线程安全。

### 2. OCR 处理器 (core/ocr.go)

使用 TomatoOCR 客户端进行文字识别，采用懒加载模式初始化。主要功能：

- `DetectText()`: 识别区域内第一个文字
- `DetectAllText()`: 识别区域内所有文字
- `FindText()`: 查找指定文字的坐标
- `WaitFor()`: 等待检测到目标文字
- `ClickWhileExists()`: 循环点击直到文字消失
- `ClickIfTextExists()`: 检测到文字则点击

文本相似度使用 Levenshtein 编辑距离算法计算。

### 3. OpenCV 处理器 (core/opencv.go)

负责图像识别和模板匹配，使用缓存机制优化性能：

- `FindImage()`: 在指定区域查找单个图像
- `FindImageAll()`: 查找所有匹配的图像
- `WaitFor()`: 等待检测到目标图像
- `ClickWhileExists()`: 循环点击直到图像消失
- `ClearTemplateCache()`: 清空模板缓存

模板图像缓存存储在 `assets.TemplateMap` 和 `assets.MaskMap` 中，支持透明图遮罩。

### 4. Color 处理器 (core/color.go)

提供颜色识别和比对功能：

- `Pixel()`: 获取指定坐标的颜色值
- `CmpColor()`: 比较指定坐标的颜色
- `FindColor()`: 在区域内查找目标颜色
- `FindMultiColors()`: 查找多点颜色序列
- `GetColorCountInRegion()`: 统计符合条件的像素数量
- `WaitForColor()`: 等待检测到目标颜色

### 5. 配置管理 (core/config.go)

使用 `ConfigManager` 管理配置，支持多层路径访问（如 `"app.packages.kit"`）。配置文件从 `assets/config/default.json` 加载。

全局数据存储使用 `DataStore` 结构体，包含用户信息和临时数据。

### 6. 资源管理 (assets/)

使用 Go embed 嵌入静态资源：

- `ImageFile`: 嵌入 `img/` 目录下的所有图片
- `ConfigFile`: 嵌入 `config/` 目录下的配置文件

### 7. AutoGo 依赖库

项目依赖本地 AutoGo 库（通过 `replace` 指令），提供底层自动化功能：

- `images`: 屏幕截图、图像处理
- `motion`: 点击、滑动等操作
- `opencv`: OpenCV 图像处理
- `ppocr`: OCR 识别
- `imgui`: UI 提示
- `storages`: 数据存储
- `system`: 系统操作
- `utils`: 工具函数

## 常用命令

### 构建项目
```bash
go build -o auto-go-test.exe
```

### 运行项目
```bash
go run main.go
```

### 运行测试文件
```bash
go run test1.go
```

### 安装依赖
```bash
go mod tidy
```

### 更新依赖
```bash
go get -u
go mod tidy
```

## 开发注意事项

### 线程安全

所有核心处理器都使用互斥锁保护，避免并发调用时出现问题。OCR 和 OpenCV 操作使用全局互斥锁 `ocrMutex` 和 `opencvMutex`。

### 资源管理

- 图像模板会被缓存以提高性能，缓存键格式为 `"文件名-是否灰度-缩放比例"`
- OpenCV Mat 对象需要手动调用 `Close()` 释放内存
- 使用 `defer` 确保资源正确释放

### 坐标系统

所有坐标参数格式为 `(x1, y1, x2, y2)`，表示矩形区域的左上角和右下角坐标。当 x2 或 y2 为 0 时，表示使用屏幕的最大宽度或高度。

### 相似度参数

- OCR 和颜色识别的相似度范围为 0.1 - 1.0
- OpenCV 图像匹配的相似度会被转换为 `0.5 + sim * 0.5` 的范围

### 辅助函数 (core/motion.go)

- `Click(x, y)`: 点击指定坐标
- `RandomClickInArea(x1, y1, x2, y2)`: 在区域内随机点击
- `Swipe(x1, y1, x2, y2, duration)`: 滑动操作
- `Sleep(ms)`: 延迟指定毫秒
- `RandomSleep(min, max)`: 随机延迟
- `Toast(context)`: 显示提示信息
- `CloseAllWindows()`: 关闭所有弹窗（通过识别 `img/close/` 目录下的关闭按钮图片）

## 项目结构

- `main.go`: 主入口文件（包含登录逻辑）
- `test1.go`: 测试文件（包含 OCR 和图像识别测试）
- `core/`: 核心功能模块
- `assets/`: 静态资源（图片、配置文件）
- `AutoGo/`: 本地依赖库
- `TomatoOCR/`: OCR 客户端封装
- `resources/`: 运行时资源（OCR 模型、native 库等）
