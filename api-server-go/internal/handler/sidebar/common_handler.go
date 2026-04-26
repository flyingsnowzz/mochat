package sidebar

import (
	"io"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/storage"
)

// CommonHandler 侧边栏通用处理器
// 处理侧边栏相关的通用操作，包括文件上传和微信JSSDK配置等

type CommonHandler struct {
	storage storage.Storage // 存储服务实例
}

// NewCommonHandler 创建通用处理器实例
// 参数:
//   - storage: 存储服务实例
// 返回值:
//   - *CommonHandler: 通用处理器实例

func NewCommonHandler(storage storage.Storage) *CommonHandler {
	return &CommonHandler{storage: storage}
}

// Upload 文件上传
// 请求方法: POST
// 请求路径: /sidebar/common/upload
// 请求参数:
//   - file: 要上传的文件
// 响应:
//   - 成功: 包含上传文件URL的对象
//   - 失败: 错误信息

func (h *CommonHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, response.ErrFileUpload, "上传文件失败")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := time.Now().Format("20060102150405") + ext
	path := "uploads/" + filename

	url, err := h.storage.Upload(file, path)
	if err != nil {
		response.Fail(c, response.ErrFileUpload, "上传文件失败")
		return
	}

	response.Success(c, gin.H{"url": url})
}

// WxJSSDKConfig 微信JSSDK配置
// 请求方法: GET
// 请求路径: /sidebar/wxJSSDK/config
// 响应:
//   - 成功: 包含微信JSSDK配置的对象

func (h *CommonHandler) WxJSSDKConfig(c *gin.Context) {
	response.Success(c, gin.H{
		"appId":     "wx1234567890",
		"timestamp": 1234567890,
		"nonceStr":  "randomstring",
		"signature": "abcdefg",
	})
}

// ReadAll 读取文件内容
// 参数:
//   - path: 文件路径
// 返回值:
//   - io.ReadCloser: 文件读取器
//   - error: 错误信息

func (h *CommonHandler) ReadAll(path string) (io.ReadCloser, error) {
	return nil, nil
}
