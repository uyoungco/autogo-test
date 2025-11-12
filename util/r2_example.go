package util

import (
	"fmt"
	"log"

	"github.com/Dasongzi1366/AutoGo/images"
)

// 使用示例（不会被编译执行，仅供参考）
func ExampleR2Usage() {
	// 1. 创建 R2 客户端
	client, err := NewR2Client(R2Config{
		AccountID:       "your-account-id",
		AccessKeyID:     "your-access-key-id",
		AccessKeySecret: "your-access-key-secret",
		BucketName:      "your-bucket-name",
	})
	if err != nil {
		log.Fatalf("创建 R2 客户端失败: %v", err)
	}

	// 2. 上传文件（基本用法）
	fileData := []byte("Hello, R2!")
	fileURL, err := client.UploadFile("test/hello.txt", fileData, "text/plain")
	if err != nil {
		log.Printf("上传文件失败: %v", err)
	} else {
		fmt.Printf("文件上传成功，URL: %s\n", fileURL)
	}

	// 3. 上传图片文件（字节数组）
	// imageData := []byte{...} // 你的图片数据
	// imageURL, err := client.UploadFile("images/photo.jpg", imageData, "image/jpeg")

	// 4. 上传文件并自动添加时间戳
	// 会生成类似 "screenshots/20231112_143025_screen.png" 的文件名
	// screenshotData := []byte{...}
	// url, err := client.UploadFileWithTimestamp("screenshots/", "screen.png", screenshotData, "image/png")

	// ========== 上传 *image.NRGBA 类型的截图 ==========

	// 5. 上传截图（*image.NRGBA 类型）- PNG 格式
	screenshot := images.CaptureScreen(0, 0, 0, 0, 0) // 返回 *image.NRGBA
	screenshotURL, err := client.UploadImage("screenshots/screen.png", screenshot, "png", 0)
	if err != nil {
		log.Printf("上传截图失败: %v", err)
	} else {
		fmt.Printf("截图上传成功，URL: %s\n", screenshotURL)
	}

	// 6. 上传截图（*image.NRGBA 类型）- JPEG 格式，质量 80
	screenshot2 := images.CaptureScreen(0, 0, 0, 0, 0)
	screenshotURL2, err := client.UploadImage("screenshots/screen.jpg", screenshot2, "jpeg", 80)
	if err != nil {
		log.Printf("上传截图失败: %v", err)
	} else {
		fmt.Printf("截图上传成功，URL: %s\n", screenshotURL2)
	}

	// 7. 上传截图并自动添加时间戳（推荐使用）
	// 会生成类似 "screenshots/20231112_143025_screen.jpg" 的文件名
	screenshot3 := images.CaptureScreen(0, 0, 0, 0, 0)
	screenshotURL3, err := client.UploadImageWithTimestamp("screenshots/", "screen.jpg", screenshot3, "jpeg", 70)
	if err != nil {
		log.Printf("上传截图失败: %v", err)
	} else {
		fmt.Printf("截图上传成功，URL: %s\n", screenshotURL3)
	}

	// ========== 其他操作 ==========

	// 8. 列出所有文件
	files, err := client.ListFiles("")
	if err != nil {
		log.Printf("列出文件失败: %v", err)
	} else {
		fmt.Printf("存储桶中的文件: %v\n", files)
	}

	// 9. 列出指定前缀的文件
	screenshots, err := client.ListFiles("screenshots/")
	if err != nil {
		log.Printf("列出截图失败: %v", err)
	} else {
		fmt.Printf("截图文件: %v\n", screenshots)
	}

	// 10. 检查文件是否存在
	exists, err := client.FileExists("test/hello.txt")
	if err != nil {
		log.Printf("检查文件失败: %v", err)
	} else {
		fmt.Printf("文件是否存在: %v\n", exists)
	}

	// 11. 删除文件
	err = client.DeleteFile("test/hello.txt")
	if err != nil {
		log.Printf("删除文件失败: %v", err)
	} else {
		fmt.Println("文件删除成功")
	}
}

// ExampleUploadScreenshot 上传截图的完整示例
func ExampleUploadScreenshot() {
	// 创建 R2 客户端
	client, err := NewR2Client(R2Config{
		AccountID:       "your-account-id",
		AccessKeyID:     "your-access-key-id",
		AccessKeySecret: "your-access-key-secret",
		BucketName:      "your-bucket-name",
	})
	if err != nil {
		log.Fatalf("创建 R2 客户端失败: %v", err)
	}

	// 截取屏幕
	screenshot := images.CaptureScreen(0, 0, 0, 0, 0) // 返回 *image.NRGBA

	// 方式 1: 上传为 PNG（无损，文件较大）
	url1, err := client.UploadImageWithTimestamp("screenshots/", "screen.png", screenshot, "png", 0)
	if err != nil {
		log.Printf("上传失败: %v", err)
	} else {
		fmt.Printf("PNG 上传成功: %s\n", url1)
	}

	// 方式 2: 上传为 JPEG（有损压缩，文件较小，推荐）
	// 质量参数: 1-100，建议 60-80
	url2, err := client.UploadImageWithTimestamp("screenshots/", "screen.jpg", screenshot, "jpeg", 70)
	if err != nil {
		log.Printf("上传失败: %v", err)
	} else {
		fmt.Printf("JPEG 上传成功: %s\n", url2)
	}
}
