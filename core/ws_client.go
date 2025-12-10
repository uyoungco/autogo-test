package core

import (
	"app/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/gorilla/websocket"
)

// WebSocket 心跳超时配置
const (
	pongWait   = 60 * time.Second // 等待 pong 消息的超时时间
	pingPeriod = 30 * time.Second // 发送 ping 消息的间隔（必须小于 pongWait）
)

// WebSocketClient WebSocket 客户端结构体
type WebSocketClient struct {
	serverURL       string          // 服务器地址，如 http://localhost:3000
	deviceCode      string          // 设备唯一标识码，用于认证
	conn            *websocket.Conn // WebSocket 连接实例
	connected       bool            // 连接状态标志，true 表示已连接
	mu              sync.RWMutex    // 读写锁
	r2Client        *util.R2Client
	reconnectTry    int           // 当前重连尝试次数（0, 1, 2）
	shouldReconnect bool          // 是否应该重连
	done            chan struct{} // 用于通知 goroutine 退出
	pingTicker      *time.Ticker  // ping 定时器
}

// DeviceAuthRequest 设备认证请求
type DeviceAuthRequest struct {
	DeviceCode string `json:"deviceCode"` // 设备唯一标识码
}

// DeviceAuthResponse 设备认证响应
type DeviceAuthResponse struct {
	Success bool   `json:"success"` // 认证是否成功
	Message string `json:"message"` // 认证结果消息
}

// HeartbeatRequest 心跳请求
type HeartbeatRequest struct {
	Timestamp int64 `json:"timestamp"` // 客户端发送心跳的时间戳（毫秒）
}

// HeartbeatResponse 心跳响应
type HeartbeatResponse struct {
	Success    bool  `json:"success"`    // 心跳是否成功
	ServerTime int64 `json:"serverTime"` // 服务器时间戳（毫秒）
}

// ScreenshotDataRequest 截图数据请求
type ScreenshotDataRequest struct {
	ScreenshotBase64 string `json:"screenshotBase64"` // Base64 编码的截图数据
	Timestamp        int64  `json:"timestamp"`        // 截图时间戳（毫秒）
}

// ScreenshotDataResponse 截图数据响应
type ScreenshotDataResponse struct {
	Success bool   `json:"success"` // 截图数据是否成功保存
	Message string `json:"message"` // 处理结果消息
}

// ScreenshotURLRequest 截图 URL 请求
type ScreenshotURLRequest struct {
	ScreenshotURL string `json:"screenshotUrl"` // 截图的 R2 存储 URL
	Timestamp     int64  `json:"timestamp"`     // 截图时间戳（毫秒）
}

// ScreenshotCommand 截图指令
type ScreenshotCommand struct {
	Timestamp int64 `json:"timestamp"` // 服务器发送指令的时间戳（毫秒）
	Timeout   int64 `json:"timeout"`   // 截图超时时间（毫秒）
}

// ErrorRequest 错误请求
type ErrorRequest struct {
	Error   string `json:"error"`             // 错误信息
	Details string `json:"details,omitempty"` // 错误详细信息（可选）
}

// NewWebSocketClient 创建新的 WebSocket 客户端
func NewWebSocketClient(serverURL, deviceCode string) *WebSocketClient {

	// 1. 创建 R2 客户端
	client, err := util.NewR2Client(util.R2Config{
		AccountID:       "227d58ddf76b97d47968d3443e1aa726",
		AccessKeyID:     "c656346d615b67abcb3f73fc6365bf17",
		AccessKeySecret: "4a3eedb529b778583d6d2ea13b7fc214a858775985891e3012ca3f6fcba95a55",
		BucketName:      "haval-coin",
		PublicDomain:    "https://haval-coin-img.uyoung.co", // R2 公开访问域名
	})
	if err != nil {
		log.Fatalf("创建 R2 客户端失败: %v", err)
	}

	return &WebSocketClient{
		r2Client:        client,
		serverURL:       serverURL,
		deviceCode:      deviceCode,
		connected:       false,
		reconnectTry:    0,
		shouldReconnect: true,
	}
}

