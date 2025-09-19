package main

import (
    "log"
    "os"

    "botgo/internal/bot"
    "botgo/internal/db"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    // Token do bot
    token := "8320983858:AAG8ID2e7ReRk81YFEF3hJiqvD9HCDvLNlU"
    if envToken := os.Getenv("TELEGRAM_BOT_TOKEN"); envToken != "" {
        token = envToken
    }

    if token == "8320983858:AAG8ID2e7ReRk81YFEF3hJiqvD9HCDvLNlU" || token == "seu_token_aqui" {
        log.Fatal("TELEGRAM_BOT_TOKEN não definido. Configure uma variável de ambiente.")
    }

    // Inicializar banco de dados (sem parâmetros)
    db.InitDB()

    // Criar bot
    b, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Panic(err)
    }

    // Desabilitar debug em produção
    if os.Getenv("DEBUG") != "true" {
        b.Debug = false
    }

    log.Printf("Bot autorizado em: %s", b.Self.UserName)

    // Iniciar bot
    botInstance := bot.NewBot(b)
    botInstance.Start()
}
