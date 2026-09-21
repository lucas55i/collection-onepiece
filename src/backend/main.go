package main

import (
	"log"
	"os"

	"github.com/collection-onepiece/backend/database"
	"github.com/collection-onepiece/backend/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Conecta ao banco e aplica migrations
	database.Connect()

	// Configura o router
	r := gin.Default()

	// CORS — permite chamadas do frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	// Rotas da API
	api := r.Group("/api")
	{
		api.GET("/volumes", handlers.GetVolumes)
		api.PATCH("/volumes/:id", handlers.UpdateVolume)
		api.POST("/sync", handlers.SyncVolumes)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
