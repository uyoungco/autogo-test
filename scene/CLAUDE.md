# scene 模块文档

[根目录](../CLAUDE.md) > **scene**

> 本文档由 AI 架构师于 2025-11-15 17:46:19 生成

## 变更记录 (Changelog)

### 2025-11-15 17:46:19
- 初始化模块文档

---

## 模块职责

`scene` 模块负责**游戏场景识别**，基于以下方式判断当前场景：

1. **图像识别**: 识别特定图标或界面元素
2. **OCR 文字识别**: 识别界面文字（支持韩文）
3. **组合判断**: 支持 `&`（且）和 `|`（或）操作符

核心能力：**高性能场景识别 + OCR 缓存优化**。

## 入口与启动

### 自动初始化

场景管理器在包初始化时自动加载配置：

```go
// scene/identify.go init()
func init() {
    manager = &SceneManager{...}
    loadSceneConfig() // 加载 scene.json
}
```

### 场景识别

```go
// 使用默认配置（scene.json）
sceneName := scene.Identify()

// 使用地图配置（scene_map.json）
mapName := scene.Identify("scene_map")
```

## 对外接口

### Identify (场景识别)

```go
func Identify(configType ...string) string
```

**参数**：
- `configType`（可选）：配置类型（`"scene"` 或 `"scene_map"`），默认为 `"scene"`

**返回值**：
- 场景名称（如 `"主界面"`、`"阿斯加德城"`）
- 未识别时返回 `"未知场景"`

**使用示例**：
```go
// 识别当前场景
sceneName := scene.Identify()
if sceneName == "主界面" {
    // 执行主界面操作
}

// 识别当前地图
mapName := scene.Identify("scene_map")
if mapName == "阿斯加德城" {
    // 在阿斯加德城
}
```

## 关键依赖与配置

### 依赖项

- `app/assets`: 静态资源（配置文件）
- `app/core`: 核心模块（OCR、OpenCV）
- `encoding/json`: JSON 解析
- `sync`: 线程安全（互斥锁）

### 配置文件

**场景配置** (`assets/config/scene.json`):
```json
{
  "主界面": {
    "region": "0,0,0,0",
    "images": "main/主城导航.png",
    "isGray": false,
    "sim": 0.8,
    "scale": 1.0
  },
  "界面仓库": {
    "regionType": "map_name",
    "text": "창고",
    "sim": 0.7
  }
}
```

**地图配置** (`assets/config/scene_map.json`):
```json
{
  "阿斯加德城": {
    "regionType": "map_name",
    "text": "아스가르드",
    "sim": 0.8
  }
}
```

### 配置字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `region` | string | 否 | 检测区域 `"x1,y1,x2,y2"`，默认全屏 `"0,0,0,0"` |
| `regionType` | string | 否 | 区域类型（如 `"map_name"`），优先于 `region` |
| `images` | string | 条件 | 图片路径，支持 `\|`（或）和 `&`（且）操作符 |
| `text` | string | 条件 | OCR 文本识别内容 |
| `isGray` | bool | 否 | 是否灰度匹配，默认 `false` |
| `sim` | float | 否 | 相似度阈值，默认 `0.8` |
| `scale` | float | 否 | 缩放比例，默认 `1.0` |

**注意**: `images` 和 `text` 至少需要提供一个。

### 预定义区域类型

```go
var regionTypes = map[string]string{
    "map_name": "57,14,188,40", // 地图名区域
}
```

可在 `identify.go` 中添加更多预定义区域。

## 数据模型

### SceneConfig

```go
type SceneConfig struct {
    Region     string  // 检测区域 "x1,y1,x2,y2"
    RegionType string  // 区域类型（优先于 region）
    Images     string  // 图片名称，支持 | 和 & 操作符
    Text       string  // OCR 文本识别内容
    IsGray     bool    // 是否灰度匹配
    Sim        float32 // 相似度阈值
    Scale      float32 // 缩放比例
}
```

### SceneManager

