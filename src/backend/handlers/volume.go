package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/collection-onepiece/backend/database"
	"github.com/collection-onepiece/backend/models"
	"github.com/gin-gonic/gin"
)

// jikanVolume estrutura de resposta da Jikan API
type jikanVolume struct {
	MalID   int    `json:"mal_id"`
	Title   string `json:"title"`
	Images  struct {
		JPG struct {
			ImageURL string `json:"image_url"`
		} `json:"jpg"`
	} `json:"images"`
}

type jikanResponse struct {
	Data []jikanVolume `json:"data"`
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
	} `json:"pagination"`
}

// updateVolumeRequest é o payload de entrada para PATCH /api/volumes/:id.
// Usa ponteiros para distinguir ausência de zero value:
//   - Collected *bool   → nil indica campo ausente (HTTP 400)
//   - AcquiredAt *string → string para validar timezone antes do parse para time.Time
type updateVolumeRequest struct {
	Collected  *bool   `json:"collected"`
	AcquiredAt *string `json:"acquired_at"`
}

// iso8601WithTZ lista os layouts aceitos — timezone obrigatório.
var iso8601WithTZ = []string{
	time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
}

// parseISO8601WithTZ faz o parse de uma string ISO 8601 com timezone obrigatório.
// Retorna erro se a string não puder ser parseada em nenhum dos formatos aceitos.
func parseISO8601WithTZ(s string) (*time.Time, error) {
	for _, layout := range iso8601WithTZ {
		if t, err := time.Parse(layout, s); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("formato inválido: esperado ISO 8601 com timezone")
}

// GetVolumes retorna todos os volumes do banco de dados
func GetVolumes(c *gin.Context) {
	var volumes []models.Volume
	result := database.DB.Order("volume_number asc").Find(&volumes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar volumes"})
		return
	}
	c.JSON(http.StatusOK, volumes)
}

// UpdateVolume atualiza os campos collected e acquired_at de um volume.
//
// Fluxo de decisão para acquired_at:
//   - collected=false → acquired_at = NULL
//   - collected=true + acquired_at explícito válido → usa valor fornecido
//   - collected=true + volume já tem acquired_at no DB → preserva existente
//   - collected=true + volume sem acquired_at → timestamp UTC do servidor
func UpdateVolume(c *gin.Context) {
	id := c.Param("id")

	var req updateVolumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON malformado ou inválido"})
		return
	}
	if req.Collected == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campo 'collected' é obrigatório"})
		return
	}

	// Valida acquired_at se presente no payload
	var acquiredAt *time.Time
	if req.AcquiredAt != nil {
		parsed, err := parseISO8601WithTZ(*req.AcquiredAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "acquired_at deve ser ISO 8601 com timezone"})
			return
		}
		acquiredAt = parsed
	}

	// Busca o volume existente para preservação de acquired_at e validação de existência
	var existing models.Volume
	if err := database.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volume não encontrado"})
		return
	}

	updates := map[string]interface{}{"collected": *req.Collected, "acquired_at": nil}

	if *req.Collected {
		switch {
		case acquiredAt != nil:
			// Valor explícito fornecido no payload
			updates["acquired_at"] = acquiredAt
		case existing.AcquiredAt != nil:
			// Preserva o valor já existente no banco
			updates["acquired_at"] = existing.AcquiredAt
		default:
			// Primeiro registro: usa timestamp UTC do servidor
			now := time.Now().UTC()
			updates["acquired_at"] = &now
		}
	}

	database.DB.Model(&existing).Updates(updates)

	// Re-fetch para garantir estado atualizado refletido na resposta
	database.DB.First(&existing, id)
	c.JSON(http.StatusOK, existing)
}

// SyncVolumes busca volumes da Jikan API e salva no banco (idempotente)
func SyncVolumes(c *gin.Context) {
	count, err := syncFromJikan()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d volumes sincronizados", count)})
}

// syncFromJikan busca todos os volumes do One Piece na Jikan API
func syncFromJikan() (int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	// ID do One Piece no MyAnimeList é 13
	const mangaID = 13
	url := fmt.Sprintf("https://api.jikan.moe/v4/manga/%d/full", mangaID)

	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao chamar Jikan API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var result struct {
		Data struct {
			MalID   int    `json:"mal_id"`
			Title   string `json:"title"`
			Volumes int    `json:"volumes"`
			Images  struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	totalVolumes := result.Data.Volumes
	if totalVolumes == 0 {
		totalVolumes = 114 // fallback com o último valor conhecido
	}

	count := 0
	for i := 1; i <= totalVolumes; i++ {
		volume := models.Volume{
			VolumeNumber: i,
			Title:        fmt.Sprintf("Volume %d", i),
			CoverImage:   result.Data.Images.JPG.ImageURL,
			Collected:    false,
		}

		// Upsert: insere apenas se não existir, preserva collected
		tx := database.DB.Where(models.Volume{VolumeNumber: i}).FirstOrCreate(&volume)
		if tx.Error != nil {
			log.Printf("Erro ao salvar volume %d: %v", i, tx.Error)
			continue
		}
		count++
	}

	return count, nil
}
