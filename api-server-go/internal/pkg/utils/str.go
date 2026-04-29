package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Explode 类似 PHP 的 explode，将字符串按分隔符拆分为切片
func Explode(sep, s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

// Implode 类似 PHP 的 implode，将字符串切片按分隔符合并
func Implode(sep string, items []string) string {
	return strings.Join(items, sep)
}

// TrimSpace 去除字符串首尾空白
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// Trim 去除字符串首尾指定字符
func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// StrPos 类似 PHP 的 strpos，返回子串在字符串中的位置，-1 表示未找到
func StrPos(haystack, needle string) int {
	return strings.Index(haystack, needle)
}

// Contains 检查字符串是否包含子串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Replace 替换字符串
func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

// MD5 计算字符串的 MD5 哈希
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// RandomString 生成指定长度的随机字符串（字母+数字）
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Substr 类似 PHP 的 substr
func Substr(s string, start int, length ...int) string {
	if start < 0 {
		start = len(s) + start
		if start < 0 {
			start = 0
		}
	}
	if start >= len(s) {
		return ""
	}
	if len(length) > 0 {
		end := start + length[0]
		if end > len(s) {
			end = len(s)
		}
		if end <= start {
			return ""
		}
		return s[start:end]
	}
	return s[start:]
}
