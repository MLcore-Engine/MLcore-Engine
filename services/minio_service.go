package services

import (
	"MLcore-Engine/common"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

var minioClient *minio.Client

// 初始化MinIO客户端
func initMinioClient() error {
	if minioClient != nil {
		return nil
	}

	// 从配置文件读取MinIO配置
	endpoint := viper.GetString("minio.endpoint")
	accessKey := viper.GetString("minio.accessKey")
	secretKey := viper.GetString("minio.secretKey")
	useSSL := viper.GetBool("minio.useSSL")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return errors.New("MinIO配置不完整")
	}

	// 创建MinIO客户端
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return fmt.Errorf("初始化MinIO客户端失败: %v", err)
	}

	return nil
}

// InitDatasetMinioStorage 初始化数据集的MinIO存储
// 创建存储桶(如果不存在)并初始化一个空的JSONL文件
func InitDatasetMinioStorage(bucketName, objectPath string) error {
	// 初始化MinIO客户端
	if err := initMinioClient(); err != nil {
		return err
	}

	ctx := context.Background()

	// 检查存储桶是否存在
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("检查存储桶失败: %v", err)
	}

	// 如果存储桶不存在，创建它
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建存储桶失败: %v", err)
		}
		common.SysLog(fmt.Sprintf("成功创建存储桶: %s", bucketName))

		// 设置存储桶策略为只读(可选)
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": "*"},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`

		err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			common.SysError(fmt.Errorf("设置存储桶策略失败: %v", err).Error())
			// 继续执行，不中断流程
		}
	}

	// 创建一个空的JSONL文件
	emptyContent := []byte("")
	reader := bytes.NewReader(emptyContent)

	// 上传到MinIO
	_, err = minioClient.PutObject(ctx, bucketName, objectPath, reader, int64(len(emptyContent)), minio.PutObjectOptions{
		ContentType: "application/jsonl",
	})

	if err != nil {
		return fmt.Errorf("创建空JSONL文件失败: %v", err)
	}

	common.SysLog(fmt.Sprintf("成功初始化数据集存储: %s/%s", bucketName, objectPath))
	return nil
}

// DeleteDatasetMinioObject 删除MinIO中的数据集对象
func DeleteDatasetMinioObject(bucketName, objectPath string) error {
	// 初始化MinIO客户端
	if err := initMinioClient(); err != nil {
		return err
	}

	ctx := context.Background()

	// 删除对象
	err := minioClient.RemoveObject(ctx, bucketName, objectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除对象失败: %v", err)
	}

	common.SysLog(fmt.Sprintf("成功删除数据集对象: %s/%s", bucketName, objectPath))
	return nil
}

// 添加读取JSONL条目的函数
func ReadJSONLFromMinio(bucketName, objectPath string, offset, limit int) ([]string, error) {
	// 初始化MinIO客户端
	if err := initMinioClient(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 获取对象
	object, err := minioClient.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取对象失败: %v", err)
	}
	defer object.Close()

	// 读取对象内容
	var lines []string
	lineCount := 0

	// 使用缓冲读取逐行处理
	buffer := make([]byte, 4096)
	lineBuffer := bytes.NewBuffer(nil)

	for {
		n, err := object.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("读取对象内容失败: %v", err)
		}

		if n == 0 {
			break
		}

		// 处理当前读取的块
		for i := 0; i < n; i++ {
			b := buffer[i]
			if b == '\n' {
				// 行结束
				if lineCount >= offset && (limit <= 0 || lineCount < offset+limit) {
					lines = append(lines, lineBuffer.String())
				}
				lineBuffer.Reset()
				lineCount++

				// 如果已经达到限制，退出
				if limit > 0 && lineCount >= offset+limit {
					return lines, nil
				}
			} else {
				// 构建当前行
				if lineCount >= offset && (limit <= 0 || lineCount < offset+limit) {
					lineBuffer.WriteByte(b)
				}
			}
		}
	}

	// 处理最后一行(如果没有换行结束)
	if lineBuffer.Len() > 0 {
		if lineCount >= offset && (limit <= 0 || lineCount < offset+limit) {
			lines = append(lines, lineBuffer.String())
		}
	}

	return lines, nil
}

// UpdateJSONLInMinio 更新MinIO中的JSONL文件的特定行
func UpdateJSONLInMinio(bucketName, objectPath string, lineIndex int, newContent string) error {
	// 读取所有行
	lines, err := ReadJSONLFromMinio(bucketName, objectPath, 0, -1)
	if err != nil {
		return err
	}

	// 更新指定行或添加新行
	if lineIndex < 0 {
		// 添加到末尾
		lines = append(lines, newContent)
	} else if lineIndex < len(lines) {
		// 更新现有行
		lines[lineIndex] = newContent
	} else {
		// 填充空行直到目标索引
		for i := len(lines); i < lineIndex; i++ {
			lines = append(lines, "{}")
		}
		// 添加新内容
		lines = append(lines, newContent)
	}

	// 构建新的内容
	content := bytes.NewBuffer(nil)
	for i, line := range lines {
		content.WriteString(line)
		if i < len(lines)-1 {
			content.WriteByte('\n')
		}
	}

	// 上传回MinIO
	reader := bytes.NewReader(content.Bytes())
	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		objectPath,
		reader,
		int64(content.Len()),
		minio.PutObjectOptions{ContentType: "application/jsonl"},
	)

	return err
}

// AppendJSONLToMinio 在MinIO中的JSONL文件末尾追加内容
func AppendJSONLToMinio(bucketName, objectPath string, content string) error {
	return UpdateJSONLInMinio(bucketName, objectPath, -1, content)
}

// 获取MinIO存储对象的预签名URL
func GetPresignedURL(bucketName, objectPath string, expires time.Duration) (string, error) {
	// 初始化MinIO客户端
	if err := initMinioClient(); err != nil {
		return "", err
	}

	ctx := context.Background()

	// 生成预签名URL
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectPath, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}

	return presignedURL.String(), nil
}

// GetMinioObject 从MinIO获取对象并返回一个可读取的对象流
func GetMinioObject(bucketName, objectPath string) (io.ReadCloser, error) {
	// 初始化MinIO客户端
	if err := initMinioClient(); err != nil {
		return nil, err
	}

	ctx := context.Background()

	// 获取对象
	object, err := minioClient.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取MinIO对象失败: %v", err)
	}

	return object, nil
}

// UploadJSONLToMinio 上传JSONL文件到MinIO
func UploadJSONLToMinio(bucketName, objectPath string, reader io.Reader) error {
	// 初始化MinIO客户端
	if err := initMinioClient(); err != nil {
		return err
	}

	ctx := context.Background()

	// 上传对象
	_, err := minioClient.PutObject(ctx, bucketName, objectPath, reader, -1, minio.PutObjectOptions{
		ContentType: "application/jsonl",
	})

	if err != nil {
		return fmt.Errorf("上传到MinIO失败: %v", err)
	}

	return nil
}
