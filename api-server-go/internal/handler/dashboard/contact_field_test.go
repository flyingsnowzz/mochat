package dashboard

import (
	"net/http"
	"testing"

	"mochat-api-server/internal/model"
	"mochat-api-server/tests/utils"
)

func TestContactFieldHandler_Index(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.ContactField{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// 创建测试数据
	field := &model.ContactField{
		Name:   "测试字段",
		Label:  "测试标签",
		Type:   1,
		Status: 1,
	}
	db.Create(field)

	handler := NewContactFieldHandler(db)

	// 测试 Index 方法
	c, w := utils.SetupTestGin("GET", "/contactField/index", "")
	handler.Index(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}

func TestContactFieldHandler_Show(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.ContactField{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// 创建测试数据
	field := &model.ContactField{
		Name:   "测试字段",
		Label:  "测试标签",
		Type:   1,
		Status: 1,
	}
	db.Create(field)

	handler := NewContactFieldHandler(db)

	// 测试 Show 方法
	c, w := utils.SetupTestGin("GET", "/contactField/show/1", "")
	handler.Show(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}
