package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/collection-onepiece/backend/database"
	"github.com/collection-onepiece/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Erro ao abrir banco de teste: %v", err)
	}
	if err := db.AutoMigrate(&models.Volume{}); err != nil {
		t.Fatalf("Erro ao migrar banco de teste: %v", err)
	}
	database.DB = db
}

func TestGetVolumes_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/volumes", nil)

	GetVolumes(c)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, recebido %d", w.Code)
	}
}

func TestGetVolumes_WithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	// Insere volumes de teste
	database.DB.Create(&models.Volume{VolumeNumber: 1, Title: "Romance Dawn", Collected: false})
	database.DB.Create(&models.Volume{VolumeNumber: 2, Title: "Buggy the Clown", Collected: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/volumes", nil)

	GetVolumes(c)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, recebido %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Romance Dawn") {
		t.Error("resposta não contém 'Romance Dawn'")
	}
}

func TestUpdateVolume_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	database.DB.Create(&models.Volume{VolumeNumber: 1, Title: "Romance Dawn", Collected: false})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/api/volumes/1", strings.NewReader(`{"collected":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	UpdateVolume(c)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, recebido %d", w.Code)
	}
	// Verifica que acquired_at foi preenchido automaticamente
	if !strings.Contains(w.Body.String(), "acquired_at") {
		t.Error("resposta não contém 'acquired_at'")
	}
}

func TestUpdateVolume_ClearsAcquiredAt_WhenUncollected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	database.DB.Create(&models.Volume{VolumeNumber: 1, Title: "Romance Dawn", Collected: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/api/volumes/1", strings.NewReader(`{"collected":false}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	UpdateVolume(c)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, recebido %d", w.Code)
	}
	// acquired_at deve ser null
	if strings.Contains(w.Body.String(), `"acquired_at":"`) {
		t.Error("acquired_at deveria ser null após desmarcar")
	}
}

func TestUpdateVolume_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/api/volumes/999", strings.NewReader(`{"collected":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	UpdateVolume(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, recebido %d", w.Code)
	}
}

func TestUpdateVolume_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/api/volumes/1", strings.NewReader(`invalid json`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	UpdateVolume(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebido %d", w.Code)
	}
}
