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

func HandleMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	mensagem := "📋 *Menu de Comandos*\n\n" +
		"*Cadastro e Gestão*:\n" +
		"/novoitem - Cadastrar novo item no estoque\n" +
		"/buscar <palavra-chave> - Buscar itens por nome/descrição\n" +
		"/atualizar <ID> - Atualizar informações de um item\n\n" +
		"*Listagem e Visualização*:\n" +
		"/listar - Lista todos os itens (paginado)\n" +
		"/listar\\_resumo - Lista resumida (nome + quantidade)\n" +
		"/listar\\_detalhado - Lista com informações completas\n" +
		"/listar\\_status <status> - Lista por status específico\n" +
		"/listar\\_baixo\\_estoque [limite] - Alerta de estoque baixo\n\n" +
		"*Filtros Avançados*:\n" +
		"/buscar status <status> - Buscar por status\n" +
		"/buscar fornecedor <fornecedor> - Buscar por fornecedor\n" +
		"/buscar data <DD/MM/AAAA> - Buscar por data\n\n" +
		"*Controle de Reparos*:\n" +
		"/enviar\\_reparo <ID> - Registrar envio para reparo\n" +
		"/retornar\\_reparo <ID> - Registrar retorno de reparo\n\n" +
		"*Relatórios e Histórico*:\n" +
		"/exportar\\_estoque - Exportar lista em CSV\n" +
		"/historico <ID> - Ver histórico de movimentações"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, mensagem)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
