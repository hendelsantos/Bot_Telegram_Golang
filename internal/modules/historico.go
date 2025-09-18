package modules

import (
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

func HandleHistorico(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use: /historico <ID>")
		bot.Send(msg)
		return
	}
	var historicos []db.Historico
	db.DB.Where("item_id = ?", args).Order("data_hora desc").Find(&historicos)
	if len(historicos) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nenhum histórico encontrado para este item.")
		bot.Send(msg)
		return
	}
	for _, h := range historicos {
		texto := fmt.Sprintf("%s - %s: %s", h.DataHora.Format("02/01/2006 15:04"), h.Acao, h.Descricao)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texto)
		bot.Send(msg)
	}
}
