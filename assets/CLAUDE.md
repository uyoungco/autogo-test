# assets 模块文档

[根目录](../CLAUDE.md) > **assets**

> 本文档由 AI 架构师于 2025-11-15 17:46:19 生成

## 变更记录 (Changelog)

### 2025-11-15 17:46:19
- 初始化模块文档

---

## 模块职责

`assets` 模块负责**静态资源管理**，使用 Go embed 将资源嵌入到编译后的可执行文件中，包括：

1. **图片资源** (img/): 用于图像识别的模板图片
2. **配置文件** (config/): JSON 格式的配置文件
3. **模板缓存**: 图像和遮罩缓存（提高识别性能）

## 入口与启动

资源在包初始化时自动嵌入，无需手动加载。

## 对外接口

### 嵌入文件系统

```go
//go:embed img/*
var ImageFile embed.FS

//go:embed config/*
var ConfigFile embed.FS
```

**使用示例**：
```go
// 读取配置文件
data, err := assets.ConfigFile.ReadFile("config/default.json")

// 读取图片（通常通过 core/opencv.go 自动处理）
imgData, err := assets.ImageFile.ReadFile("img/main/主城导航.png")
```

### 缓存映射

```go
var TemplateMap = make(map[string]opencv.Mat) // 模板图像缓存
var MaskMap = make(map[string]opencv.Mat)     // 遮罩图像缓存
```

**缓存键格式**: `"文件名-是否灰度-缩放比例"`
- 例如: `"main/主城导航.png-false-1.0"`

**缓存管理**: 由 `core/opencv.go` 自动管理，可调用 `core.OpenCV.ClearTemplateCache()` 清空。

## 关键依赖与配置

### 依赖项

- `embed`: Go 标准库，用于嵌入静态资源
- `gocv.io/x/gocv`: OpenCV 绑定（opencv.Mat）

### 目录结构

```
assets/
├── img/                    # 图像资源目录
│   ├── close/              # 关闭按钮图标
│   │   ├── 关闭1.png
│   │   └── 关闭2.png
│   ├── login/              # 登录相关图标
│   │   ├── 游戏更新_资源弹框.png
│   │   ├── 通知弹框.png
│   │   └── 谷歌登录弹窗.png
│   ├── main/               # 主界面图标
│   │   ├── 主城导航.png
│   │   ├── 回程卷轴.png
│   │   ├── 药水商人.png
│   │   ├── 礼包.png
│   │   ├── 退出地图.png
│   │   └── 跳跃.png
│   ├── sys/                # 系统图标
│   │   └── kit_start.png
│   ├── 仓库-哈夫币-icon.png
│   ├── 仓库-三角币-icon.png
│   ├── 仓库_保险箱图标.png
│   ├── 仓库_安全箱_问号.png
│   ├── 界面-仓库.png
│   ├── 主界面-取消行动.png
│   ├── 主界面-出发.png
│   ├── gerenzhongxin.png
│   └── editicon.png
├── config/                 # 配置文件目录
│   ├── default.json        # 默认配置
│   ├── scene.json          # 场景识别配置
│   └── scene_map.json      # 地图场景配置
└── definitions.go          # 资源定义文件
```

## 数据模型

### embed.FS

```go
type FS interface {
    Open(name string) (fs.File, error)
    ReadFile(name string) ([]byte, error)
    ReadDir(name string) ([]fs.DirEntry, error)
}
```

Go embed 提供的文件系统接口，支持读取嵌入的文件。

### opencv.Mat

OpenCV 图像数据结构，由 `gocv.io/x/gocv` 提供。

**重要提示**: Mat 对象需要手动调用 `Close()` 释放内存。

## 资源管理策略

### 图像资源

**命名规范**：
- 使用中文或英文描述性命名
- 按功能分组到子目录（如 `close/`, `main/`, `login/`）
- 支持透明 PNG（用于遮罩识别）

**使用方式**：
```go
// 在 core/opencv.go 中自动加载
x, y := core.OpenCV.FindImage(0, 0, 0, 0, "img/main/主城导航.png", false, 1, 0.8)
```

**图片路径**：
- 相对于 `assets/` 的路径（如 `img/main/主城导航.png`）
- 在配置文件中可省略 `img/` 前缀（如 `main/主城导航.png`）

### 配置文件

**default.json**：
```json
{
  "app_packages": {
    "kit": "fun.kitsunebi.kitsunebi4android",
    "ymir": "com.wemade.ymir",
    "google": "com.google.android.gms"
  },
  "ocr": {
    "license_key": "..."
  },
  "sls": {
    "access_key_id": "...",
    "endpoint": "..."
  }
}
```

**scene.json**：
```json
{
  "主界面": {
    "region": "0,0,0,0",
    "images": "main/主城导航.png",
    "sim": 0.8
  }
}
```

**scene_map.json**：
```json
{
  "阿斯加德城": {
    "regionType": "map_name",
    "text": "아스가르드",
    "sim": 0.8
  }
}
```

### 缓存管理

**缓存策略**：
- 首次加载图片时自动缓存到 `TemplateMap`
- 透明图片的遮罩缓存到 `MaskMap`
- 缓存键包含文件名、灰度标志、缩放比例

**清空缓存**：
```go
core.OpenCV.ClearTemplateCache()
```

**适用场景**：
- 内存不足时
- 图片文件被更新后
- 识别异常时

## 测试与质量

### 当前测试覆盖

- 无自动化测试
- 通过实际运行验证资源加载

### 建议补充

1. 单元测试：资源文件是否存在
2. 集成测试：图像识别准确性
3. 性能测试：缓存命中率

## 常见问题 (FAQ)

### 1. 如何添加新图片？

1. 将图片放到 `assets/img/` 相应子目录
2. 重新编译项目（Go embed 会自动包含新文件）
3. 在代码中使用相对路径引用

### 2. 图片识别失败怎么办？

- 检查图片路径是否正确
- 检查图片质量（清晰度、大小）
- 调整相似度阈值
- 尝试灰度匹配
- 调整缩放比例

### 3. 配置文件如何修改？

修改 `assets/config/` 下的 JSON 文件，重新编译项目。

**注意**: 运行时修改配置文件无效，需要重新编译。

### 4. 如何查看嵌入的资源？

嵌入的资源在编译后无法直接查看，建议保留源文件用于开发调试。

### 5. 透明图片如何处理？

- 保存为 PNG 格式
- 使用透明通道（Alpha）
- OpenCV 会自动生成遮罩并缓存到 `MaskMap`

### 6. 缓存占用内存过多怎么办？

定期调用 `core.OpenCV.ClearTemplateCache()` 清空缓存。

### 7. 如何组织大量图片？

- 按功能分组到子目录
- 使用描述性命名
- 避免重复图片
- 定期清理未使用的图片

## 相关文件清单

### 核心源码

- `definitions.go`: 资源定义

### 资源目录

- `img/`: 图像资源
  - `close/`: 关闭按钮（2 张）
  - `login/`: 登录相关（3 张）
  - `main/`: 主界面（6 张）
  - `sys/`: 系统图标（1 张）
  - 其他图标（9 张）
- `config/`: 配置文件（3 个）

### 相关模块

- `../core/`: 核心模块（使用图像和配置资源）
- `../scene/`: 场景识别（使用配置资源）

---

**模块统计**:
- 文件数: 1 (definitions.go)
- 图片数: 21+
- 配置文件数: 3
- 缓存映射数: 2 (TemplateMap, MaskMap)
