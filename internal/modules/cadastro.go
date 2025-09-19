package modules

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"botgo/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserStateCadastro struct {
	Nome       string
	Descricao  string
	Quantidade string
	FotoPath   string
	Step       string
}

var userStateCadastro = make(map[int64]*UserStateCadastro)

func HandleNovoItem(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userStateCadastro[chatID] = &UserStateCadastro{Step: "nome"}
	msg := tgbotapi.NewMessage(chatID, "📋 **Cadastro de Novo Item**\n\nQual o nome do item?")
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func ProcessCadastroFlow(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	state, exists := userStateCadastro[chatID]
	if !exists {
		return
	}

	switch state.Step {
	case "nome":
		state.Nome = message.Text
		state.Step = "descricao"
		msg := tgbotapi.NewMessage(chatID, "Qual a descrição do item?")
		bot.Send(msg)
	case "descricao":
		state.Descricao = message.Text
		state.Step = "quantidade"
		msg := tgbotapi.NewMessage(chatID, "Qual a quantidade inicial?")
		bot.Send(msg)
	case "quantidade":
		state.Quantidade = message.Text
		state.Step = "foto"
		msg := tgbotapi.NewMessage(chatID, "Envie uma foto do item (ou digite 'pular').")
		bot.Send(msg)
	case "foto":
		if message.Photo != nil && len(message.Photo) > 0 {
			// Lógica para salvar a foto...
			fileID := message.Photo[len(message.Photo)-1].FileID
			state.FotoPath = fileID // Apenas um exemplo, idealmente faria o download
			msg := tgbotapi.NewMessage(chatID, "Foto recebida!")
			bot.Send(msg)
		} else if message.Text != "pular" {
			msg := tgbotapi.NewMessage(chatID, "Por favor, envie uma foto ou digite 'pular'.")
			bot.Send(msg)
			return // Permanece no mesmo passo
		}

		// Finalizar cadastro
		quantidade, err := strconv.Atoi(state.Quantidade)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "❌ Quantidade inválida. O cadastro foi cancelado.")
			bot.Send(msg)
			delete(userStateCadastro, chatID)
			return
		}

		item := db.Item{
			Nome:       state.Nome,
			Descricao:  state.Descricao,
			Quantidade: quantidade,
			Status:     "Em Estoque",
			FotoPath:   state.FotoPath,
		}

		result := db.DB.Create(&item)
		if result.Error != nil {
			log.Printf("Erro ao salvar item: %v", result.Error)
			msg := tgbotapi.NewMessage(chatID, "❌ Ocorreu um erro ao salvar o item.")
			bot.Send(msg)
		} else {
			// Registrar movimentação
			movimentacao := db.Movimentacao{
				ItemID:    item.ID,
				Tipo:      "cadastro",
				Descricao: fmt.Sprintf("Item '%s' cadastrado com quantidade inicial %d.", item.Nome, item.Quantidade),
				DataHora:  time.Now(),
			}
			db.DB.Create(&movimentacao)

			msgText := fmt.Sprintf("✅ **Item Cadastrado com Sucesso!**\n\n**ID:** %d\n**Nome:** %s\n**Descrição:** %s\n**Quantidade:** %d", item.ID, item.Nome, item.Descricao, item.Quantidade)
			msg := tgbotapi.NewMessage(chatID, msgText)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
		delete(userStateCadastro, chatID)
	}
}
