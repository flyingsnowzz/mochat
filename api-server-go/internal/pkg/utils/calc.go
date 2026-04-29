package utils

import "time"

// CalcFissionStatus 根据 endTime 计算活动状态
// 已过期返回 2（已结束），未过期或无截止时间返回 1（进行中）
func CalcFissionStatus(endTime *time.Time) int {
	if endTime != nil && time.Now().After(*endTime) {
		return 2
	}
	return 1
}
