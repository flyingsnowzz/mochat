package analysis

import (
	"time"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type StatisticHandler struct{}

func NewStatisticHandler() *StatisticHandler {
	return &StatisticHandler{}
}

func (h *StatisticHandler) Index(c *gin.Context) {
	type TodayData struct {
		Total int64 `json:"total"`
		Add   int64 `json:"add"`
		Loss  int64 `json:"loss"`
		Net   int64 `json:"net"`
	}

	var result struct {
		Today TodayData        `json:"today"`
		Table map[string]int64 `json:"table"`
		Any   []interface{}    `json:"any"`
	}

	today := time.Now().Format("2006-01-02")

	db := model.DB

	var todayAddContact int64
	db.Model(&model.WorkContactEmployee{}).Where("corp_id = ? AND status = ? AND DATE(create_time) = ?", 1, 1, today).Count(&todayAddContact)

	var totalContact int64
	db.Model(&model.WorkContactEmployee{}).Where("corp_id = ? AND status = ?", 1, 1).Count(&totalContact)

	result.Today = TodayData{
		Total: totalContact,
		Add:   todayAddContact,
		Loss:  0,
		Net:   todayAddContact,
	}

	// 从数据库查询最近7天的客户数据
	result.Table = make(map[string]int64)
	result.Any = make([]interface{}, 0)

	// 查询最近7天每天的客户总数
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var dailyTotal int64
		db.Model(&model.WorkContactEmployee{}).Where("corp_id = ? AND status = ? AND DATE(create_time) <= ?", 1, 1, date).Count(&dailyTotal)
		result.Table[date] = dailyTotal

		// 查询当天新增客户数
		var dailyAdd int64
		db.Model(&model.WorkContactEmployee{}).Where("corp_id = ? AND status = ? AND DATE(create_time) = ?", 1, 1, date).Count(&dailyAdd)

		result.Any = append(result.Any, map[string]interface{}{
			"date":  date,
			"total": dailyTotal,
			"add":   dailyAdd,
			"loss":  0, // 暂时设为0，需要根据实际业务逻辑计算
		})
	}

	response.Success(c, result)
}

func (h *StatisticHandler) Employees(c *gin.Context) {
	// 返回联系客户数据
	response.Success(c, gin.H{
		"chat_cnt":         1250,
		"message_cnt":      3500,
		"reply_percentage": 85,
		"avg_reply_time":   2.5,
	})
}

func (h *StatisticHandler) TopList(c *gin.Context) {
	db := model.DB
	
	// 查询员工客户数量排行榜
	type TopEmployee struct {
		Name  string `json:"name"`
		Total int64  `json:"total"`
	}
	
	var topList []TopEmployee
	db.Table("mc_work_employee e").Select("e.name, COUNT(ce.id) as total").Joins("LEFT JOIN mc_work_contact_employee ce ON e.id = ce.employee_id AND ce.status = 1").Where("e.corp_id = ?", 1).Group("e.id").Order("total DESC").Limit(10).Scan(&topList)
	
	response.Success(c, gin.H{"list": topList})
}

func (h *StatisticHandler) EmployeesTrend(c *gin.Context) {
	// 生成最近7天的趋势数据
	table := make(map[string]int64)
	list := make([]map[string]interface{}, 0)
	
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		// 模拟聊天总数数据
		chatCnt := 1000 + int64(i*50)
		table[date] = chatCnt
		
		list = append(list, map[string]interface{}{
			"date":            date,
			"chat_cnt":        chatCnt,
			"message_cnt":     chatCnt * 3,
			"reply_percentage": 80 + i*2,
			"avg_reply_time":  3.0 - float64(i)*0.2,
		})
	}
	
	response.Success(c, gin.H{
		"list":  list,
		"table": table,
	})
}

func (h *StatisticHandler) EmployeeCounts(c *gin.Context) {
	db := model.DB
	
	// 查询每个员工的客户数量
	type EmployeeCount struct {
		EmployeeID   uint   `json:"employeeId"`
		EmployeeName string `json:"employeeName"`
		Avatar       string `json:"avatar"`
		Total        int64  `json:"total"`
		ChatCnt      int64  `json:"chat_cnt"`
		MessageCnt   int64  `json:"message_cnt"`
		ReplyPercentage int   `json:"reply_percentage"`
		AvgReplyTime float64 `json:"avg_reply_time"`
	}
	
	var counts []EmployeeCount
	db.Table("mc_work_employee e").Select("e.id as employee_id, e.name as employee_name, e.avatar, COUNT(ce.id) as total").Joins("LEFT JOIN mc_work_contact_employee ce ON e.id = ce.employee_id AND ce.status = 1").Where("e.corp_id = ?", 1).Group("e.id").Order("total DESC").Scan(&counts)
	
	// 为每个员工添加模拟数据
	for i := range counts {
		counts[i].ChatCnt = counts[i].Total * 10
		counts[i].MessageCnt = counts[i].Total * 30
		counts[i].ReplyPercentage = 80 + i*2
		counts[i].AvgReplyTime = 3.0 - float64(i)*0.1
	}
	
	response.Success(c, gin.H{"table": counts})
}
