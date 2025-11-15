# core 模块文档

[根目录](../CLAUDE.md) > **core**

> 本文档由 AI 架构师于 2025-11-15 17:46:19 生成

## 变更记录 (Changelog)

### 2025-11-15 17:46:19
- 初始化模块文档

---

## 模块职责

`core` 模块是整个自动化框架的**核心基础层**，提供以下能力：

1. **OCR 文字识别** (ocr.go)
2. **OpenCV 图像识别** (opencv.go)
3. **颜色检测与比对** (color.go)
4. **配置管理** (config.go)
5. **API 接口调用** (api.go)
6. **WebSocket 通信** (ws_client.go)
7. **设备操作（点击、滑动等）** (motion.go)

所有处理器均采用**全局单例模式**，通过互斥锁保证线程安全。

## 入口与启动

### 全局实例初始化

```go
// core/global.go
var OCR = NewOCRHandler()          // OCR 处理器
var OpenCV = NewOpenCVHandler()    // OpenCV 处理器
var Color = NewColorHandler()      // 颜色处理器
var API = NewApiHandler()          // API 处理器
```

### 配置加载

配置在包初始化时自动加载：

```go
// core/config.go init()
func init() {
    Config = &ConfigManager{}
    GlobalStore = NewDataStore()
    loadConfig() // 从 assets/config/default.json 加载
}
```

## 对外接口

### OCR 处理器 (OCRHandler)

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `DetectText(x1, y1, x2, y2)` | 识别区域内第一个文字 | 坐标区域 | string |
| `DetectAllText(x1, y1, x2, y2)` | 识别区域内所有文字 | 坐标区域 | []string |
| `FindText(x1, y1, x2, y2, text, sim)` | 查找指定文字坐标 | 区域, 文字, 相似度 | (x, y, bool) |
| `WaitFor(text, x1, y1, x2, y2, sim, interval, maxRetries, desc)` | 等待文字出现 | 文字, 区域, 相似度, 间隔, 最大重试, 描述 | (x, y, bool) |
| `ClickWhileExists(text, x1, y1, x2, y2, sim, maxRetries)` | 循环点击直到文字消失 | 文字, 区域, 相似度, 最大重试 | bool |
| `ClickIfTextExists(text, x1, y1, x2, y2, sim)` | 检测到文字则点击 | 文字, 区域, 相似度 | bool |

**线程安全**: 所有方法使用全局 `ocrMutex` 保护。

### OpenCV 处理器 (OpenCVHandler)

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `FindImage(x1, y1, x2, y2, imgPath, isGray, scale, sim)` | 查找单个图像 | 区域, 图片路径, 灰度, 缩放, 相似度 | (x, y) |
| `FindImageAll(x1, y1, x2, y2, imgPath, isGray, scale, sim)` | 查找所有匹配图像 | 同上 | []Point |
| `WaitFor(x1, y1, x2, y2, imgPath, isGray, scale, sim, interval, maxRetries, desc)` | 等待图像出现 | 区域, 图片路径, 灰度, 缩放, 相似度, 间隔, 最大重试, 描述 | (x, y, bool) |
| `ClickWhileExists(x1, y1, x2, y2, imgPath, isGray, scale, sim, maxRetries)` | 循环点击直到图像消失 | 区域, 图片路径, 灰度, 缩放, 相似度, 最大重试 | bool |
| `ClearTemplateCache()` | 清空模板缓存 | 无 | 无 |

**缓存机制**: 模板图像缓存在 `assets.TemplateMap` 和 `assets.MaskMap` 中，缓存键格式为 `"文件名-是否灰度-缩放比例"`。

**线程安全**: 所有方法使用全局 `opencvMutex` 保护。

### 颜色处理器 (ColorHandler)

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Pixel(x, y)` | 获取指定坐标颜色值 | x, y | string (十六进制) |
| `CmpColor(x, y, color, sim)` | 比较指定坐标颜色 | x, y, 颜色, 相似度 | bool |
| `FindColor(x1, y1, x2, y2, color, sim, dir)` | 查找目标颜色 | 区域, 颜色, 相似度, 方向 | (x, y) |
| `FindMultiColors(x1, y1, x2, y2, firstColor, offsetColors, sim, dir)` | 查找多点颜色序列 | 区域, 首点颜色, 偏移颜色数组, 相似度, 方向 | (x, y) |
| `GetColorCountInRegion(x1, y1, x2, y2, color, sim)` | 统计符合条件的像素数 | 区域, 颜色, 相似度 | int |
| `WaitForColor(x1, y1, x2, y2, color, sim, interval, maxRetries, desc)` | 等待颜色出现 | 区域, 颜色, 相似度, 间隔, 最大重试, 描述 | (x, y, bool) |

**颜色格式**: 十六进制字符串（如 `"ba0404"`），不带 `#` 前缀。

### 配置管理器 (ConfigManager)

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Get(path)` | 获取配置值 | 路径（如 "app.packages.kit"） | interface{} |
| `GetString(path)` | 获取字符串配置 | 路径 | string |
| `GetInt(path)` | 获取整数配置 | 路径 | int |
| `GetBool(path)` | 获取布尔配置 | 路径 | bool |
| `GetMap(path)` | 获取 map 配置 | 路径 | map[string]interface{} |
| `GetArray(path)` | 获取数组配置 | 路径 | []interface{} |
| `Exists(path)` | 检查路径是否存在 | 路径 | bool |

**全局实例**: `core.Config`

### WebSocket 客户端 (WebSocketClient)

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Connect()` | 连接到服务器 | 无 | error |
| `Emit(event, data)` | 发送事件 | 事件名, 数据 | error |
| `Disconnect()` | 断开连接 | 无 | error |
| `IsConnected()` | 检查连接状态 | 无 | bool |
| `StartHeartbeat()` | 启动心跳 | 无 | 无 |
| `StopReconnect()` | 停止自动重连 | 无 | 无 |

