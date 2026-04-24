package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mochat-api-server/internal/pkg/logger"
)

// CoreMiddleware 核心中间件，记录请求开始时间
func CoreMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("startTime", time.Now())
		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Sugar.Infof("[%s] %s %d %v",
			method,
			path,
			statusCode,
			latency,
		)
	}
}

// RateLimitMiddleware 简易限流中间件（基于内存，单实例适用）
func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	type client struct {
		count    int
		lastTime time.Time
	}
	clients := make(map[string]*client)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		cl, exists := clients[ip]
		if !exists || now.Sub(cl.lastTime) > duration {
			clients[ip] = &client{count: 1, lastTime: now}
			c.Next()
			return
		}

		cl.count++
		cl.lastTime = now

		if cl.count > maxRequests {
			logger.Sugar.Warnw("rate limit exceeded",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)
			c.AbortWithStatusJSON(429, gin.H{
				"code": 429,
				"msg":  "请求过于频繁",
				"data": nil,
			})
			return
		}

		c.Next()
	}
}