// Connect 连接到 WebSocket 服务器
func (c *WebSocketClient) Connect() error {
	log.Println("========== 开始连接 WebSocket 服务器 ==========")

	// 构建 WebSocket URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		log.Printf("❌ 解析服务器 URL 失败: %v", err)
		return fmt.Errorf("解析服务器 URL 失败: %v", err)
	}

	// 将 http/https 转换为 ws/wss
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	// 添加 Socket.IO 路径
	u.Path = "/socket.io/"

	// 添加查询参数
	q := u.Query()
	q.Set("EIO", "4")
	q.Set("transport", "websocket")
	u.RawQuery = q.Encode()

	log.Printf("📤 连接到: %s", u.String())

	// 建立 WebSocket 连接
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(u.String(), http.Header{})
	if err != nil {
		log.Printf("❌ WebSocket 连接失败: %v", err)
		return fmt.Errorf("WebSocket 连接失败: %v", err)
	}

	// 设置最大消息大小为 100MB，以支持大型截图数据的接收
	conn.SetReadLimit(100 * 1024 * 1024) // 100MB

	// 初始化心跳相关资源
	c.mu.Lock()
	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}
	c.signalDone() // 关闭旧的 done channel
	c.done = make(chan struct{})
	c.pingTicker = time.NewTicker(pingPeriod)
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	log.Println("✅ WebSocket 连接建立成功")

	// 启动消息接收循环
	go c.readLoop()
	// 启动 WebSocket ping 循环
	go c.pingLoop()

	// 等待接收服务器的握手消息
	time.Sleep(500 * time.Millisecond)

	// 发送 Socket.IO 连接请求 (消息类型 40)
	log.Println("📤 发送 Socket.IO 连接请求...")
	c.mu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, []byte("40"))
	c.mu.Unlock()
	if err != nil {
		log.Printf("❌ 发送连接请求失败: %v", err)
		return fmt.Errorf("发送连接请求失败: %v", err)
	}
	log.Println("✓ Socket.IO 连接请求已发送")

	// 等待连接确认
	time.Sleep(500 * time.Millisecond)

	// 发送设备认证
	log.Println("开始设备认证...")
	if err := c.authenticate(); err != nil {
		log.Printf("❌ 设备认证失败: %v", err)
		return fmt.Errorf("设备认证失败: %v", err)
	}

	log.Println("✅ WebSocket 客户端连接成功")
	log.Println("========== WebSocket 连接完成 ==========")
	return nil
}

// readLoop 读取消息循环
func (c *WebSocketClient) readLoop() {
	log.Println("📨 消息接收循环已启动")

	// 设置读超时和 pong 处理器
	c.mu.Lock()
	if c.conn != nil {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		// 仅设置 PongHandler，保持默认 PingHandler（自动回复 pong）
		c.conn.SetPongHandler(func(string) error {
			log.Println("💓 收到服务器 WebSocket pong，刷新读超时")
			return c.conn.SetReadDeadline(time.Now().Add(pongWait))
		})
	}
	c.mu.Unlock()

	defer func() {
		log.Println("📨 消息接收循环已停止")
		c.signalDone() // 通知其他 goroutine 退出
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
		if c.conn != nil {
			c.conn.Close()
		}

		// 触发自动重连
		if c.shouldReconnect {
			go c.autoReconnect()
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// 检查是否是连接关闭错误
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("⚠️ WebSocket 连接已关闭，准备重连...")
				return
			}

			// 检查是否是意外的连接关闭
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("⚠️ WebSocket 连接意外断开，准备重连...")
				return
			}

			// 其他错误只记录日志，继续接收
			log.Printf("⚠️ 读取消息时出现错误（继续运行）: %v", err)
			continue
		}

		// 连接正常，重置重连计数
		c.mu.Lock()
		c.reconnectTry = 0
		c.mu.Unlock()

		// 每次收到消息时刷新读超时
		c.mu.Lock()
		if c.conn != nil {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
		}
		c.mu.Unlock()

		log.Printf("📥 收到消息: %s", string(message))
		c.handleMessage(message)
	}
}