**自动重连**: 支持 1分钟、3分钟、10分钟间隔的自动重连机制。

**事件处理**:
- `device_auth`: 设备认证
- `screenshot_command`: 截图指令
- `heartbeat`: 心跳响应

### 设备操作 (motion.go)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `Click(x, y)` | 点击指定坐标 | x, y | 无 |
| `RandomClickInArea(x1, y1, x2, y2)` | 在区域内随机点击 | 区域 | 无 |
| `Swipe(x1, y1, x2, y2, duration)` | 滑动操作 | 起点, 终点, 时长 | 无 |
| `Sleep(ms)` | 延迟 | 毫秒 | 无 |
| `RandomSleep(min, max)` | 随机延迟 | 最小值, 最大值 | 无 |
| `Toast(context)` | 显示提示信息 | 内容 | 无 |
| `CloseAllWindows()` | 关闭所有弹窗 | 无 | 无 |
| `Log(message)` | 记录日志 | 消息 | 无 |

## 关键依赖与配置

### 依赖项

- **AutoGo**: 本地依赖库，提供底层自动化能力
  - `github.com/Dasongzi1366/AutoGo/images`
  - `github.com/Dasongzi1366/AutoGo/motion`
  - `github.com/Dasongzi1366/AutoGo/opencv`
  - `github.com/Dasongzi1366/AutoGo/ppocr`
  - `github.com/Dasongzi1366/AutoGo/imgui`
  - `github.com/Dasongzi1366/AutoGo/storages`

- **TomatoOCR**: OCR 客户端
- **gorilla/websocket**: WebSocket 通信
- **aws-sdk-go-v2**: R2 存储（用于截图上传）

### 配置文件

- `assets/config/default.json`: 默认配置
  - `app_packages`: 应用包名配置
  - `ocr.license_key`: OCR 授权密钥
  - `sls`: 阿里云 SLS 日志配置

## 数据模型

### ConfigManager

```go
type ConfigManager struct {
    data map[string]interface{}
}
```

支持多层路径访问，如 `"app.packages.kit"` 会依次查找 `app` → `packages` → `kit`。

### DataStore

```go
type DataStore struct {
    UserInfo UserInfo                `json:"user_info"`
    TempData map[string]interface{}  `json:"temp_data"`
}
```

全局数据存储，包含用户信息和临时数据。

### WebSocketClient

```go
type WebSocketClient struct {
    serverURL       string
    deviceCode      string
    conn            *websocket.Conn
    connected       bool
    mu              sync.RWMutex
    r2Client        *util.R2Client
    reconnectTry    int
    shouldReconnect bool
}
```

WebSocket 客户端，支持自动重连和截图上传到 R2。

### OCRHandler / OpenCVHandler / ColorHandler

均采用**懒加载模式**，首次调用时初始化底层客户端，使用互斥锁保证线程安全。

## 测试与质量

### 当前测试覆盖

- 手动测试：`test1.go` 包含 OCR 和图像识别测试
- 无自动化单元测试

### 建议补充

1. 单元测试：每个处理器的核心方法
2. Mock 测试：外部依赖（TomatoOCR、WebSocket）
3. 集成测试：多个处理器协同工作的场景

## 常见问题 (FAQ)

### 1. OCR 识别不准确怎么办？

- 检查区域坐标是否正确
- 调整相似度阈值（建议 0.6-0.8）
- 确保 OCR 密钥有效

### 2. 图像识别失败的可能原因？

- 图片路径错误（应使用相对于 `assets/img/` 的路径）
- 相似度阈值过高（建议 0.7-0.9）
- 需要灰度匹配时未设置 `isGray=true`
- 模板缓存问题（调用 `ClearTemplateCache()` 清空）

### 3. WebSocket 连接断开怎么办？

- 框架已实现自动重连机制（1分钟、3分钟、10分钟间隔）
- 检查服务器地址和设备编号是否正确
- 查看日志了解断连原因

### 4. 如何添加新的配置项？

1. 在 `assets/config/default.json` 添加配置
2. 使用 `core.Config.Get()` 或 `GetString()` 等方法读取
3. 支持多层路径，如 `"new_section.new_key"`

### 5. 颜色检测不工作？

- 确保颜色格式为十六进制字符串（如 `"ba0404"`，不带 `#`）
- 调整相似度（建议 0.9-0.98）
- 检查坐标区域是否正确

## 相关文件清单

### 核心源码

- `global.go`: 全局实例定义
- `ocr.go`: OCR 处理器
- `opencv.go`: OpenCV 处理器
- `color.go`: 颜色处理器
- `config.go`: 配置管理器
- `api.go`: API 处理器
- `ws_client.go`: WebSocket 客户端
- `motion.go`: 设备操作

### 配置文件

- `../assets/config/default.json`: 默认配置

### 依赖库

- `../AutoGo/`: 本地依赖库
- `../TomatoOCR/`: OCR 客户端

---

**模块统计**:
- 文件数: 8
- 导出接口数: 50+
- 全局实例数: 4 (OCR, OpenCV, Color, API)
- 线程安全: 是（所有处理器使用互斥锁）
