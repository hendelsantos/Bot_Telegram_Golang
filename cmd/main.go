package main

import (
	"log"
	"os"
	"path/filepath"

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

	// Define caminho do banco de dados
	dbPath := "estoque.db"
	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		dbPath = filepath.Join(dataDir, "estoque.db")
	}

	// Inicializa banco de dados
	db.InitDB(dbPath)

	// Cria bot
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Debug apenas em desenvolvimento
	b.Debug = os.Getenv("DEBUG") == "true"

	log.Printf("Bot autorizado em: %s", b.Self.UserName)

	// Limpa qualquer webhook existente para evitar conflitos
	deleteWebhookConfig := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	}
	_, err = b.Request(deleteWebhookConfig)
	if err != nil {
		log.Printf("Erro ao deletar webhook: %v", err)
	}

	log.Println("Bot iniciando...")
	bot.Start(b)
}