// handleMessage 处理接收到的消息
func (c *WebSocketClient) handleMessage(message []byte) {
	if len(message) == 0 {
		return
	}

	msgStr := string(message)

	// 处理 Socket.IO ping 心跳消息 (2)
	// 服务器发送 "2"，客户端需要回复 "3" (pong)
	if msgStr == "2" {
		log.Println("💓 收到服务器心跳 ping (2)，回复 pong (3)")
		c.mu.Lock()
		err := c.conn.WriteMessage(websocket.TextMessage, []byte("3"))
		c.mu.Unlock()
		if err != nil {
			log.Printf("❌ 回复心跳失败: %v (继续运行)", err)
		} else {
			log.Println("✓ 心跳 pong 已发送")
		}
		return
	}

	// 处理 Socket.IO 事件消息 (42)
	if len(msgStr) >= 2 && msgStr[0:2] == "42" {
		jsonData := msgStr[2:]
		var eventData []json.RawMessage
		if err := json.Unmarshal([]byte(jsonData), &eventData); err != nil {
			log.Printf("❌ 解析消息失败: %v (继续运行)", err)
			return
		}

		if len(eventData) < 2 {
			return
		}

		var eventName string
		if err := json.Unmarshal(eventData[0], &eventName); err != nil {
			log.Printf("❌ 解析事件名称失败: %v (继续运行)", err)
			return
		}

		log.Printf("📥 收到事件: %s", eventName)

		switch eventName {
		case "device_auth":
			c.handleDeviceAuthResponse(eventData[1])
		case "screenshot_command":
			c.handleScreenshotCommandData(eventData[1])
		case "heartbeat":
			c.handleHeartbeatResponse(eventData[1])
		}
	}
}

// Emit 发送事件
func (c *WebSocketClient) Emit(event string, data interface{}) error {
	c.mu.RLock()
	if !c.connected || c.conn == nil {
		c.mu.RUnlock()
		return fmt.Errorf("未连接到服务器")
	}
	c.mu.RUnlock()

	payload := []interface{}{event, data}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	message := "42" + string(jsonData)
	log.Printf("📤 发送消息: %s", message)

	c.mu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	c.mu.Unlock()

	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	return nil
}

// handleDeviceAuthResponse 处理设备认证响应
func (c *WebSocketClient) handleDeviceAuthResponse(data []byte) {
	var response DeviceAuthResponse
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("❌ 解析认证响应失败: %v", err)
		return
	}

	if response.Success {
		log.Printf("✅ 设备认证成功: %s", response.Message)
	} else {
		log.Printf("❌ 设备认证失败: %s", response.Message)
	}
}

// handleHeartbeatResponse 处理心跳响应
func (c *WebSocketClient) handleHeartbeatResponse(data []byte) {
	var response HeartbeatResponse
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("❌ 解析心跳响应失败: %v", err)
		return
	}

	if response.Success {
		log.Printf("💓 心跳成功，服务器时间: %d", response.ServerTime)
	}
}

// authenticate 设备认证
func (c *WebSocketClient) authenticate() error {
	log.Printf("正在进行设备认证: %s", c.deviceCode)

	request := DeviceAuthRequest{
		DeviceCode: c.deviceCode,
	}

	log.Printf("📤 发送认证请求: %+v", request)

	// 发送认证请求
	err := c.Emit("device_auth", request)
	if err != nil {
		return fmt.Errorf("发送认证请求失败: %v", err)
	}

	log.Println("✓ 设备认证请求已发送")
	return nil
}

// StartHeartbeat 启动心跳
func (c *WebSocketClient) StartHeartbeat() {
	ticker := time.NewTicker(60 * time.Second) // 每60秒发送一次心跳
	defer ticker.Stop()

	log.Println("💓 心跳循环已启动")

	for range ticker.C {
		if !c.connected {
			log.Println("⚠️ 未连接，跳过心跳")
			continue
		}

		if err := c.sendHeartbeat(); err != nil {
			log.Printf("❌ 发送心跳失败: %v", err)
		}
	}
}

// sendHeartbeat 发送心跳
func (c *WebSocketClient) sendHeartbeat() error {
	request := HeartbeatRequest{
		Timestamp: time.Now().UnixMilli(),
	}

	log.Printf("💓 发送心跳: %+v", request)

	// 发送心跳请求
	err := c.Emit("heartbeat", request)
	if err != nil {
		return fmt.Errorf("发送心跳请求失败: %v", err)
	}

	return nil
}

// handleScreenshotCommandData 处理截图指令数据
func (c *WebSocketClient) handleScreenshotCommandData(data []byte) {
	log.Println("📸 收到截图指令")

	var command ScreenshotCommand
	if err := json.Unmarshal(data, &command); err != nil {
		log.Printf("❌ 解析截图指令失败: %v", err)
		return
	}

	log.Printf("截图指令详情 - 时间戳: %d, 超时: %d ms", command.Timestamp, command.Timeout)

	// 执行截图操作
	go c.takeScreenshot()
}

