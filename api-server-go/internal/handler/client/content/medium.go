package content

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

// MediumHandler 侧边栏媒体处理器
// 处理侧边栏相关的媒体操作，包括媒体列表和媒体ID更新等

type MediumHandler struct{}

// NewMediumHandler 创建媒体处理器实例
// 返回值:
//   - *MediumHandler: 媒体处理器实例

func NewMediumHandler() *MediumHandler {
	return &MediumHandler{}
}

// Index 获取媒体列表
// 请求方法: GET
// 请求路径: /sidebar/medium/index
// 响应:
//   - 成功: 包含媒体列表的对象

func (h *MediumHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// MediaIdUpdate 更新媒体ID
// 请求方法: PUT
// 请求路径: /sidebar/medium/mediaIdUpdate/:id
// 请求参数:
//   - id: 媒体ID
// 响应:
//   - 成功: 更新成功消息

func (h *MediumHandler) MediaIdUpdate(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}

// MediumGroupHandler 侧边栏媒体分组处理器
// 处理侧边栏相关的媒体分组操作，包括媒体分组列表等

type MediumGroupHandler struct{}

// NewMediumGroupHandler 创建媒体分组处理器实例
// 返回值:
//   - *MediumGroupHandler: 媒体分组处理器实例

func NewMediumGroupHandler() *MediumGroupHandler {
	return &MediumGroupHandler{}
}

// Index 获取媒体分组列表
// 请求方法: GET
// 请求路径: /sidebar/mediumGroup/index
// 响应:
//   - 成功: 包含媒体分组列表的对象

func (h *MediumGroupHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
