package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mochat-api-server/internal/config"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/logger"
	"mochat-api-server/internal/pkg/response"
	mochatRedis "mochat-api-server/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// DashboardClaims 管理后台 JWT Claims
type DashboardClaims struct {
	UserID   uint   `json:"userId"`
	Phone    string `json:"phone"`
	TenantID uint   `json:"tenantId"`
	jwt.RegisteredClaims
}

// SidebarClaims 侧边栏 JWT Claims
type SidebarClaims struct {
	EmployeeID uint `json:"employeeId"`
	CorpID     uint `json:"corpId"`
	jwt.RegisteredClaims
}

// GenerateDashboardToken 生成管理后台 JWT Token
func GenerateDashboardToken(userID uint, phone string, tenantID uint, secret string) (string, error) {
	claims := DashboardClaims{
		UserID:   userID,
		Phone:    phone,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateSidebarToken 生成侧边栏 JWT Token
func GenerateSidebarToken(employeeID uint, corpID uint, secret string) (string, error) {
	claims := SidebarClaims{
		EmployeeID: employeeID,
		CorpID:     corpID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// DashboardAuthMiddleware 管理后台 JWT 认证中间件
func DashboardAuthMiddleware(cfg config.JWTConfig, whiteRoutes []string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, route := range whiteRoutes {
			if matchRoute(path, route) {
				c.Next()
				return
			}
		}

		tokenStr := extractToken(c)
		if tokenStr == "" {
			response.FailWithHTTP(c, http.StatusUnauthorized, response.ErrAuth, "未登录")
			c.Abort()
			return
		}

		claims := &DashboardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.DashboardSecret), nil
		})

		if err != nil || !token.Valid {
			// 添加详细的错误信息
			fmt.Printf("Token 验证失败: err=%v, tokenStr=%s, secret=%s\n", err, tokenStr, cfg.DashboardSecret)
			response.FailWithHTTP(c, http.StatusUnauthorized, response.ErrTokenInvalid, "token无效或已过期")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("tenantId", claims.TenantID)

		// 从 Redis 获取用户当前选择的企业，如不存在则从 DB 回退查询
		corpID, employeeID := resolveCorpBinding(claims.UserID, db)
		if corpID > 0 {
			c.Set("corpId", corpID)
		}
		if employeeID > 0 {
			c.Set("employeeId", employeeID)
		}

		c.Set("guard", "dashboard")
		c.Next()
	}
}

// SidebarAuthMiddleware 侧边栏 JWT 认证中间件
func SidebarAuthMiddleware(cfg config.JWTConfig, whiteRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, route := range whiteRoutes {
			if matchRoute(path, route) {
				c.Next()
				return
			}
		}

		tokenStr := extractToken(c)
		if tokenStr == "" {
			response.FailWithHTTP(c, http.StatusUnauthorized, response.ErrAuth, "未登录")
			c.Abort()
			return
		}

		claims := &SidebarClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.SidebarSecret), nil
		})

		if err != nil || !token.Valid {
			response.FailWithHTTP(c, http.StatusUnauthorized, response.ErrTokenInvalid, "token无效或已过期")
			c.Abort()
			return
		}

		c.Set("employeeId", claims.EmployeeID)
		c.Set("corpId", claims.CorpID)
		c.Set("guard", "sidebar")
		c.Next()
	}
}

// extractToken 从请求头中提取 Bearer Token
func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	// 处理没有 Bearer 前缀的情况
	return auth
}

// matchRoute 匹配路由路径，支持通配符 *
func matchRoute(path string, pattern string) bool {
	if pattern == path {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	return false
}

// resolveCorpBinding 从 Redis 获取用户-企业绑定，缓存不存在时从 DB 回退查询并写入
func resolveCorpBinding(userID uint, db *gorm.DB) (uint, uint) {
	fmt.Printf("[resolveCorpBinding] start userID=%d\n", userID)
	logger.Sugar.Info("resolveCorpBinding start userID=%d", userID)
	logger.Sugar.Infof("resolveCorpBinding db=%v", db)
	logger.Sugar.Infof("resolveCorpBinding mochatRedis.RDB=%v", mochatRedis.RDB)
	if mochatRedis.RDB != nil {
		cacheKey := "mc:user." + strconv.Itoa(int(userID))
		if cached, err := mochatRedis.RDB.Get(context.Background(), cacheKey).Result(); err == nil && cached != "" {
			parts := strings.SplitN(cached, "-", 2)
			if len(parts) >= 1 {
				if corpID, parseErr := strconv.ParseUint(parts[0], 10, 64); parseErr == nil && corpID > 0 {
					if len(parts) == 2 {
						if employeeID, parseErr := strconv.ParseUint(parts[1], 10, 64); parseErr == nil && employeeID > 0 {
							fmt.Printf("[resolveCorpBinding] redis hit: corpId=%d employeeId=%d\n", corpID, employeeID)
							logger.Sugar.Infof("resolveCorpBinding redis hit: corpId=%d employeeId=%d", corpID, employeeID)
							return uint(corpID), uint(employeeID)
						}
					}
					fmt.Printf("[resolveCorpBinding] redis hit (no employee): corpId=%d\n", corpID)
					logger.Sugar.Infof("resolveCorpBinding redis hit (no employee): corpId=%d", corpID)
					return uint(corpID), 0
				}
			}
		}
	}
	fmt.Println("[resolveCorpBinding] redis miss, falling back to DB")
	logger.Sugar.Info("resolveCorpBinding redis miss, falling back to DB")

	// Redis 无缓存，从 DB 回退查询
	if db == nil {
		fmt.Println("[resolveCorpBinding] db is nil")
		return 0, 0
	}
	var employee model.WorkEmployee
	if err := db.Where("log_user_id = ?", userID).First(&employee).Error; err != nil {
		fmt.Printf("[resolveCorpBinding] DB query failed: %v\n", err)
		return 0, 0
	}
	fmt.Printf("[resolveCorpBinding] DB found employee: id=%d corpId=%d\n", employee.ID, employee.CorpID)

	if employee.CorpID == 0 || employee.ID == 0 {
		fmt.Println("[resolveCorpBinding] corpID or employeeID is 0")
		return 0, 0
	}

	// 写入 Redis 缓存
	if mochatRedis.RDB != nil {
		cacheKey := "mc:user." + strconv.Itoa(int(userID))
		cacheValue := strconv.Itoa(int(employee.CorpID)) + "-" + strconv.Itoa(int(employee.ID))
		if err := mochatRedis.RDB.Set(context.Background(), cacheKey, cacheValue, 0).Err(); err != nil {
			fmt.Printf("[resolveCorpBinding] Redis set failed: %v\n", err)
			logger.Sugar.Errorf("resolveCorpBinding Redis set failed: %v", err)
		} else {
			fmt.Printf("[resolveCorpBinding] Redis set success: %s = %s\n", cacheKey, cacheValue)
			logger.Sugar.Infof("resolveCorpBinding Redis set success: %s = %s", cacheKey, cacheValue)
		}
	}

	return employee.CorpID, employee.ID
}
