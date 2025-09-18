package modules

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

type CadastroState struct {
	Step      int
	FotoPath  string
	Nome      string
	Descricao string
	Qtd       int
}

var userStates = make(map[int64]*CadastroState)

func HandleNovoItem(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	
	// Se é um comando /novoitem, inicia o fluxo (não solicita foto imediatamente)
	if update.Message.IsCommand() && update.Message.Command() == "novoitem" {
		// Limpa qualquer estado anterior
		delete(userStates, userID)
		userStates[userID] = &CadastroState{Step: 0}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Envie uma foto do item para cadastro.")
		bot.Send(msg)
		return
	}
	
	// Se não existe estado, não está no fluxo de cadastro
	state, exists := userStates[userID]
	if !exists {
		return
	}

	switch state.Step {
	case 0:
		if len(update.Message.Photo) > 0 {
			// Pegar a foto com melhor qualidade (última do array)
			photos := update.Message.Photo
			fileID := photos[len(photos)-1].FileID
			
			// Salvar o FileID para referência futura
			photoPath := fmt.Sprintf("photos/items/%s.jpg", fileID)
			state.FotoPath = photoPath
			
			state.Step++
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Foto recebida! Agora envie o nome do item.")
			bot.Send(msg)
			log.Printf("Foto recebida e processada com ID: %s", fileID)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Por favor, envie uma foto válida.")
			bot.Send(msg)
		}
	case 1:
		state.Nome = update.Message.Text
		state.Step++
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Descrição/Observações?")
		bot.Send(msg)
	case 2:
		state.Descricao = update.Message.Text
		state.Step++
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Quantidade inicial?")
		bot.Send(msg)
	case 3:
		var qtd int
		_, err := fmt.Sscanf(update.Message.Text, "%d", &qtd)
		if err != nil || qtd < 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Quantidade inválida. Envie um número inteiro.")
			bot.Send(msg)
			return
		}
		state.Qtd = qtd
		state.Step++
		resumo := fmt.Sprintf("Resumo:\nNome: %s\nDescrição: %s\nQtd: %d\nConfirma? (sim/não)", state.Nome, state.Descricao, state.Qtd)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resumo)
		bot.Send(msg)
	case 4:
		if strings.ToLower(update.Message.Text) == "sim" {
			item := db.Item{
				Nome:      state.Nome,
				Descricao: state.Descricao,
				Quantidade: state.Qtd,
				Status:    "Em Estoque",
				FotoPath:  state.FotoPath,
			}
			err := db.DB.Create(&item).Error
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro ao salvar item.")
				bot.Send(msg)
				log.Println("Erro ao salvar item:", err)
			} else {
				db.RegistrarHistorico(item.ID, "Cadastro", "Item cadastrado no estoque")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Item cadastrado com sucesso!")
				bot.Send(msg)
			}
			delete(userStates, userID)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Cadastro cancelado.")
			bot.Send(msg)
			delete(userStates, userID)
		}
	}
}
