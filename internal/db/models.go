package db

import (
    "time"
    "gorm.io/gorm"
)

type Item struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    Nome        string    `json:"nome"`
    Descricao   string    `json:"descricao"`
    Quantidade  int       `json:"quantidade"`
    Status      string    `json:"status"` // "Em Estoque", "Em Reparo Externo"
    FotoPath    string    `json:"foto_path"`
    Fornecedor  string    `json:"fornecedor"`  // Para reparos externos
    DataEnvio   *time.Time `json:"data_envio"` // Para reparos externos
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Movimentacao struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    ItemID    uint      `json:"item_id"`
    Item      Item      `gorm:"foreignKey:ItemID" json:"item"`
    Tipo      string    `json:"tipo"` // "cadastro", "atualizacao", "envio_reparo", "retorno_reparo"
    Descricao string    `json:"descricao"`
    DataHora  time.Time `json:"data_hora"`
    CreatedAt time.Time `json:"created_at"`
}
