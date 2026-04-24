package model

import (
	"testing"

	"gorm.io/gorm"
	"mochat-api-server/tests/utils"
)

func TestUserCRUD(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Test Create
	user := &User{
		Name:     "Test User",
		Phone:    "13800138000",
		Password: "password123",
		TenantID: 1,
		Status:   1,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}

	// Test Read
	var retrievedUser User
	if err := db.First(&retrievedUser, user.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrievedUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, retrievedUser.Name)
	}

	// Test Update
	retrievedUser.Name = "Updated User"
	if err := db.Save(&retrievedUser).Error; err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	var updatedUser User
	if err := db.First(&updatedUser, user.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}

	if updatedUser.Name != "Updated User" {
		t.Errorf("Expected updated user name 'Updated User', got %s", updatedUser.Name)
	}

	// Test Delete
	if err := db.Delete(&updatedUser).Error; err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	var deletedUser User
	result := db.First(&deletedUser, user.ID)
	if result.Error != gorm.ErrRecordNotFound {
		t.Error("Expected user to be deleted")
	}
}

func TestCorpCRUD(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &Corp{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Test Create
	corp := &Corp{
		Name:     "Test Corp",
		TenantID: 1,
	}

	if err := db.Create(corp).Error; err != nil {
		t.Fatalf("Failed to create corp: %v", err)
	}

	if corp.ID == 0 {
		t.Error("Corp ID should be set after creation")
	}

	// Test Read
	var retrievedCorp Corp
	if err := db.First(&retrievedCorp, corp.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve corp: %v", err)
	}

	if retrievedCorp.Name != corp.Name {
		t.Errorf("Expected corp name %s, got %s", corp.Name, retrievedCorp.Name)
	}
}

func TestWorkContactCRUD(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &WorkContact{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Test Create
	contact := &WorkContact{
		Name:             "Test Contact",
		WxExternalUserID: "external123",
		CorpID:           1,
	}

	if err := db.Create(contact).Error; err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	if contact.ID == 0 {
		t.Error("Contact ID should be set after creation")
	}

	// Test Read
	var retrievedContact WorkContact
	if err := db.First(&retrievedContact, contact.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve contact: %v", err)
	}

	if retrievedContact.Name != contact.Name {
		t.Errorf("Expected contact name %s, got %s", contact.Name, retrievedContact.Name)
	}
}

func TestWorkRoomCRUD(t *testing.T) {
	db := utils.SetupTestDB(t)

	// 自动迁移模型
	err := utils.AutoMigrateModels(db, &WorkRoom{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Test Create
	room := &WorkRoom{
		Name:     "Test Room",
		WxChatID: "room123",
		CorpID:   1,
		Status:   1,
	}

	if err := db.Create(room).Error; err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	if room.ID == 0 {
		t.Error("Room ID should be set after creation")
	}

	// Test Read
	var retrievedRoom WorkRoom
	if err := db.First(&retrievedRoom, room.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve room: %v", err)
	}

	if retrievedRoom.Name != room.Name {
		t.Errorf("Expected room name %s, got %s", room.Name, retrievedRoom.Name)
	}
}
