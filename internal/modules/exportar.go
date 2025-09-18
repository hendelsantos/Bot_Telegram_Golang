package modules

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

func HandleExportarEstoque(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	filename := "estoque_export_" + time.Now().Format("20060102_150405") + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro ao criar arquivo de exportação.")
		bot.Send(msg)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Nome", "Descrição", "Quantidade", "Status", "Fornecedor", "DataEnvio"})

	var itens []db.Item
	db.DB.Find(&itens)
	for _, item := range itens {
		writer.Write([]string{
			strconv.Itoa(int(item.ID)),
			item.Nome,
			item.Descricao,
			strconv.Itoa(item.Quantidade),
			item.Status,
			item.Fornecedor,
			item.DataEnvio,
		})
	}

	sentFile := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(filename))
	bot.Send(sentFile)
	os.Remove(filename)
}
