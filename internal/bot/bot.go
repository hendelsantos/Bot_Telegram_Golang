package bot

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/modules"
)

// generateSessionID gera um ID único para cada sessão do bot
func generateSessionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// StartWithSmartRetry inicia o bot com retry inteligente
func Start(bot *tgbotapi.BotAPI) {
	// Limpa webhook primeiro
	deleteWebhookConfig := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	}
	_, err := bot.Request(deleteWebhookConfig)
	if err != nil {
		log.Printf("Erro ao deletar webhook: %v", err)
	}
	
	// Aguarda estabilização
	time.Sleep(3 * time.Second)
	
	sessionID := generateSessionID()
	log.Printf("Iniciando nova sessão com ID: %s", sessionID)
	
	retryCount := 0
	maxRetries := 10
	baseDelay := 5 * time.Second
	
	for retryCount < maxRetries {
		log.Printf("Tentativa %d/%d de conexão", retryCount+1, maxRetries)
		
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30
		updateConfig.Limit = 5
		
		// Tenta obter updates
		updates := bot.GetUpdatesChan(updateConfig)
		
		// Processa mensagens
		conflictDetected := false
		for update := range updates {
			if update.Message == nil {
				continue
			}
			
			// Verifica se é um erro de conflito no log
			log.Printf("Mensagem recebida de %s: %s", update.Message.From.UserName, update.Message.Text)
			
			// Processa normalmente
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					modules.HandleStart(bot, update)
				case "menu":
					modules.HandleMenu(bot, update)
				case "novoitem":
					modules.HandleNovoItem(bot, update)
				case "buscar":
					modules.HandleBuscar(bot, update)
				case "atualizar":
					modules.HandleAtualizar(bot, update)
				case "enviar_reparo":
					modules.HandleEnviarReparo(bot, update)
				case "retornar_reparo":
					modules.HandleRetornarReparo(bot, update)
				case "exportar_estoque":
					modules.HandleExportarEstoque(bot, update)
				case "historico":
					modules.HandleHistorico(bot, update)
				// Novos comandos de listagem
				case "listar", "listar_resumo", "listar_status", "listar_baixo_estoque", "listar_detalhado":
					modules.HandleListar(bot, update)
				case "ajuda_listagem":
					modules.HandleAjudaListagem(bot, update)
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Comando não reconhecido. Use /menu para ver todos os comandos disponíveis.")
					bot.Send(msg)
				}
			} else {
				// Verifica em qual fluxo o usuário está
				if modules.IsUserInCadastroFlow(update.Message.From.ID) {
					modules.HandleNovoItem(bot, update)
				} else if modules.IsUserInAtualizacaoFlow(update.Message.From.ID) {
					modules.HandleAtualizar(bot, update)
				} else if modules.IsUserInReparoFlow(update.Message.From.ID) {
					modules.HandleEnviarReparo(bot, update)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use /menu para ver todos os comandos disponíveis.")
					bot.Send(msg)
				}
			}
		}
		
		// Se chegou aqui, houve erro na conexão
		if !conflictDetected {
			retryCount++
			delay := time.Duration(retryCount) * baseDelay
			log.Printf("Conexão perdida, aguardando %v antes da próxima tentativa...", delay)
			time.Sleep(delay)
			
			// Tenta limpar novamente
			bot.Request(deleteWebhookConfig)
			time.Sleep(2 * time.Second)
		}
	}
	
	log.Printf("Máximo de tentativas (%d) atingido. Encerrando...", maxRetries)
}
