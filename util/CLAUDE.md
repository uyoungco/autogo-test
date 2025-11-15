# util 模块文档

[根目录](../CLAUDE.md) > **util**

> 本文档由 AI 架构师于 2025-11-15 17:46:19 生成

## 变更记录 (Changelog)

### 2025-11-15 (更新)
- 新增 `inventory_detector.go`: 仓库物品检测与可视化工具
- 新增 `InventoryItem` 和 `ColorRange` 数据类型
- 新增 `FindByColor()`: 基于 HSV 颜色范围的物品识别
- 新增 `VisualizeResults()`: 物品检测结果可视化
- 新增 `CalculateOccupiedCells()` 和 `CreateCellArray()`: 格数计算工具
- 完善 `imageHandler.go` 文档说明

### 2025-11-15 17:46:19
- 初始化模块文档

---

## 模块职责

`util` 模块提供**通用工具函数**，包括：

1. **HTTP 请求** (http.go): 封装 HTTP GET/POST 请求
2. **R2 存储** (r2.go): Cloudflare R2 对象存储客户端
3. **坐标处理** (coordinate.go): 坐标计算和转换
4. **图像处理** (imageHandler.go): 图像格式转换工具
5. **仓库物品检测** (inventory_detector.go): 基于颜色的物品识别与可视化
6. **数组工具** (contains.go): 数组包含判断

## 入口与启动

工具函数无需初始化，直接调用即可。

## 对外接口

### HTTP 请求 (HttpRequest)

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `PostJSON(url, data, headers, timeout)` | 发送 POST 请求（JSON） | URL, 数据, 请求头, 超时 | (响应数据, 状态码, error) |
| `Get(url, headers, timeout)` | 发送 GET 请求 | URL, 请求头, 超时 | (响应数据, 状态码, error) |

**使用示例**：
```go
data := map[string]interface{}{
    "deviceCode": "device123",
    "accountId": "user456",
}
response, statusCode, err := util.HttpRequest.PostJSON(
    "http://example.com/api",
    data,
    nil,
    30000, // 30 秒超时
)
```

### R2 存储客户端 (R2Client)

#### 创建客户端

```go
client, err := util.NewR2Client(util.R2Config{
    AccountID:       "your-account-id",
    AccessKeyID:     "your-access-key-id",
    AccessKeySecret: "your-access-key-secret",
    BucketName:      "your-bucket-name",
    PublicDomain:    "https://your-domain.com", // 可选
})
```

#### 主要方法

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `UploadFile(filePath, data, contentType)` | 上传文件 | 路径, 数据, MIME 类型 | (URL, error) |
| `UploadFileWithTimestamp(prefix, filename, data, contentType)` | 上传文件（自动添加时间戳） | 前缀, 文件名, 数据, MIME 类型 | (URL, error) |
| `UploadImage(filePath, img, format, quality)` | 上传图片 | 路径, 图片, 格式, 质量 | (URL, error) |
| `UploadImageWithTimestamp(prefix, filename, img, format, quality)` | 上传图片（自动添加时间戳） | 前缀, 文件名, 图片, 格式, 质量 | (URL, error) |
| `DeleteFile(filePath)` | 删除文件 | 路径 | error |
| `ListFiles(prefix)` | 列出文件 | 前缀 | ([]string, error) |
| `FileExists(filePath)` | 检查文件是否存在 | 路径 | (bool, error) |

**使用示例**：
```go
// 上传截图（JPEG 格式，质量 70）
screenshot := images.CaptureScreen(0, 0, 0, 0, 0)
url, err := client.UploadImageWithTimestamp(
    "screenshots/",
    "device123_screen.jpg",
    screenshot,
    "jpeg",
    70,
)
```

### 坐标工具 (coordinate.go)

*具体功能待补充*

### 图像处理 (imageHandler.go)

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `ImageToMat(img)` | 将 *image.NRGBA 转换为 opencv.Mat (BGR) | 图像 | opencv.Mat |
| `MatToImage(mat)` | 将 opencv.Mat (BGR) 转换为 *image.NRGBA | Mat | *image.NRGBA |
| `MaskToImage(mask)` | 将单通道掩膜转换为 *image.NRGBA | 掩膜 | *image.NRGBA |
| `Mat3ToImage(mat)` | 将三通道 Mat 转换为 *image.NRGBA | Mat | *image.NRGBA |
| `CreateHSVScalars(hsvRange)` | 创建 HSV 上下限 Scalar | HSV 范围 | (lower, upper) |

