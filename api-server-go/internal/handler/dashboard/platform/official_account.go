package platform

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

type OfficialAccountHandler struct{}

func NewOfficialAccountHandler() *OfficialAccountHandler {
	return &OfficialAccountHandler{}
}

func (h *OfficialAccountHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *OfficialAccountHandler) GetPreAuthUrl(c *gin.Context) {
	response.Success(c, gin.H{"url": "https://mp.weixin.qq.com/cgi-bin/componentloginpage"})
}

func (h *OfficialAccountHandler) AuthEventCallback(c *gin.Context) {
	c.String(200, "success")
}

func (h *OfficialAccountHandler) MessageEventCallback(c *gin.Context) {
	c.String(200, "success")
}

func (h *OfficialAccountHandler) Set(c *gin.Context) {
	response.SuccessMsg(c, "设置成功")
}

func (h *OfficialAccountHandler) AuthRedirect(c *gin.Context) {
	c.Redirect(302, "/")
}
