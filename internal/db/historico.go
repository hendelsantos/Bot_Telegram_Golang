package db

import (
	"gorm.io/gorm"
	"time"
)

type Historico struct {
	ID        uint      `gorm:"primaryKey"`
	ItemID    uint
	Acao      string
	Descricao string
	DataHora  time.Time
}

func AutoMigrateHistorico(db *gorm.DB) error {
	return db.AutoMigrate(&Historico{})
}

func RegistrarHistorico(itemID uint, acao, descricao string) {
	db := DB
	h := Historico{
		ItemID:    itemID,
		Acao:      acao,
		Descricao: descricao,
		DataHora:  time.Now(),
	}
	db.Create(&h)
}
