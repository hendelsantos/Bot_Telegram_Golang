package modules

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

type AtualizaState struct {
	Step   int
	ItemID uint
	Campo  string
}

var atualizaStates = make(map[int64]*AtualizaState)

func HandleAtualizar(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	state, exists := atualizaStates[userID]
	if !exists {
		args := strings.TrimSpace(update.Message.CommandArguments())
		id, err := strconv.Atoi(args)
		if err != nil || id <= 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use: /atualizar <ID>")
			bot.Send(msg)
			return
		}
		atualizaStates[userID] = &AtualizaState{Step: 0, ItemID: uint(id)}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "O que deseja alterar? (nome, descricao, quantidade, foto)")
		bot.Send(msg)
		return
	}

	switch state.Step {
	case 0:
		campo := strings.ToLower(update.Message.Text)
		if campo == "nome" || campo == "descricao" || campo == "quantidade" || campo == "foto" {
			state.Campo = campo
			state.Step++
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Novo valor para %s:", campo))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Opção inválida. Escolha: nome, descricao, quantidade, foto.")
			bot.Send(msg)
		}
	case 1:
		var item db.Item
		db.DB.First(&item, state.ItemID)
		if item.ID == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item não encontrado.")
			bot.Send(msg)
			delete(atualizaStates, userID)
			return
		}
		switch state.Campo {
		case "nome":
			item.Nome = update.Message.Text
		case "descricao":
			item.Descricao = update.Message.Text
		case "quantidade":
			qtd, err := strconv.Atoi(update.Message.Text)
			if err != nil || qtd < 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Quantidade inválida.")
				bot.Send(msg)
				return
			}
			item.Quantidade = qtd
		case "foto":
			if len(update.Message.Photo) > 0 {
				photos := update.Message.Photo
				fileID := photos[len(photos)-1].FileID
				photoPath := fmt.Sprintf("photos/items/%s.jpg", fileID)
				item.FotoPath = photoPath
				log.Printf("Foto atualizada com ID: %s", fileID)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Envie uma foto válida.")
				bot.Send(msg)
				return
			}
		}
			db.DB.Save(&item)
			db.RegistrarHistorico(item.ID, "Atualização", "Campo atualizado: "+state.Campo)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item atualizado com sucesso!")
			bot.Send(msg)
			delete(atualizaStates, userID)
	}
}
