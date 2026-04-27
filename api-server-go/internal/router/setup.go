package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/config"
)

// SetupRouter 初始化并配置应用的路由（兼容旧接口）
// 这个函数保持原有的接口签名，内部使用新的 Router 结构体
func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := NewRouter(cfg, db)
	return router.Setup()
}