```go
type SceneManager struct {
    scenes   map[string]SceneConfig // 默认场景配置
    sceneMap map[string]SceneConfig // 地图场景配置
    opencv   *core.OpenCVHandler
    ocr      *core.OCRHandler
    mutex    sync.Mutex             // 互斥锁
    ocrCache map[string]string      // OCR 结果缓存
}
```

**OCR 缓存优化**：
- 同一区域的 OCR 结果会被缓存
- 每次 `Identify()` 调用前自动清空缓存
- 避免同一帧内重复 OCR 识别

### OCRCache

```go
type OCRCache struct {
    Region string // 区域坐标
    Text   string // OCR 识别结果
}
```

## 场景识别机制

### 图像识别模式

**单个图片**：
```json
{
  "主界面": {
    "images": "main/主城导航.png"
  }
}
```

**或操作（`|`）**：只要找到其中一个即可
```json
{
  "关闭弹窗": {
    "images": "close/关闭1.png|close/关闭2.png"
  }
}
```

**且操作（`&`）**：所有图片都必须找到
```json
{
  "特定界面": {
    "images": "icon1.png&icon2.png&icon3.png"
  }
}
```

### OCR 识别模式

```json
{
  "界面仓库": {
    "regionType": "map_name",
    "text": "창고",
    "sim": 0.7
  }
}
```

**文本相似度计算**：
- 使用 **Levenshtein 编辑距离算法**
- 自动清理地图名（去掉数字及后面的字符，如 `폴크방1층` → `폴크방`）
- 去除空格并转换为小写，提高匹配准确性

### 混合模式

可以同时使用图像和文本识别：
```json
{
  "复杂场景": {
    "images": "icon.png",
    "text": "标题文字",
    "sim": 0.8
  }
}
```

**注意**: 当前实现中，如果同时配置了 `text` 和 `images`，优先使用 `text`（OCR 模式）。

## 测试与质量

### 当前测试覆盖

- 无自动化测试
- 通过实际运行验证场景识别准确性

### 建议补充

1. 单元测试：文本相似度计算、区域解析
2. Mock 测试：模拟 OCR 和 OpenCV 结果
3. 集成测试：完整场景识别流程

## 常见问题 (FAQ)

### 1. 场景识别失败怎么办？

- 检查配置文件中的图片路径是否正确
- 检查相似度阈值是否合适（建议 0.7-0.9）
- 检查区域坐标是否正确
- 查看 OCR 缓存是否正确（可能需要重启）

### 2. 如何添加新场景？

1. 在 `assets/config/scene.json` 添加场景配置
2. 使用 `scene.Identify()` 识别
3. 无需修改代码

### 3. OCR 识别韩文不准确？

- 调整相似度阈值（建议 0.6-0.8）
- 检查检测区域是否包含目标文字
- 确保 OCR 密钥有效

### 4. 如何提高识别性能？

- 使用 `regionType` 限定检测区域
- 优先使用图像识别（比 OCR 快）
- OCR 结果会自动缓存，同一帧内不会重复识别

### 5. 如何调试场景识别？

```go
sceneName := scene.Identify()
core.Log("当前场景: " + sceneName)
core.Toast("当前场景: " + sceneName)
```

### 6. 图像或操作符（`|`）的优先级？

从左到右依次检测，找到第一个匹配的即返回 `true`。

### 7. 图像且操作符（`&`）的顺序重要吗？

重要。从左到右依次检测，任意一个未找到则返回 `false`。建议将最可能失败的放在前面，提高性能。

## 相关文件清单

### 核心源码

- `identify.go`: 场景识别实现

### 配置文件

- `../assets/config/scene.json`: 默认场景配置
- `../assets/config/scene_map.json`: 地图场景配置

### 相关模块

- `../core/`: 核心模块（OCR、OpenCV）
- `../assets/`: 静态资源

---

**模块统计**:
- 文件数: 1
- 导出接口数: 1 (Identify)
- 配置文件数: 2
- 线程安全: 是（使用互斥锁）
- OCR 缓存: 是
