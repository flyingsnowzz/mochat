package service

import (
	"testing"

	"mochat-api-server/internal/model"
	"mochat-api-server/tests/utils"
)

func TestUserService(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	svc := NewUserService(db)

	// Test Create
	user := &model.User{
		Name:     "Test User",
		Phone:    "13800138000",
		Password: "password123",
		TenantID: 1,
		Status:   1,
	}

	err = svc.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}

	// Test GetByID
	retrievedUser, err := svc.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}

	// Test GetByPhone
	phoneUser, err := svc.GetByPhone("13800138000")
	if err != nil {
		t.Fatalf("Failed to get user by phone: %v", err)
	}

	if phoneUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, phoneUser.ID)
	}

	// Test List
	users, total, err := svc.List(1, 1, 10)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected total 1, got %d", total)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Test Update
	retrievedUser.Name = "Updated User"
	err = svc.Update(retrievedUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	updatedUser, err := svc.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updatedUser.Name != "Updated User" {
		t.Errorf("Expected user name 'Updated User', got %s", updatedUser.Name)
	}

	// Test UpdateStatus
	err = svc.UpdateStatus(user.ID, 0)
	if err != nil {
		t.Fatalf("Failed to update user status: %v", err)
	}

	statusUser, err := svc.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user after status update: %v", err)
	}

	if statusUser.Status != 0 {
		t.Errorf("Expected user status 0, got %d", statusUser.Status)
	}
}

func TestCorpService(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &model.Corp{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	svc := NewCorpService(db)

	// Test Create
	corp := &model.Corp{
		Name:     "Test Corp",
		TenantID: 1,
	}

	err = svc.Create(corp)
	if err != nil {
		t.Fatalf("Failed to create corp: %v", err)
	}

	if corp.ID == 0 {
		t.Error("Corp ID should be set after creation")
	}

	// Test GetByID
	retrievedCorp, err := svc.GetByID(corp.ID)
	if err != nil {
		t.Fatalf("Failed to get corp by ID: %v", err)
	}

	if retrievedCorp.ID != corp.ID {
		t.Errorf("Expected corp ID %d, got %d", corp.ID, retrievedCorp.ID)
	}

	// Test List
	corps, total, err := svc.List(1, 1, 10)
	if err != nil {
		t.Fatalf("Failed to list corps: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected total 1, got %d", total)
	}

	if len(corps) != 1 {
		t.Errorf("Expected 1 corp, got %d", len(corps))
	}

	// Test Select
	selectCorps, err := svc.Select(1)
	if err != nil {
		t.Fatalf("Failed to select corps: %v", err)
	}

	if len(selectCorps) != 1 {
		t.Errorf("Expected 1 corp in select, got %d", len(selectCorps))
	}

	// Test Update
	retrievedCorp.Name = "Updated Corp"
	err = svc.Update(retrievedCorp)
	if err != nil {
		t.Fatalf("Failed to update corp: %v", err)
	}

	updatedCorp, err := svc.GetByID(corp.ID)
	if err != nil {
		t.Fatalf("Failed to get updated corp: %v", err)
	}

	if updatedCorp.Name != "Updated Corp" {
		t.Errorf("Expected corp name 'Updated Corp', got %s", updatedCorp.Name)
	}
}
