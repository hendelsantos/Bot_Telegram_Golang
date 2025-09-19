package modules

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "botgo/internal/db"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleExportarEstoque(bot *tgbotapi.BotAPI, chatID int64) {
    // Buscar todos os itens
    var itens []db.Item
    result := db.DB.Find(&itens)
    if result.Error != nil {
        log.Printf("Erro ao buscar itens para exportação: %v", result.Error)
        msg := tgbotapi.NewMessage(chatID, "❌ Erro ao buscar itens para exportação.")
        bot.Send(msg)
        return
    }

    // Criar arquivo CSV temporário
    fileName := fmt.Sprintf("estoque_%d.csv", time.Now().Unix())
    file, err := os.Create(fileName)
    if err != nil {
        log.Printf("Erro ao criar arquivo CSV: %v", err)
        msg := tgbotapi.NewMessage(chatID, "❌ Erro ao criar arquivo de exportação.")
        bot.Send(msg)
        return
    }
    defer file.Close()
    defer os.Remove(fileName) // Limpar arquivo após envio

    // Criar writer CSV
    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Escrever cabeçalho
    header := []string{"ID", "Nome", "Descrição", "Quantidade", "Status", "Fornecedor", "Data Envio"}
    writer.Write(header)

    // Escrever dados
    for _, item := range itens {
        record := []string{
            strconv.Itoa(int(item.ID)),
            item.Nome,
            item.Descricao,
            strconv.Itoa(item.Quantidade),
            item.Status,
            item.Fornecedor,
            formatDataEnvio(item.DataEnvio),
        }
        writer.Write(record)
    }

    // Enviar arquivo
    doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(fileName))
    doc.Caption = fmt.Sprintf("📊 Exportação do Estoque\n\n📋 Total de itens: %d\n📅 Gerado em: %s", 
        len(itens), time.Now().Format("02/01/2006 15:04"))
    
    _, err = bot.Send(doc)
    if err != nil {
        log.Printf("Erro ao enviar arquivo CSV: %v", err)
        msg := tgbotapi.NewMessage(chatID, "❌ Erro ao enviar arquivo de exportação.")
        bot.Send(msg)
        return
    }

    msg := tgbotapi.NewMessage(chatID, "✅ Estoque exportado com sucesso!")
    bot.Send(msg)
}

// Função helper para formatar data de envio
func formatDataEnvio(dataEnvio *time.Time) string {
    if dataEnvio == nil {
        return ""
    }
    return dataEnvio.Format("02/01/2006")
}
