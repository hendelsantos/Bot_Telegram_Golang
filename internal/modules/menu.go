package modules

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	mensagem := "👋 *Bem-vindo ao Assistente de Estoque*!\n\n" +
		"Este bot permite que você gerencie seu estoque de forma simples e eficiente.\n\n" +
		"Use /menu para ver todos os comandos disponíveis."

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, mensagem)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// HandleMenu agora aceita chatID diretamente.
func HandleMenu(bot *tgbotapi.BotAPI, chatID int64) {
	text := `
📋 *Menu de Comandos*

*Cadastro e Gestão*:
/novoitem - Cadastrar novo item
/buscar <palavra> - Buscar itens
/atualizar <ID> - Atualizar um item

*Listagem e Visualização*:
/listar - Listar todos os itens
/listar_status <status> - Listar por status
/listar_baixo_estoque - Itens com estoque baixo

*Controle de Reparos*:
/enviar_reparo <ID> - Registrar envio para reparo
/retornar_reparo <ID> - Registrar retorno de reparo

*Relatórios e Histórico*:
/exportar_estoque - Exportar estoque em CSV
/historico <ID> - Ver histórico de um item
`
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
