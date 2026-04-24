package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mochat-api-server/internal/config"
)

type Storage interface {
	Upload(src io.Reader, dstPath string) (string, error)
	GetURL(path string) string
	Delete(path string) error
}

var DefaultStorage Storage

func InitStorage(cfg config.FileConfig, apiBaseURL string) {
	switch cfg.Driver {
	case "local":
		DefaultStorage = &LocalStorage{BaseURL: apiBaseURL}
	case "oss":
		DefaultStorage = &OSSStorage{Config: cfg.OSS}
	case "cos":
		DefaultStorage = &COSStorage{Config: cfg.COS}
	case "s3", "minio":
		DefaultStorage = &S3Storage{Config: cfg.S3}
	case "qiniu":
		DefaultStorage = &QiniuStorage{Config: cfg.Qiniu}
	default:
		DefaultStorage = &LocalStorage{BaseURL: apiBaseURL}
	}
}

type LocalStorage struct {
	BaseURL string
}

func (s *LocalStorage) Upload(src io.Reader, dstPath string) (string, error) {
	fullPath := filepath.Join("storage", "upload", dstPath)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create dir failed: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file failed: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, src); err != nil {
		return "", fmt.Errorf("copy file failed: %w", err)
	}

	return dstPath, nil
}

func (s *LocalStorage) GetURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.BaseURL, "/"), path)
}

func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join("storage", "upload", path)
	return os.Remove(fullPath)
}

type OSSStorage struct {
	Config config.OSSConfig
}

func (s *OSSStorage) Upload(src io.Reader, dstPath string) (string, error) {
	return dstPath, fmt.Errorf("OSS upload not implemented yet, use aliyun-oss-sdk")
}

func (s *OSSStorage) GetURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return fmt.Sprintf("https://%s.%s/%s", s.Config.Bucket, s.Config.Endpoint, path)
}

func (s *OSSStorage) Delete(path string) error {
	return fmt.Errorf("OSS delete not implemented yet")
}

type COSStorage struct {
	Config config.COSConfig
}

func (s *COSStorage) Upload(src io.Reader, dstPath string) (string, error) {
	return dstPath, fmt.Errorf("COS upload not implemented yet")
}

func (s *COSStorage) GetURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com/%s", s.Config.Bucket, s.Config.AppID, s.Config.Region, path)
}

func (s *COSStorage) Delete(path string) error {
	return fmt.Errorf("COS delete not implemented yet")
}

type S3Storage struct {
	Config config.S3Config
}

func (s *S3Storage) Upload(src io.Reader, dstPath string) (string, error) {
	return dstPath, fmt.Errorf("S3 upload not implemented yet")
}

func (s *S3Storage) GetURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.Config.Bucket, s.Config.Region, path)
}

func (s *S3Storage) Delete(path string) error {
	return fmt.Errorf("S3 delete not implemented yet")
}

type QiniuStorage struct {
	Config config.QiniuConfig
}

func (s *QiniuStorage) Upload(src io.Reader, dstPath string) (string, error) {
	return dstPath, fmt.Errorf("Qiniu upload not implemented yet")
}

func (s *QiniuStorage) GetURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.Config.Domain, "/"), path)
}

func (s *QiniuStorage) Delete(path string) error {
	return fmt.Errorf("Qiniu delete not implemented yet")
}

func GenerateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}
