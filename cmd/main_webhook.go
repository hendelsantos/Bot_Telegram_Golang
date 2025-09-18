package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/bot"
	"botgo/internal/db"
)

func main() {
	// Token do bot
	token := "8320983858:AAG8ID2e7ReRk81YFEF3hJiqvD9HCDvLNlU"
	if envToken := os.Getenv("TELEGRAM_BOT_TOKEN"); envToken != "" {
		token = envToken
	}

	dbPath := "estoque.db"
	db.InitDB(dbPath)

	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	b.Debug = true

	log.Printf("Bot autorizado em: %s", b.Self.UserName)

	// Configuração para webhook (alternativa ao polling)
	useWebhook := os.Getenv("USE_WEBHOOK") == "true"
	
	if useWebhook {
		// Configurar webhook
		webhookURL := os.Getenv("WEBHOOK_URL") // Ex: https://seudominio.com/webhook
		if webhookURL == "" {
			log.Fatal("WEBHOOK_URL não configurada")
		}
		
		webhookConfig := tgbotapi.NewWebhook(webhookURL)
		_, err = b.Request(webhookConfig)
		if err != nil {
			log.Fatal(err)
		}
		
		info, err := b.GetWebhookInfo()
		if err != nil {
			log.Fatal(err)
		}
		
		if info.LastErrorDate != 0 {
			log.Printf("Erro no webhook: %s", info.LastErrorMessage)
		}
		
		// Configurar handler HTTP para receber webhooks
		http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
			update, err := b.HandleUpdate(r)
			if err != nil {
				log.Printf("Erro no webhook: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			// Processar update usando a mesma lógica
			bot.ProcessUpdate(b, update)
		})
		
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		
		log.Printf("Iniciando servidor webhook na porta %s", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
		
	} else {
		// Usar polling (método atual com melhorias)
		bot.StartWithImprovedPolling(b)
	}
}