**使用示例**：
```go
// NRGBA 图像转换为 Mat
img := images.CaptureScreen(0, 0, 1920, 1080, 0)
mat := util.ImageToMat(img)
defer mat.Close()

// Mat 转换回 NRGBA 图像
nrgba := util.MatToImage(mat)
images.Save(nrgba, "/sdcard/output.png", 100)
```

### 仓库物品检测 (inventory_detector.go)

#### 数据类型

**InventoryItem**: 仓库物品信息
```go
type InventoryItem struct {
    X, Y             int    // 左上角坐标
    W, H             int    // 宽度和高度
    CenterX, CenterY int    // 中心点坐标
    Size             int    // 占用格数 (1-9)
    Color            string // 颜色标签
}
```

**ColorRange**: 颜色检测范围配置
```go
type ColorRange struct {
    Name  string   // 颜色名称标签（如"黄色"）
    HSV   HSVRange // HSV 范围
    Color string   // 颜色英文标识（如"yellow"）
}
```

#### 主要函数

| 函数 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `FindByColor(img, colorRanges, minSize, cellSize)` | 按 HSV 范围筛选轮廓并识别物品 | 图像, 颜色范围列表, 最小尺寸, 单元格大小 | ([]InventoryItem, error) |
| `VisualizeResults(img, items)` | 在图上绘制矩形与编号 | 图像, 物品列表 | (opencv.Mat, error) |
| `CalculateOccupiedCells(x, y, w, h, cellSize)` | 依据矩形面积匹配占用格数 | 坐标和尺寸, 单元格大小 | int (1-9, 0表示无匹配) |
| `CreateCellArray(cellSize, length, tolerance)` | 生成面积区间数组 | 单元格大小, 数量, 容差 | [][2]float64 |

**使用示例**：
```go
// 定义颜色范围
colorRanges := []util.ColorRange{
    {
        Name: "黄色",
        HSV: util.HSVRange{
            LowerH: 11, LowerS: 45, LowerV: 58,
            UpperH: 14, UpperS: 66, UpperV: 77,
        },
        Color: "yellow",
    },
    {
        Name: "绿色",
        HSV: util.HSVRange{
            LowerH: 84, LowerS: 58, LowerV: 40,
            UpperH: 88, UpperS: 91, UpperV: 63,
        },
        Color: "green",
    },
}

// 截取屏幕区域
img := images.CaptureScreen(570, 78, 1201, 638, 0)

// 查找物品（最小尺寸 65 像素，单元格大小 68 像素）
items, err := util.FindByColor(img, colorRanges, 65, 68)
if err != nil {
    log.Fatal(err)
}

// 打印物品信息
for _, item := range items {
    fmt.Println(item) // 输出: InventoryItem(x:100, y:200, w:70, h:70, centerX:135, centerY:235, size:1, color:黄色)
}

// 可视化结果
vis, err := util.VisualizeResults(img, items)
if err != nil {
    log.Fatal(err)
}
defer vis.Close()

// 保存可视化结果
visImg := util.MatToImage(vis)
images.Save(visImg, "/sdcard/result.png", 100)
```

**工作原理**：
1. 将输入图像从 BGR 转换为 HSV 色彩空间
2. 对每种颜色范围生成掩膜（mask）
3. 查找掩膜中的外部轮廓
4. 过滤小于最小尺寸的轮廓
5. 根据轮廓面积计算占用格数（容差 5%）
6. 返回符合条件的物品列表

**注意事项**：
- HSV 范围需要根据实际游戏截图调整
- `minSize` 用于过滤噪点和小图标
- `cellSize` 是游戏中单个格子的像素大小
- `VisualizeResults` 返回的 Mat 需要调用方手动 `Close()`

### 数组工具 (contains.go)

*具体功能待补充*

## 关键依赖与配置

### 依赖项

- **AWS SDK for Go v2**: 用于 R2 存储
  - `github.com/aws/aws-sdk-go-v2/aws`
  - `github.com/aws/aws-sdk-go-v2/config`
  - `github.com/aws/aws-sdk-go-v2/credentials`
  - `github.com/aws/aws-sdk-go-v2/service/s3`