// takeScreenshot 执行截图并上传到 R2
func (c *WebSocketClient) takeScreenshot() {
	log.Println("📸 开始执行截图...")

	// 1. 截取屏幕
	screenshot := images.CaptureScreen(0, 0, 0, 0, 0)
	if screenshot == nil {
		log.Println("❌ 截图失败：截图数据为空")
		c.sendError("截图失败", "截图数据为空")
		return
	}

	// 2. 上传截图到 R2（使用 JPEG 格式，质量 70）
	log.Println("📤 正在上传截图到 R2...")
	// 文件名格式：设备编号_screen.jpg（时间戳会自动添加）
	filename := fmt.Sprintf("%s_screen.jpg", c.deviceCode)
	imageURL, err := c.r2Client.UploadImageWithTimestamp("screenshots/", filename, screenshot, "jpeg", 70)
	if err != nil {
		log.Printf("❌ 上传截图到 R2 失败: %v", err)
		c.sendError("上传截图失败", err.Error())
		return
	}

	log.Printf("✅ 截图上传成功，URL: %s", imageURL)

	// 3. 发送图片 URL 给 WebSocket 服务器
	if err := c.sendScreenshotData(imageURL); err != nil {
		log.Printf("❌ 发送截图 URL 失败: %v", err)
		c.sendError("发送截图 URL 失败", err.Error())
	} else {
		log.Println("✅ 截图 URL 已发送到服务器")
	}
}

// sendScreenshotData 发送截图数据
func (c *WebSocketClient) sendScreenshotData(screenshotBase64 string) error {
	request := ScreenshotDataRequest{
		ScreenshotBase64: screenshotBase64,
		Timestamp:        time.Now().UnixMilli(),
	}

	// 发送截图数据
	err := c.Emit("screenshot_data", request)
	if err != nil {
		return fmt.Errorf("发送截图数据失败: %v", err)
	}

	log.Println("✓ 截图数据已发送")
	return nil
}

// sendError 发送错误信息
func (c *WebSocketClient) sendError(errorMsg, details string) {
	request := ErrorRequest{
		Error:   errorMsg,
		Details: details,
	}

	if err := c.Emit("error", request); err != nil {
		log.Printf("❌ 发送错误信息失败: %v", err)
	}
}

// pingLoop 定期发送 WebSocket ping 帧
func (c *WebSocketClient) pingLoop() {
	c.mu.RLock()
	ticker := c.pingTicker
	c.mu.RUnlock()

	if ticker == nil {
		log.Println("⚠️ pingTicker 未初始化，pingLoop 退出")
		return
	}

	defer func() {
		log.Println("📡 WebSocket ping 循环已停止")
	}()

	log.Println("📡 WebSocket ping 循环已启动")
	for {
		select {
		case <-c.done:
			// 收到退出信号
			return
		case <-ticker.C:
			// 定时发送 ping
			c.mu.RLock()
			conn := c.conn
			connected := c.connected
			c.mu.RUnlock()

			if !connected || conn == nil {
				log.Println("⚠️ 连接已断开，pingLoop 退出")
				return
			}

			// 发送 WebSocket 层 ping 帧
			c.mu.Lock()
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second))
			c.mu.Unlock()

			if err != nil {
				log.Printf("⚠️ 发送 WebSocket ping 失败: %v", err)
				// ping 失败说明连接有问题，关闭连接触发重连
				c.signalDone()
				if conn != nil {
					conn.Close()
				}
				return
			}
			log.Println("📡 WebSocket ping 已发送")
		}
	}
}

// Disconnect 断开连接
func (c *WebSocketClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connected = false

	// 停止心跳定时器
	c.signalDone()
	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}

	// 关闭连接
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected 检查是否已连接
func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// WebSocketConfig WebSocket 客户端配置
type WebSocketConfig struct {
	ServerURL  string // 服务器地址，如 http://localhost:3000
	DeviceCode string // 设备唯一标识码
}

