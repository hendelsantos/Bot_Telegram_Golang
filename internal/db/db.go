package db

import (
"log"
"os"

"gorm.io/driver/sqlite"
"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) {
var err error

// Para Railway (produção), tratar erros de CGO graciosamente
if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
log.Println("Detectado ambiente Railway")
// Tenta primeiro com arquivo normal
DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
if err != nil {
log.Printf("Erro ao criar banco de arquivo, usando in-memory: %v", err)
// Fallback para in-memory
DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
if err != nil {
log.Fatalf("Erro crítico ao conectar ao banco: %v", err)
}
}
} else {
// Desenvolvimento local - usa arquivo
DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
if err != nil {
log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
}
}

// Migra as tabelas
if err := AutoMigrate(DB); err != nil {
log.Fatalf("Erro ao migrar banco de dados: %v", err)
}
if err := AutoMigrateHistorico(DB); err != nil {
log.Fatalf("Erro ao migrar tabela de histórico: %v", err)
}

log.Println("Banco de dados inicializado com sucesso!")
}
