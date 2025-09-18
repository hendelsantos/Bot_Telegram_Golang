package modules

import (
	"strings"
	"time"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

type ReparoState struct {
	Step      int
	ItemID    uint
	Fornecedor string
	DataEnvio string
}

var reparoStates = make(map[int64]*ReparoState)

func HandleEnviarReparo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	state, exists := reparoStates[userID]
	if !exists {
		args := strings.TrimSpace(update.Message.CommandArguments())
		id, err := strconv.Atoi(args)
		if err != nil || id <= 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use: /enviar_reparo <ID>")
			bot.Send(msg)
			return
		}
		reparoStates[userID] = &ReparoState{Step: 0, ItemID: uint(id)}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Para qual fornecedor/local?")
		bot.Send(msg)
		return
	}

	switch state.Step {
	case 0:
		state.Fornecedor = update.Message.Text
		state.Step++
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Data de envio? (formato: DD/MM/AAAA)")
		bot.Send(msg)
	case 1:
		// Validação simples de data
		_, err := time.Parse("02/01/2006", update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Data inválida. Use o formato DD/MM/AAAA.")
			bot.Send(msg)
			return
		}
		state.DataEnvio = update.Message.Text
		var item db.Item
		db.DB.First(&item, state.ItemID)
		if item.ID == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item não encontrado.")
			bot.Send(msg)
			delete(reparoStates, userID)
			return
		}
		if item.Quantidade <= 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sem estoque disponível para envio.")
			bot.Send(msg)
			delete(reparoStates, userID)
			return
		}
		item.Status = "Em Reparo Externo"
		item.Fornecedor = state.Fornecedor
		item.DataEnvio = state.DataEnvio
		item.Quantidade--
		db.DB.Save(&item)
		db.RegistrarHistorico(item.ID, "Envio para Reparo", "Enviado para: "+state.Fornecedor+", Data: "+state.DataEnvio)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item enviado para reparo externo!")
		bot.Send(msg)
		delete(reparoStates, userID)
	}
}

func HandleRetornarReparo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	id, err := strconv.Atoi(args)
	if err != nil || id <= 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use: /retornar_reparo <ID>")
		bot.Send(msg)
		return
	}
	var item db.Item
	db.DB.First(&item, uint(id))
	if item.ID == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item não encontrado.")
		bot.Send(msg)
		return
	}
	item.Status = "Em Estoque"
	item.Quantidade++
	db.DB.Save(&item)
	db.RegistrarHistorico(item.ID, "Retorno de Reparo", "Item retornou do reparo externo")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item retornou do reparo e está em estoque!")
	bot.Send(msg)
}