// autoReconnect 自动重连机制
// 按照 1分钟、3分钟、10分钟的间隔进行重试
func (c *WebSocketClient) autoReconnect() {
	c.mu.Lock()
	currentTry := c.reconnectTry
	c.mu.Unlock()

	// 定义重连间隔：1分钟、3分钟、10分钟
	reconnectIntervals := []time.Duration{
		1 * time.Minute,  // 第一次重连：1分钟后
		3 * time.Minute,  // 第二次重连：3分钟后
		10 * time.Minute, // 第三次重连：10分钟后
	}

	// 如果已经尝试了3次，不再重连
	if currentTry >= len(reconnectIntervals) {
		log.Println("⚠️ 已达到最大重连次数，停止重连")
		return
	}

	// 获取当前重连间隔
	interval := reconnectIntervals[currentTry]
	log.Printf("🔄 将在 %v 后尝试第 %d 次重连...", interval, currentTry+1)

	// 等待指定时间
	time.Sleep(interval)

	// 尝试重连
	log.Printf("🔄 开始第 %d 次重连尝试...", currentTry+1)

	err := c.Connect()
	if err != nil {
		log.Printf("❌ 第 %d 次重连失败: %v", currentTry+1, err)

		// 增加重连计数
		c.mu.Lock()
		c.reconnectTry++
		c.mu.Unlock()

		// 继续下一次重连
		go c.autoReconnect()
	} else {
		log.Printf("✅ 第 %d 次重连成功！", currentTry+1)

		// 重连成功，重置计数
		c.mu.Lock()
		c.reconnectTry = 0
		c.mu.Unlock()

		// 重新启动心跳
		go c.StartHeartbeat()
	}
}

// StopReconnect 停止自动重连
func (c *WebSocketClient) StopReconnect() {
	c.mu.Lock()
	c.shouldReconnect = false
	c.mu.Unlock()
	log.Println("⚠️ 已停止自动重连")
}

// StartWebSocketClient 启动 WebSocket 客户端（封装函数）
func StartWebSocketClient(config WebSocketConfig) (*WebSocketClient, error) {
	// 验证配置参数
	if config.ServerURL == "" {
		return nil, fmt.Errorf("服务器地址不能为空")
	}
	if config.DeviceCode == "" {
		return nil, fmt.Errorf("设备编号不能为空")
	}

	log.Printf("========================================")
	log.Printf("正在启动 WebSocket 客户端")
	log.Printf("服务器地址: %s", config.ServerURL)
	log.Printf("设备编号: %s", config.DeviceCode)
	log.Printf("========================================")

	// 创建客户端
	client := NewWebSocketClient(config.ServerURL, config.DeviceCode)

	// 连接到服务器（带重试机制）
	err := client.connectWithRetry()
	if err != nil {
		log.Printf("⚠️ 初始连接失败，将在后台继续尝试重连")
		// 不返回错误，而是在后台继续尝试重连
		go client.autoReconnect()
	} else {
		// 启动心跳（在单独的 goroutine 中）
		go client.StartHeartbeat()
		log.Println("✅ WebSocket 客户端启动成功")
	}

	return client, nil
}

// connectWithRetry 带重试机制的连接（用于启动时）
func (c *WebSocketClient) connectWithRetry() error {
	// 定义重连间隔：立即、1分钟、3分钟、10分钟
	reconnectIntervals := []time.Duration{
		0,                // 第一次：立即尝试
		1 * time.Minute,  // 第二次：1分钟后
		3 * time.Minute,  // 第三次：3分钟后
		10 * time.Minute, // 第四次：10分钟后
	}

	var lastErr error
	for i := 0; i < len(reconnectIntervals); i++ {
		if i > 0 {
			interval := reconnectIntervals[i]
			log.Printf("🔄 将在 %v 后尝试第 %d 次连接...", interval, i+1)
			time.Sleep(interval)
			log.Printf("🔄 开始第 %d 次连接尝试...", i+1)
		}

		err := c.Connect()
		if err == nil {
			if i > 0 {
				log.Printf("✅ 第 %d 次连接成功！", i+1)
			}
			// 连接成功，重置重连计数
			c.mu.Lock()
			c.reconnectTry = 0
			c.mu.Unlock()
			return nil
		}

		lastErr = err
		log.Printf("❌ 第 %d 次连接失败: %v", i+1, err)
	}

	// 所有尝试都失败
	return fmt.Errorf("连接失败，已尝试 %d 次: %v", len(reconnectIntervals), lastErr)
}

// signalDone 安全地关闭 done channel，通知所有 goroutine 退出
// 使用 select 避免重复关闭已关闭的 channel
func (c *WebSocketClient) signalDone() {
	if c.done == nil {
		return
	}
	select {
	case <-c.done:
		// channel 已关闭，不做任何操作
	default:
		// channel 未关闭，关闭它
		close(c.done)
	}
}
