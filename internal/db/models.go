package db

import (
	"gorm.io/gorm"
)

type Item struct {
	ID          uint   `gorm:"primaryKey"`
	Nome        string
	Descricao   string
	Quantidade  int
	Status      string
	FotoPath    string
	Fornecedor  string
	DataEnvio   string
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Item{})
}
