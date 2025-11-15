# AutoGo 模块文档

[根目录](../CLAUDE.md) > **AutoGo**

> 本文档由 AI 架构师于 2025-11-15 17:46:19 生成

## 变更记录 (Changelog)

### 2025-11-15 17:46:19
- 初始化模块文档

---

## 模块职责

`AutoGo` 是一个**独立的 Go 模块**，提供 Android 自动化测试的底层能力，包括：

1. **图像处理** (images): 屏幕截图、图像读写
2. **设备操作** (motion): 点击、滑动、按键
3. **OpenCV** (opencv): 图像识别、模板匹配
4. **OCR** (ppocr): 文字识别
5. **UI 提示** (imgui): 悬浮窗提示
6. **数据存储** (storages): 键值存储
7. **系统操作** (system): 应用启动、停止
8. **工具函数** (utils): 通用工具
9. **无障碍服务** (uiacc): 无障碍 API
10. **虚拟显示** (vdisplay): 虚拟屏幕
11. **YOLO** (yolo): 目标检测
12. **Rhino** (rhino): JavaScript 引擎

## 入口与启动

AutoGo 通过 `go.mod` 的 `replace` 指令引用本地路径：

```go
replace github.com/Dasongzi1366/AutoGo => ./AutoGo
```

**模块声明** (AutoGo/go.mod):
```
module github.com/Dasongzi1366/AutoGo

go 1.25.0
```

## 对外接口

### images (图像处理)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `CaptureScreen(x1, y1, x2, y2, quality)` | 截取屏幕 | 区域, 质量 | image.Image |
| `LoadImage(path)` | 加载图片 | 路径 | image.Image |
| `SaveImage(path, img)` | 保存图片 | 路径, 图片 | error |

### motion (设备操作)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Click(x, y)` | 点击坐标 | x, y | 无 |
| `Swipe(x1, y1, x2, y2, duration)` | 滑动操作 | 起点, 终点, 时长 | 无 |
| `PressKey(keycode)` | 按键 | 键码 | 无 |
| `Input(text)` | 输入文字 | 文字 | 无 |

### opencv (OpenCV)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `MatchTemplate(source, template, threshold)` | 模板匹配 | 源图, 模板, 阈值 | (x, y, similarity) |
| `FindContours(img)` | 查找轮廓 | 图片 | []Contour |
| `CvtColor(img, code)` | 颜色转换 | 图片, 转换码 | Mat |

### ppocr (OCR)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Detect(img)` | 检测文字区域 | 图片 | []TextRegion |
| `Recognize(img)` | 识别文字 | 图片 | string |

### imgui (UI 提示)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Toast(message)` | 显示提示信息 | 消息 | 无 |
| `ShowFloatWindow(content)` | 显示悬浮窗 | 内容 | 无 |

### storages (数据存储)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Put(db, key, value)` | 存储键值对 | 数据库, 键, 值 | 无 |
| `Get(db, key)` | 获取值 | 数据库, 键 | string |
| `Remove(db, key)` | 删除键值对 | 数据库, 键 | 无 |

### system (系统操作)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `StartApp(packageName)` | 启动应用 | 包名 | error |
| `StopApp(packageName)` | 停止应用 | 包名 | error |
| `IsAppRunning(packageName)` | 检查应用是否运行 | 包名 | bool |

### utils (工具函数)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `RandomInt(min, max)` | 随机整数 | 最小值, 最大值 | int |
| `Sleep(ms)` | 延迟 | 毫秒 | 无 |
| `S2i(str)` | 字符串转整数 | 字符串 | int |

## 关键依赖与配置

### 依赖项

- **gocv.io/x/gocv**: OpenCV Go 绑定
- **标准库**: sync, time, math 等

### Native 库

AutoGo 依赖以下 Native 库（位于 `resources/libs/`）：

- `libopencv_core.so`: OpenCV 核心库
- `libopencv_imgcodecs.so`: 图像编解码
- `libopencv_imgproc.so`: 图像处理
- `libppocr.so`: OCR 库
- `libimgui.so`: UI 库
- `libyolo.so`: YOLO 目标检测
- `libncnn.so`: NCNN 推理框架
- `libc++_shared.so`: C++ 标准库

