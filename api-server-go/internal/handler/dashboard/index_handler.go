package dashboard

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
)

type DashboardIndexHandler struct {
	db *gorm.DB
}

func NewDashboardIndexHandler(db *gorm.DB) *DashboardIndexHandler {
	return &DashboardIndexHandler{db: db}
}

func (h *DashboardIndexHandler) Index(c *gin.Context) {
	tenantID, exists := c.Get("tenantId")
	if !exists {
		response.Fail(c, response.ErrAuth, "未获取到租户信息")
		return
	}

	// 获取今日开始时间
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

	response.Success(c, gin.H{
		"todayAddContact":  todayAddContact,
		"todayAddRoom":     todayAddRoom,
		"todayAddIntoRoom": todayAddIntoRoom,
		"todayLossContact": todayLossContact,
		"todayQuitRoom":    todayQuitRoom,
		"totalContact":     totalContact,
		"totalRoom":        totalRoom,
	})
}

func (h *DashboardIndexHandler) LineChat(c *gin.Context) {
	tenantID, exists := c.Get("tenantId")
	if !exists {
		response.Fail(c, response.ErrAuth, "未获取到租户信息")
		return
	}

	// 获取最近7天的日期
	dates := make([]string, 7)
	data := make([]int, 7)
	lossData := make([]int, 7)

	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dates[6-i] = date.Format("01-02")

		// 统计当天新增客户
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		var count int64
		h.db.Model(&model.WorkContact{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, dayStart, dayEnd).Count(&count)
		data[6-i] = int(count)

		// 统计当天流失客户
		var lossCount int64
		h.db.Model(&model.WorkContact{}).Where("tenant_id = ? AND deleted_at >= ? AND deleted_at < ?", tenantID, dayStart, dayEnd).Count(&lossCount)
		lossData[6-i] = int(lossCount)
	}

	response.Success(c, gin.H{
		"xAxis": dates,
		"series": []map[string]interface{}{
			{"name": "新增客户", "data": data},
			{"name": "流失客户", "data": lossData},
		},
	})
}
