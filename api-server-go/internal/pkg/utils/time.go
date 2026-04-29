package utils

import "time"

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

// FormatTime 格式化 time.Time，零值返回空字符串
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateTimeFormat)
}

// FormatDate 格式化 time.Time 为日期字符串
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateFormat)
}

// FormatTimePtr 格式化 *time.Time，nil 返回空字符串
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(DateTimeFormat)
}

// Now 获取当前时间字符串
func Now() string {
	return time.Now().Format(DateTimeFormat)
}

// TodayDate 获取今天的日期字符串
func TodayDate() string {
	return time.Now().Format(DateFormat)
}

// BeginOfDay 获取指定日期的 00:00:00
func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取指定日期的 23:59:59
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}
