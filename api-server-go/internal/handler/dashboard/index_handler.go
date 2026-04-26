// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含仪表盘首页的处理器，提供数据统计和图表数据功能
package dashboard

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
)

// DashboardIndexHandler 仪表盘首页处理器
// 提供仪表盘首页的数据统计和图表数据功能
// 主要职责：
// 1. 获取首页统计数据
// 2. 获取客户增长趋势图表数据
//
// 依赖服务：
// - gorm.DB: 数据库连接

type DashboardIndexHandler struct {
	db *gorm.DB // 数据库连接
}

// NewDashboardIndexHandler 创建仪表盘首页处理器实例
// 参数：db - GORM 数据库连接
// 返回：仪表盘首页处理器实例
func NewDashboardIndexHandler(db *gorm.DB) *DashboardIndexHandler {
	return &DashboardIndexHandler{db: db}
}

// Index 获取仪表盘首页统计数据
// 获取仪表盘首页的各种统计数据，包括今日新增、流失等指标
// 处理流程：
// 1. 获取租户 ID
// 2. 计算今日开始和结束时间
// 3. 统计今日新增客户
// 4. 统计今日新增群聊
// 5. 统计今日加入群聊人数
// 6. 统计今日流失客户
// 7. 统计今日退出群聊人数
// 8. 统计总客户数
// 9. 统计总群聊数
// 10. 返回统计结果
//
// 返回：包含各种统计数据的响应
func (h *DashboardIndexHandler) Index(c *gin.Context) {
	// 获取租户 ID
	tenantID, exists := c.Get("tenantId")
	if !exists {
		response.Fail(c, response.ErrAuth, "未获取到租户信息")
		return
	}

	// 计算今日开始和结束时间
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	todayEnd := todayStart.Add(24 * time.Hour)

	// 统计今日新增客户
	var todayAddContact int64
	h.db.Model(&model.WorkContact{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, todayStart, todayEnd).Count(&todayAddContact)

	// 统计今日新增群聊
	var todayAddRoom int64
	h.db.Model(&model.WorkRoom{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, todayStart, todayEnd).Count(&todayAddRoom)

	// 统计今日加入群聊人数
	var todayAddIntoRoom int64
	h.db.Model(&model.WorkContactRoom{}).Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Count(&todayAddIntoRoom)

	// 统计今日流失客户
	var todayLossContact int64
	h.db.Model(&model.WorkContact{}).Where("tenant_id = ? AND deleted_at >= ? AND deleted_at < ?", tenantID, todayStart, todayEnd).Count(&todayLossContact)

	// 统计今日退出群聊人数
	var todayQuitRoom int64
	h.db.Model(&model.WorkContactRoom{}).Where("out_time >= ? AND out_time < ?", todayStart, todayEnd).Count(&todayQuitRoom)

	// 统计总客户数
	var totalContact int64
	h.db.Model(&model.WorkContact{}).Where("tenant_id = ?", tenantID).Count(&totalContact)

	// 统计总群聊数
	var totalRoom int64
	h.db.Model(&model.WorkRoom{}).Where("tenant_id = ?", tenantID).Count(&totalRoom)

	// 返回统计结果
	response.Success(c, gin.H{
		"todayAddContact":  todayAddContact,  // 今日新增客户数
		"todayAddRoom":     todayAddRoom,     // 今日新增群聊数
		"todayAddIntoRoom": todayAddIntoRoom, // 今日加入群聊人数
		"todayLossContact": todayLossContact, // 今日流失客户数
		"todayQuitRoom":    todayQuitRoom,    // 今日退出群聊人数
		"totalContact":     totalContact,     // 总客户数
		"totalRoom":        totalRoom,        // 总群聊数
	})
}

// LineChat 获取客户增长趋势图表数据
// 获取最近7天的客户新增和流失数据，用于生成趋势图表
// 处理流程：
// 1. 获取租户 ID
// 2. 初始化日期、新增客户数据和流失客户数据数组
// 3. 循环计算最近7天的日期和对应的数据
// 4. 统计每天的新增客户数
// 5. 统计每天的流失客户数
// 6. 返回图表数据
//
// 返回：包含最近7天日期、新增客户数据和流失客户数据的响应
func (h *DashboardIndexHandler) LineChat(c *gin.Context) {
	// 获取租户 ID
	tenantID, exists := c.Get("tenantId")
	if !exists {
		response.Fail(c, response.ErrAuth, "未获取到租户信息")
		return
	}

	// 初始化日期、新增客户数据和流失客户数据数组
	dates := make([]string, 7)
	data := make([]int, 7)
	lossData := make([]int, 7)

	// 循环计算最近7天的日期和对应的数据
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dates[6-i] = date.Format("01-02")

		// 计算当天的开始和结束时间
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		// 统计当天新增客户
		var count int64
		h.db.Model(&model.WorkContact{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, dayStart, dayEnd).Count(&count)
		data[6-i] = int(count)

		// 统计当天流失客户
		var lossCount int64
		h.db.Model(&model.WorkContact{}).Where("tenant_id = ? AND deleted_at >= ? AND deleted_at < ?", tenantID, dayStart, dayEnd).Count(&lossCount)
		lossData[6-i] = int(lossCount)
	}

	// 返回图表数据
	response.Success(c, gin.H{
		"xAxis": dates, // 最近7天的日期
		"series": []map[string]interface{}{
			{"name": "新增客户", "data": data},     // 新增客户数据
			{"name": "流失客户", "data": lossData}, // 流失客户数据
		},
	})
}
