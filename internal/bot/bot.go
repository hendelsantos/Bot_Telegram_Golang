package bot

import (
	"log"
	"strings"
	"time"

	"botgo/internal/modules"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

// Start inicia o loop principal do bot usando polling.
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Aguardando mensagens...")
	for update := range updates {
		b.ProcessUpdate(update)
	}
}

// ProcessUpdate é o roteador central de todas as atualizações recebidas.
func (b *Bot) ProcessUpdate(update tgbotapi.Update) {
	// Ignora qualquer coisa que não seja uma mensagem
	if update.Message == nil {
		return
	}

	message := update.Message
	chatID := message.Chat.ID

	log.Printf("[%s] %s", message.From.UserName, message.Text)

	// Se for um comando, processa o comando.
	if message.IsCommand() {
		// Limpa qualquer estado de fluxo anterior ao receber um novo comando
		modules.ClearUserState(chatID)
		modules.ClearCadastroState(chatID)
		modules.ClearAtualizacaoState(chatID)
		modules.ClearReparoState(chatID)

		switch message.Command() {
		case "start", "menu", "help":
			modules.HandleMenu(b.api, chatID)
		case "novoitem":
			modules.HandleNovoItem(b.api, message)
		case "buscar":
			modules.HandleBuscar(b.api, message)
		case "atualizar":
			modules.HandleAtualizar(b.api, message)
		case "enviar_reparo":
			modules.HandleEnviarReparo(b.api, message)
		case "retornar_reparo":
			modules.HandleRetornarReparo(b.api, message)
		case "exportar_estoque":
			modules.HandleExportarEstoque(b.api, chatID)
		case "historico":
			modules.HandleHistorico(b.api, message)
		case "listar", "listar_resumo", "listar_detalhado", "listar_status", "listar_baixo_estoque":
			modules.HandleListagem(b.api, message)
		case "ajuda_listagem":
			modules.HandleAjudaListagem(b.api, chatID)
		default:
			msg := tgbotapi.NewMessage(chatID, "Comando não reconhecido. Use /menu para ver a lista de comandos.")
			b.api.Send(msg)
		}
		return
	}

	// Se não for um comando, verifica se o usuário está em algum fluxo de conversa.
	if modules.IsUserInCadastroFlow(chatID) {
		modules.ProcessCadastroFlow(b.api, message)
		return
	}
	if modules.IsUserInAtualizacaoFlow(chatID) {
		modules.ProcessAtualizacaoFlow(b.api, message)
		return
	}
	if modules.IsUserInReparoFlow(chatID) {
		modules.ProcessReparoFlow(b.api, message)
		return
	}

	// Se não for comando e não estiver em nenhum fluxo, ignora.
}

// StartWithImprovedPolling é uma alternativa com mais configurações.
func (b *Bot) StartWithImprovedPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	u.Limit = 5

	// Limpa webhooks pendentes
	_, _ = b.api.Request(tgbotapi.DeleteWebhookConfig{})

	updates := b.api.GetUpdatesChan(u)
	log.Println("Aguardando mensagens com polling aprimorado...")

	for {
		select {
		case update := <-updates:
			b.ProcessUpdate(update)
		case <-time.After(30 * time.Second):
			// Timeout, apenas para garantir que o loop não bloqueie para sempre
		}
	}
}
