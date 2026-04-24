package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB 创建内存数据库用于测试
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

// AutoMigrateModels 自动迁移模型
func AutoMigrateModels(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// SetupTestGin 创建测试用的 Gin 上下文
func SetupTestGin(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, path, nil)
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
}

// SetupTestContext 创建测试用的上下文
func SetupTestContext() context.Context {
	return context.Background()
}
