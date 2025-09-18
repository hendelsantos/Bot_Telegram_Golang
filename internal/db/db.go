package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB(path string) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	DB = database
	if err := AutoMigrate(DB); err != nil {
		log.Fatalf("Erro ao migrar banco de dados: %v", err)
	}
	if err := AutoMigrateHistorico(DB); err != nil {
		log.Fatalf("Erro ao migrar tabela de histórico: %v", err)
	}
}
