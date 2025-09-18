package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/bot"
	"botgo/internal/db"
)


func main() {
	// Carrega token da variável de ambiente
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN não definida. Configure a variável de ambiente.")
	}

	dbPath := "estoque.db"
	db.InitDB(dbPath)

	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	b.Debug = true

	log.Printf("Bot autorizado em: %s", b.Self.UserName)

	// Limpa qualquer webhook existente para evitar conflitos
	deleteWebhookConfig := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	}
	_, err = b.Request(deleteWebhookConfig)
	if err != nil {
		log.Printf("Erro ao deletar webhook: %v", err)
	}
	
	bot.Start(b)
}
