package modules

import (
    "fmt"
    "log"
    "strconv"
    "strings"
    "time"

    "botgo/internal/db"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserStateReparo struct {
    ItemID      uint
    Fornecedor  string
    DataEnvio   string
    Step        string
}

var userStateReparo = make(map[int64]*UserStateReparo)

func HandleEnviarReparo(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    chatID := message.Chat.ID
    args := strings.Fields(message.Text)
    
    if len(args) < 2 {
        msg := tgbotapi.NewMessage(chatID, "❌ Use: /enviar_reparo <ID>")
        bot.Send(msg)
        return
    }
    
    id, err := strconv.Atoi(args[1])
    if err != nil {
        msg := tgbotapi.NewMessage(chatID, "❌ ID inválido.")
        bot.Send(msg)
        return
    }
    
    // Verificar se item existe
    var item db.Item
    result := db.DB.First(&item, id)
    if result.Error != nil {
        msg := tgbotapi.NewMessage(chatID, "❌ Item não encontrado.")
        bot.Send(msg)
        return
    }
    
    // Verificar se já está em reparo
    if item.Status == "Em Reparo Externo" {
        msg := tgbotapi.NewMessage(chatID, "⚠️ Este item já está em reparo externo.")
        bot.Send(msg)
        return
    }
    
    // Iniciar fluxo
    userStateReparo[chatID] = &UserStateReparo{
        ItemID: item.ID,
        Step:   "fornecedor",
    }
    
    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🔧 Enviando item para reparo:\n\n📦 **%s** (ID: %d)\n\n👤 Para qual fornecedor/local será enviado?", item.Nome, item.ID))
    msg.ParseMode = "Markdown"
    bot.Send(msg)
}

func HandleRetornarReparo(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    chatID := message.Chat.ID
    args := strings.Fields(message.Text)
    
    if len(args) < 2 {
        msg := tgbotapi.NewMessage(chatID, "❌ Use: /retornar_reparo <ID>")
        bot.Send(msg)
        return
    }
    
    id, err := strconv.Atoi(args[1])
    if err != nil {
        msg := tgbotapi.NewMessage(chatID, "❌ ID inválido.")
        bot.Send(msg)
        return
    }
    
    // Verificar se item existe
    var item db.Item
    result := db.DB.First(&item, id)
    if result.Error != nil {
        msg := tgbotapi.NewMessage(chatID, "❌ Item não encontrado.")
        bot.Send(msg)
        return
    }
    
    // Verificar se está em reparo
    if item.Status != "Em Reparo Externo" {
        msg := tgbotapi.NewMessage(chatID, "⚠️ Este item não está em reparo externo.")
        bot.Send(msg)
        return
    }
    
    // Atualizar item
    item.Status = "Em Estoque"
    item.Quantidade++
    item.Fornecedor = ""
    item.DataEnvio = nil
    
    err = db.DB.Save(&item).Error
    if err != nil {
        log.Printf("Erro ao retornar item do reparo: %v", err)
        msg := tgbotapi.NewMessage(chatID, "❌ Erro ao processar retorno do reparo.")
        bot.Send(msg)
        return
    }
    
    // Registrar movimentação
    movimentacao := db.Movimentacao{
        ItemID:    item.ID,
        Tipo:      "retorno_reparo",
        Descricao: fmt.Sprintf("Item retornado do reparo externo"),
        DataHora:  time.Now(),
    }
    db.DB.Create(&movimentacao)
    
    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ **%s** retornado do reparo!\n\n📦 Quantidade em estoque: %d\n📋 Status: %s", item.Nome, item.Quantidade, item.Status))
    msg.ParseMode = "Markdown"
    bot.Send(msg)
}

func ProcessReparoFlow(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
    chatID := message.Chat.ID
    state, exists := userStateReparo[chatID]
    if !exists {
        return false
    }
    
    switch state.Step {
    case "fornecedor":
        state.Fornecedor = message.Text
        state.Step = "data"
        
        msg := tgbotapi.NewMessage(chatID, "📅 Qual a data de envio? (formato: DD/MM/AAAA)")
        bot.Send(msg)
        return true
        
    case "data":
        state.DataEnvio = message.Text
        
        // Validar e converter data
        dataEnvio, err := time.Parse("02/01/2006", state.DataEnvio)
        if err != nil {
            msg := tgbotapi.NewMessage(chatID, "❌ Data inválida. Use o formato DD/MM/AAAA")
            bot.Send(msg)
            return true
        }
        
        // Buscar item
        var item db.Item
        result := db.DB.First(&item, state.ItemID)
        if result.Error != nil {
            msg := tgbotapi.NewMessage(chatID, "❌ Erro ao buscar item.")
            bot.Send(msg)
            delete(userStateReparo, chatID)
            return true
        }
        
        // Atualizar item
        item.Status = "Em Reparo Externo"
        item.Quantidade--
        item.Fornecedor = state.Fornecedor
        item.DataEnvio = &dataEnvio
        
        err = db.DB.Save(&item).Error
        if err != nil {
            log.Printf("Erro ao enviar item para reparo: %v", err)
            msg := tgbotapi.NewMessage(chatID, "❌ Erro ao processar envio para reparo.")
            bot.Send(msg)
            delete(userStateReparo, chatID)
            return true
        }
        
        // Registrar movimentação
        movimentacao := db.Movimentacao{
            ItemID:    item.ID,
            Tipo:      "envio_reparo",
            Descricao: fmt.Sprintf("Item enviado para reparo externo - Fornecedor: %s", state.Fornecedor),
            DataHora:  time.Now(),
        }
        db.DB.Create(&movimentacao)
        
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ **%s** enviado para reparo!\n\n👤 **Fornecedor:** %s\n📅 **Data de envio:** %s\n📦 **Quantidade em estoque:** %d\n📋 **Status:** %s", 
            item.Nome, state.Fornecedor, state.DataEnvio, item.Quantidade, item.Status))
        msg.ParseMode = "Markdown"
        bot.Send(msg)
        
        delete(userStateReparo, chatID)
        return true
    }
    
    return false
}

func IsUserInReparoFlow(chatID int64) bool {
    _, exists := userStateReparo[chatID]
    return exists
}

func ClearReparoState(chatID int64) {
    delete(userStateReparo, chatID)
}