- **标准库**:
  - `net/http`: HTTP 请求
  - `image`: 图像处理
  - `image/jpeg`: JPEG 编码
  - `image/png`: PNG 编码

### 配置参数

**R2 配置**：
```go
type R2Config struct {
    AccountID       string // Cloudflare 账户 ID
    AccessKeyID     string // R2 访问密钥 ID
    AccessKeySecret string // R2 访问密钥
    BucketName      string // 存储桶名称
    PublicDomain    string // 公开访问域名（可选）
}
```

## 数据模型

### R2Client

```go
type R2Client struct {
    client       *s3.Client
    bucketName   string
    publicDomain string // 公开访问域名
}
```

**端点配置**：
- R2 端点: `https://{AccountID}.r2.cloudflarestorage.com`
- 公开访问域名: 自定义域名或 `https://pub-{AccountID}.r2.dev`

### R2Config

详见"配置参数"章节。

## R2 存储详解

### 上传流程

1. **创建客户端**: 配置账户 ID、密钥、存储桶名称
2. **准备数据**: 图片或文件数据
3. **上传**: 调用 `UploadFile()` 或 `UploadImage()`
4. **获取 URL**: 返回公开访问 URL

### 文件路径

**手动指定路径**：
```go
url, err := client.UploadFile("screenshots/2025/screen.jpg", data, "image/jpeg")
```

**自动添加时间戳**：
```go
// 生成路径: screenshots/20250115_174619_device123_screen.jpg
url, err := client.UploadFileWithTimestamp(
    "screenshots/",
    "device123_screen.jpg",
    data,
    "image/jpeg",
)
```

### 图片格式

支持的格式：
- **PNG**: 无损压缩，文件较大
- **JPEG**: 有损压缩，可调整质量（1-100）

**推荐配置**：
- 截图: JPEG 格式，质量 70
- 图标: PNG 格式

### 权限设置

R2 存储桶需要配置为**公开访问**或提供公开访问域名：
1. 在 Cloudflare 控制台启用公开访问
2. 配置自定义域名（推荐）
3. 使用 `PublicDomain` 参数指定域名

## 测试与质量

### 当前测试覆盖

- 无自动化测试
- `r2_example.go` 提供使用示例

### 建议补充

1. 单元测试：HTTP 请求、R2 上传/删除
2. Mock 测试：模拟 S3 客户端
3. 集成测试：真实 R2 存储桶操作

## 常见问题 (FAQ)

### 1. R2 上传失败怎么办？

- 检查账户 ID、密钥是否正确
- 检查存储桶名称是否存在
- 检查网络连接
- 查看错误信息了解具体原因

### 2. 如何获取 R2 配置信息？

1. 登录 Cloudflare 控制台
2. 进入 R2 存储桶页面
3. 创建 API 令牌获取密钥
4. 记录账户 ID 和存储桶名称

### 3. 上传的文件如何访问？

通过返回的 URL 直接访问，格式为：
```
https://{PublicDomain}/{filePath}
```

### 4. 如何设置文件过期时间？

当前实现不支持，需要在 Cloudflare 控制台配置**生命周期规则**。

### 5. HTTP 请求超时怎么办？

调整 `timeout` 参数（单位：毫秒）：
```go
response, statusCode, err := util.HttpRequest.PostJSON(
    url,
    data,
    nil,
    60000, // 60 秒
)
```

### 6. 如何上传大文件？

当前实现使用 `bytes.Reader`，适合中小文件（< 100MB）。大文件建议使用分片上传（需扩展实现）。

## 相关文件清单

### 核心源码

- `http.go`: HTTP 请求工具
- `r2.go`: R2 存储客户端
- `r2_example.go`: R2 使用示例
- `coordinate.go`: 坐标工具
- `imageHandler.go`: 图像格式转换（NRGBA ↔ Mat, HSV Scalar 创建）
- `inventory_detector.go`: 仓库物品检测与可视化
- `contains.go`: 数组工具

### 相关模块

- `../core/`: 核心模块（WebSocket 客户端使用 R2）

---

**模块统计**:
- 文件数: 7
- 导出接口数: 20+
- 主要功能: HTTP 请求, R2 存储, 仓库物品检测, 图像转换
