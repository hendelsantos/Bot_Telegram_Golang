package db

import (
    "log"
    "os"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    _ "modernc.org/sqlite" // Pure Go SQLite driver
)

var DB *gorm.DB

func InitDB() {
    var err error
    
    // Verificar se estamos no Railway (ambiente de produção)
    if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
        log.Println("Detectado ambiente Railway")
        // No Railway, usar banco em memória
        DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
        if err != nil {
            log.Printf("Erro crítico ao conectar ao banco: %v", err)
            os.Exit(1)
        }
        log.Println("Conectado ao banco: :memory:")
    } else {
        // Localmente, usar arquivo
        dbPath := os.Getenv("DB_PATH")
        if dbPath == "" {
            dbPath = "estoque.db"
        }
        
        DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
        if err != nil {
            log.Printf("Erro ao conectar ao banco local: %v", err)
            // Fallback para memória se arquivo falhar
            DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
            if err != nil {
                log.Printf("Erro crítico ao conectar ao banco: %v", err)
                os.Exit(1)
            }
            log.Println("Conectado ao banco: :memory: (fallback)")
        } else {
            log.Printf("Conectado ao banco: %s", dbPath)
        }
    }

    // Migrar as tabelas
    err = DB.AutoMigrate(&Item{}, &Movimentacao{})
    if err != nil {
        log.Printf("Erro ao migrar tabelas: %v", err)
        os.Exit(1)
    }
    
    log.Println("Database inicializado com sucesso")
}
