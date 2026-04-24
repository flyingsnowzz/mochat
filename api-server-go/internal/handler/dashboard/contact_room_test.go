package dashboard

import (
	"net/http"
	"testing"

	"mochat-api-server/internal/model"
	"mochat-api-server/tests/utils"
)

func TestWorkContactRoomHandler_Index(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.WorkContactRoom{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	handler := NewWorkContactRoomHandler(db)

	// 测试 Index 方法
	c, w := utils.SetupTestGin("GET", "/workContactRoom/index?roomId=1", "")
	handler.Index(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}
