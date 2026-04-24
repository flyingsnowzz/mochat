package dashboard

import (
	"net/http"
	"testing"

	"mochat-api-server/tests/utils"
)

func TestDashboardIndexHandler(t *testing.T) {
	db := utils.SetupTestDB(t)

	// Setup test context
	c, w := utils.SetupTestGin("GET", "/dashboard", "")

	// Create handler
	handler := NewDashboardIndexHandler(db)

	// Test Index method
	handler.Index(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	// Test LineChat method
	handler.LineChat(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}

func TestCorpHandler(t *testing.T) {
	db := utils.SetupTestDB(t)
	handler := NewCorpHandler(db)

	// Test Index method
	c, w := utils.SetupTestGin("GET", "/corps", "")
	c.Set("tenantId", uint(1))

	handler.Index(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	// Test Select method
	c, w = utils.SetupTestGin("GET", "/corps/select", "")
	c.Set("tenantId", uint(1))

	handler.Select(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}
