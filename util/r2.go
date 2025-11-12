package util

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Client Cloudflare R2 存储客户端
type R2Client struct {
	client       *s3.Client
	bucketName   string
	publicDomain string // 公开访问域名
}

// R2Config R2 客户端配置
type R2Config struct {
	AccountID       string // Cloudflare 账户 ID
	AccessKeyID     string // R2 访问密钥 ID
	AccessKeySecret string // R2 访问密钥
	BucketName      string // 存储桶名称
	PublicDomain    string // 公开访问域名（可选，如 "pub-xxxxx.r2.dev" 或自定义域名）
}

// NewR2Client 创建新的 R2 客户端
func NewR2Client(cfg R2Config) (*R2Client, error) {
	// 验证配置参数
	if cfg.AccountID == "" {
		return nil, fmt.Errorf("账户 ID 不能为空")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("访问密钥 ID 不能为空")
	}
	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("访问密钥不能为空")
	}
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("存储桶名称不能为空")
	}

	// 加载 AWS SDK 配置
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	// 创建 S3 客户端，指向 R2 端点
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})

	// 如果没有提供公开域名，使用默认的 R2.dev 域名格式
	publicDomain := cfg.PublicDomain
	if publicDomain == "" {
		// 默认使用 R2 的公开访问格式（需要在 Cloudflare 控制台启用公开访问）
		publicDomain = fmt.Sprintf("https://pub-%s.r2.dev", cfg.AccountID)
	}

	return &R2Client{
		client:       client,
		bucketName:   cfg.BucketName,
		publicDomain: publicDomain,
	}, nil
}

// UploadFile 上传文件到 R2
// filePath: 文件在 R2 中的路径（如 "images/screenshot.png"）
// data: 文件数据（字节数组）
// contentType: 文件 MIME 类型（如 "image/png"）
// 返回: 文件的公开访问 URL 和错误信息
func (r *R2Client) UploadFile(filePath string, data []byte, contentType string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("文件数据不能为空")
	}

	// 上传文件
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(filePath),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	// 构建文件的公开访问 URL
	fileURL := fmt.Sprintf("%s/%s", r.publicDomain, filePath)

	return fileURL, nil
}

// UploadFileWithReader 使用 io.Reader 上传文件
func (r *R2Client) UploadFileWithReader(filePath string, reader io.Reader, contentType string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(filePath),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	// 构建文件的公开访问 URL
	fileURL := fmt.Sprintf("%s/%s", r.publicDomain, filePath)
	return fileURL, nil
}

// UploadFileWithTimestamp 上传文件并自动添加时间戳到文件名
// prefix: 文件路径前缀（如 "screenshots/"）
// filename: 原始文件名（如 "screen.png"）
// data: 文件数据
// contentType: 文件 MIME 类型
func (r *R2Client) UploadFileWithTimestamp(prefix, filename string, data []byte, contentType string) (string, error) {
	// 生成带时间戳的文件名
	timestamp := time.Now().Format("20060102_150405")
	filePath := fmt.Sprintf("%s%s_%s", prefix, timestamp, filename)

	return r.UploadFile(filePath, data, contentType)
}

// DeleteFile 删除 R2 中的文件
func (r *R2Client) DeleteFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// ListFiles 列出存储桶中的所有文件
func (r *R2Client) ListFiles(prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucketName),
	}

	// 如果指定了前缀，只列出该前缀下的文件
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	output, err := r.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %v", err)
	}

	var files []string
	for _, object := range output.Contents {
		if object.Key != nil {
			files = append(files, *object.Key)
		}
	}

	return files, nil
}

// FileExists 检查文件是否存在
func (r *R2Client) FileExists(filePath string) (bool, error) {
	if filePath == "" {
		return false, fmt.Errorf("文件路径不能为空")
	}

	_, err := r.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		// 如果是 NotFound 错误，返回 false
		return false, nil
	}

	return true, nil
}

// UploadImage 上传 image.Image 类型的图片（支持 *image.NRGBA）
// filePath: 文件在 R2 中的路径（如 "screenshots/screen.png"）
// img: 图片对象（*image.NRGBA 或其他 image.Image 类型）
// format: 图片格式 "png" 或 "jpeg"
// quality: JPEG 质量（1-100），PNG 格式时忽略此参数
func (r *R2Client) UploadImage(filePath string, img image.Image, format string, quality int) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}
	if img == nil {
		return "", fmt.Errorf("图片数据不能为空")
	}

	// 将图片编码为字节数组
	var buf bytes.Buffer
	var contentType string

	switch format {
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return "", fmt.Errorf("PNG 编码失败: %v", err)
		}
		contentType = "image/png"

	case "jpeg", "jpg":
		// 确保质量在有效范围内
		if quality < 1 {
			quality = 1
		}
		if quality > 100 {
			quality = 100
		}

		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return "", fmt.Errorf("JPEG 编码失败: %v", err)
		}
		contentType = "image/jpeg"

	default:
		return "", fmt.Errorf("不支持的图片格式: %s（仅支持 png 或 jpeg）", format)
	}

	// 上传到 R2
	return r.UploadFile(filePath, buf.Bytes(), contentType)
}

// UploadImageWithTimestamp 上传图片并自动添加时间戳
// prefix: 文件路径前缀（如 "screenshots/"）
// filename: 原始文件名（如 "screen.png"）
// img: 图片对象（*image.NRGBA 或其他 image.Image 类型）
// format: 图片格式 "png" 或 "jpeg"
// quality: JPEG 质量（1-100），PNG 格式时忽略此参数
func (r *R2Client) UploadImageWithTimestamp(prefix, filename string, img image.Image, format string, quality int) (string, error) {
	// 生成带时间戳的文件名
	timestamp := time.Now().Format("20060102_150405")
	filePath := fmt.Sprintf("%s%s_%s", prefix, timestamp, filename)

	return r.UploadImage(filePath, img, format, quality)
}
