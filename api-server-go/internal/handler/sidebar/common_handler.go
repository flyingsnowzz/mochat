package sidebar

import (
	"io"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/storage"
)

type CommonHandler struct {
	storage storage.Storage
}

func NewCommonHandler(storage storage.Storage) *CommonHandler {
	return &CommonHandler{storage: storage}
}

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

func (h *CommonHandler) WxJSSDKConfig(c *gin.Context) {
	response.Success(c, gin.H{
		"appId":     "wx1234567890",
		"timestamp": 1234567890,
		"nonceStr":  "randomstring",
		"signature": "abcdefg",
	})
}

func (h *CommonHandler) ReadAll(path string) (io.ReadCloser, error) {
	return nil, nil
}
