package modules

import (
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

func HandleBuscar(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use: /buscar <palavra-chave>\n/buscar_status <status>\n/buscar_fornecedor <fornecedor>\n/buscar_data <DD/MM/AAAA>")
		bot.Send(msg)
		return
	}

	var itens []db.Item
	if strings.HasPrefix(args, "status ") {
		status := strings.TrimSpace(strings.TrimPrefix(args, "status "))
		db.DB.Where("status = ?", status).Find(&itens)
	} else if strings.HasPrefix(args, "fornecedor ") {
		fornecedor := strings.TrimSpace(strings.TrimPrefix(args, "fornecedor "))
		db.DB.Where("fornecedor LIKE ?", "%"+fornecedor+"%").Find(&itens)
	} else if strings.HasPrefix(args, "data ") {
		data := strings.TrimSpace(strings.TrimPrefix(args, "data "))
		// Busca por data de envio
		db.DB.Where("data_envio = ?", data).Find(&itens)
	} else {
		db.DB.Where("nome LIKE ? OR descricao LIKE ?", "%"+args+"%", "%"+args+"%").Find(&itens)
	}

	if len(itens) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nenhum item encontrado.")
		bot.Send(msg)
		return
	}

	for _, item := range itens {
		texto := fmt.Sprintf("ID: %d\nNome: %s\nDescrição: %s\nQtd: %d\nStatus: %s", item.ID, item.Nome, item.Descricao, item.Quantidade, item.Status)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texto)
		if item.FotoPath != "" {
			photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(item.FotoPath))
			photo.Caption = texto
			bot.Send(photo)
		} else {
			bot.Send(msg)
		}
	}
}