**支持架构**：
- arm64-v8a
- x86
- x86_64

## 数据模型

### Mat (OpenCV 图像)

```go
type Mat struct {
    p unsafe.Pointer
}
```

OpenCV 图像数据结构，需要手动调用 `Close()` 释放内存。

### Point (坐标点)

```go
type Point struct {
    X int
    Y int
}
```

### Rect (矩形区域)

```go
type Rect struct {
    X      int
    Y      int
    Width  int
    Height int
}
```

## 架构设计

### 模块分层

```
AutoGo/
├── app/           # 应用管理
├── console/       # 控制台
├── files/         # 文件操作
├── device/        # 设备信息
├── https/         # HTTPS 请求
├── images/        # 图像处理
├── ime/           # 输入法
├── imgui/         # UI 提示
├── media/         # 媒体操作
├── motion/        # 设备操作
├── opencv/        # OpenCV
├── ppocr/         # OCR
├── rhino/         # JavaScript 引擎
├── storages/      # 数据存储
├── system/        # 系统操作
├── uiacc/         # 无障碍服务
├── utils/         # 工具函数
├── vdisplay/      # 虚拟显示
└── yolo/          # YOLO 目标检测
```

### 调用关系

```
主项目 (auto-go-test)
    ↓
core/ 模块
    ↓
AutoGo/ 库
    ↓
Native 库 (resources/libs/)
```

## 测试与质量

### 当前测试覆盖

- 无自动化测试
- 依赖实际设备测试

### 建议补充

1. 单元测试：工具函数、数据结构
2. Mock 测试：模拟设备和 Native 调用
3. 集成测试：完整自动化流程

## 常见问题 (FAQ)

### 1. AutoGo 如何安装？

AutoGo 是本地模块，通过 `replace` 指令引用，无需安装。

### 2. Native 库如何部署？

Native 库位于 `resources/libs/`，需要部署到 Android 设备的对应目录。

### 3. 支持哪些 Android 架构？

- arm64-v8a（主要）
- x86（模拟器）
- x86_64（模拟器）

### 4. 如何调试 AutoGo？

- 使用 `imgui.Toast()` 显示调试信息
- 查看 Android logcat 日志
- 使用 GoLand 等 IDE 断点调试

### 5. OpenCV Mat 为什么需要手动释放？

Mat 对象包含 Native 内存，Go GC 无法自动回收，必须手动调用 `Close()`。

### 6. 如何更新 AutoGo？

直接修改 `AutoGo/` 目录下的代码，重新编译主项目即可。

### 7. 可以独立使用 AutoGo 吗？

可以。AutoGo 是独立模块，可以在其他项目中引用。

## 相关文件清单

### 核心源码

- `images/images.go`: 图像处理
- `motion/motion.go`: 设备操作
- `motion/keycode.go`: 键码定义
- `opencv/opencv.go`: OpenCV 封装
- `opencv/core.go`: OpenCV 核心
- `opencv/imgcodecs.go`: 图像编解码
- `opencv/imgproc.go`: 图像处理
- `opencv/imgproc_colorcodes.go`: 颜色代码
- `ppocr/ppocr.go`: OCR 封装
- `imgui/imgui.go`: UI 提示
- `storages/storages.go`: 数据存储
- `system/system.go`: 系统操作
- `utils/utils.go`: 工具函数
- `uiacc/uiacc.go`: 无障碍服务
- `vdisplay/vdisplay.go`: 虚拟显示
- `yolo/yolo.go`: YOLO 目标检测
- `rhino/rhino.go`: JavaScript 引擎
- `app/app.go`: 应用管理
- `console/console.go`: 控制台
- `files/files.go`: 文件操作
- `device/device.go`: 设备信息
- `https/https.go`: HTTPS 请求
- `media/media.go`: 媒体操作
- `ime/ime.go`: 输入法

### Native 库

- `../resources/libs/`: Native 库目录

### 相关模块

- `../core/`: 核心模块（调用 AutoGo）
- `../TomatoOCR/`: OCR 客户端（独立于 AutoGo）

---

**模块统计**:
- 文件数: 25+
- 子模块数: 17
- 支持架构数: 3
- Native 库数: 7+
