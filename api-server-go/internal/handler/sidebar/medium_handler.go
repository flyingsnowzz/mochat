package sidebar

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

type MediumHandler struct{}

func NewMediumHandler() *MediumHandler {
	return &MediumHandler{}
}

func (h *MediumHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *MediumHandler) MediaIdUpdate(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}

type MediumGroupHandler struct{}

func NewMediumGroupHandler() *MediumGroupHandler {
	return &MediumGroupHandler{}
}

func (h *MediumGroupHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
