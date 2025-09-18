package modules

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleAjudaListagem fornece ajuda detalhada sobre os comandos de listagem
func HandleAjudaListagem(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	mensagem := "📖 **Guia Completo de Listagem**\n\n" +
		"**Comandos Básicos:**\n" +
		"`/listar` - Lista todos os itens com paginação\n" +
		"`/listar_resumo` - Lista apenas nome e quantidade\n" +
		"`/listar_detalhado` - Lista com informações completas\n\n" +
		"**Comandos Especializados:**\n" +
		"`/listar_status <nome>` - Lista por status específico\n" +
		"`/listar_baixo_estoque [limite]` - Alerta de estoque baixo\n\n" +
		"**Parâmetros Avançados:**\n" +
		"`/listar pagina=2` - Vai para página específica\n" +
		"`/listar limite=20` - Define itens por página (máx: 50)\n" +
		"`/listar ordenar_nome` - Ordena por nome A-Z\n" +
		"`/listar ordenar_qtd` - Ordena por quantidade\n" +
		"`/listar ordenar_data` - Ordena por data de criação\n\n" +
		"**Exemplos Práticos:**\n" +
		"• `/listar pagina=3 limite=5` - 5 itens na página 3\n" +
		"• `/listar_status Disponível` - Itens disponíveis\n" +
		"• `/listar_baixo_estoque 10` - Itens com ≤ 10 unidades\n" +
		"• `/listar_detalhado 2` - Detalhes específicos do item ID 2\n\n" +
		"**Dicas:**\n" +
		"🔹 Use `/listar_resumo` para visão geral rápida\n" +
		"🔹 Use `/listar_baixo_estoque` para controle de reposição\n" +
		"🔹 Use `/listar_detalhado` para informações completas\n" +
		"🔹 Navegue pelas páginas com os comandos sugeridos"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, mensagem)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
