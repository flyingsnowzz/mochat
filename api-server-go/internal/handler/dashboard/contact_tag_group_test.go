package dashboard

import (
	"net/http"
	"testing"

	"mochat-api-server/internal/model"
	"mochat-api-server/tests/utils"
)

func TestWorkContactTagGroupHandler_Index(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.WorkContactTagGroup{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	handler := NewWorkContactTagGroupHandler(db)

	// 测试 Index 方法
	c, w := utils.SetupTestGin("GET", "/workContactTagGroup/index", "")
	c.Set("corpId", uint(1))
	handler.Index(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}

func TestWorkContactTagGroupHandler_Detail(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.WorkContactTagGroup{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// 创建测试数据
	group := &model.WorkContactTagGroup{
		GroupName: "测试标签组",
		CorpID:    1,
	}
	db.Create(group)

	handler := NewWorkContactTagGroupHandler(db)

	// 测试 Detail 方法
	c, w := utils.SetupTestGin("GET", "/workContactTagGroup/detail/1", "")
	handler.Detail(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}
