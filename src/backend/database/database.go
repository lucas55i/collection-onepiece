package database

import (
	"fmt"
	"log"
	"os"

	"github.com/collection-onepiece/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect inicializa a conexão com o PostgreSQL e executa as migrations
func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "onepiece"),
		getEnv("DB_PASSWORD", "onepiece"),
		getEnv("DB_NAME", "onepiece"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Auto migration — cria/atualiza a tabela volumes
	if err := db.AutoMigrate(&models.Volume{}); err != nil {
		log.Fatalf("Erro ao executar migration: %v", err)
	}

	DB = db
	log.Println("Banco de dados conectado e migrations aplicadas.")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
