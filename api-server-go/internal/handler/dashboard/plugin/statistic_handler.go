package plugin

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

type StatisticHandler struct{}

func NewStatisticHandler() *StatisticHandler {
	return &StatisticHandler{}
}

func (h *StatisticHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{
		"todayAddContact":  0,
		"todayAddRoom":     0,
		"todayLossContact": 0,
		"totalContact":     0,
		"totalRoom":        0,
	})
}

func (h *StatisticHandler) Employees(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *StatisticHandler) TopList(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}