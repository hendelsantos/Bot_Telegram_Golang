package modules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"
)

// ListarConfig define configurações para listagem
type ListarConfig struct {
	Limite     int
	Pagina     int
	Ordenacao  string
	Filtros    map[string]interface{}
}

// ResultadoListagem representa o resultado de uma consulta
type ResultadoListagem struct {
	Itens      []db.Item
	Total      int64
	Pagina     int
	TotalPags  int
	HasProx    bool
	HasAnt     bool
}

// HandleListar gerencia todos os comandos de listagem
func HandleListar(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	args := strings.Fields(update.Message.Text)
	comando := args[0]

	switch comando {
	case "/listar":
		handleListarTodos(bot, update, args[1:])
	case "/listar_resumo":
		handleListarResumo(bot, update, args[1:])
	case "/listar_status":
		handleListarPorStatus(bot, update, args[1:])
	case "/listar_baixo_estoque":
		handleListarBaixoEstoque(bot, update, args[1:])
	case "/listar_detalhado":
		handleListarDetalhado(bot, update, args[1:])
	}
}

// handleListarTodos lista todos os itens com paginação
func handleListarTodos(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	
	// Configuração padrão
	config := ListarConfig{
		Limite:    10,
		Pagina:    1,
		Ordenacao: "id DESC",
		Filtros:   make(map[string]interface{}),
	}
	
	// Parse dos argumentos
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "pagina="):
			if p, err := strconv.Atoi(strings.TrimPrefix(arg, "pagina=")); err == nil && p > 0 {
				config.Pagina = p
			}
		case strings.HasPrefix(arg, "limite="):
			if l, err := strconv.Atoi(strings.TrimPrefix(arg, "limite=")); err == nil && l > 0 && l <= 50 {
				config.Limite = l
			}
		case arg == "ordenar_nome":
			config.Ordenacao = "nome ASC"
		case arg == "ordenar_qtd":
			config.Ordenacao = "quantidade ASC"
		case arg == "ordenar_data":
			config.Ordenacao = "created_at DESC"
		}
	}
	
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		return
	}
	
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Nenhum item cadastrado no estoque.")
		bot.Send(msg)
		return
	}
	
	enviarListagemFormatada(bot, chatID, resultado, "📋 **Lista Completa do Estoque**")
}

// handleListarResumo lista apenas informações básicas
func handleListarResumo(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	
	config := ListarConfig{
		Limite:    20,
		Pagina:    1,
		Ordenacao: "nome ASC",
		Filtros:   make(map[string]interface{}),
	}
	
	// Parse da página se fornecida
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil && p > 0 {
			config.Pagina = p
		}
	}
	
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		return
	}
	
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Estoque vazio.")
		bot.Send(msg)
		return
	}
	
	enviarResumoFormatado(bot, chatID, resultado)
}

// handleListarPorStatus lista itens por status específico
func handleListarPorStatus(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	
	if len(args) == 0 {
		enviarStatusDisponiveis(bot, chatID)
		return
	}
	
	status := strings.Join(args, " ")
	
	config := ListarConfig{
		Limite:    15,
		Pagina:    1,
		Ordenacao: "nome ASC",
		Filtros:   map[string]interface{}{"status": status},
	}
	
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		return
	}
	
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📭 Nenhum item encontrado com status: **%s**", status))
		msg.ParseMode = "Markdown"
		bot.Send(msg)
		return
	}
	
	titulo := fmt.Sprintf("📊 **Itens com Status: %s**", status)
	enviarListagemFormatada(bot, chatID, resultado, titulo)
}

// handleListarBaixoEstoque lista itens com quantidade baixa
func handleListarBaixoEstoque(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	
	limite := 5 // Quantidade mínima padrão
	if len(args) > 0 {
		if l, err := strconv.Atoi(args[0]); err == nil && l >= 0 {
			limite = l
		}
	}
	
	config := ListarConfig{
		Limite:    50,
		Pagina:    1,
		Ordenacao: "quantidade ASC",
		Filtros:   map[string]interface{}{"quantidade_max": limite},
	}
	
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		return
	}
	
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Não há itens com estoque baixo (≤ %d)", limite))
		bot.Send(msg)
		return
	}
	
	titulo := fmt.Sprintf("⚠️ **Alerta: Estoque Baixo (≤ %d)**", limite)
	enviarListagemFormatada(bot, chatID, resultado, titulo)
}

// handleListarDetalhado lista itens com todas as informações
func handleListarDetalhado(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	
	// Suporte a ID específico
	if len(args) > 0 {
		if id, err := strconv.Atoi(args[0]); err == nil {
			enviarItemDetalhado(bot, chatID, uint(id))
			return
		}
	}
	
	config := ListarConfig{
		Limite:    5,
		Pagina:    1,
		Ordenacao: "id DESC",
		Filtros:   make(map[string]interface{}),
	}
	
	// Parse da página
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil && p > 0 {
			config.Pagina = p
		}
	}
	
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		return
	}
	
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Nenhum item cadastrado.")
		bot.Send(msg)
		return
	}
	
	enviarListagemDetalhada(bot, chatID, resultado)
}

// buscarItensComPaginacao executa a consulta com base na configuração
func buscarItensComPaginacao(config ListarConfig) (*ResultadoListagem, error) {
	query := db.DB.Model(&db.Item{})
	
	// Aplicar filtros
	for campo, valor := range config.Filtros {
		switch campo {
		case "status":
			query = query.Where("status LIKE ?", "%"+valor.(string)+"%")
		case "quantidade_max":
			query = query.Where("quantidade <= ?", valor.(int))
		case "nome":
			query = query.Where("nome LIKE ?", "%"+valor.(string)+"%")
		case "fornecedor":
			query = query.Where("fornecedor LIKE ?", "%"+valor.(string)+"%")
		}
	}
	
	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// Calcular offset
	offset := (config.Pagina - 1) * config.Limite
	
	// Buscar itens
	var itens []db.Item
	if err := query.Offset(offset).Limit(config.Limite).Order(config.Ordenacao).Find(&itens).Error; err != nil {
		return nil, err
	}
	
	// Calcular informações de paginação
	totalPags := int((total + int64(config.Limite) - 1) / int64(config.Limite))
	
	return &ResultadoListagem{
		Itens:     itens,
		Total:     total,
		Pagina:    config.Pagina,
		TotalPags: totalPags,
		HasProx:   config.Pagina < totalPags,
		HasAnt:    config.Pagina > 1,
	}, nil
}

// enviarListagemFormatada envia lista formatada com paginação
func enviarListagemFormatada(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem, titulo string) {
	var mensagem strings.Builder
	
	mensagem.WriteString(titulo + "\n\n")
	
	for i, item := range resultado.Itens {
		status := item.Status
		if status == "" {
			status = "Disponível"
		}
		
		mensagem.WriteString(fmt.Sprintf("**%d.** `%s` (ID: %d)\n", 
			i+1, item.Nome, item.ID))
		mensagem.WriteString(fmt.Sprintf("   📦 Qtd: **%d**  |  📋 Status: **%s**\n", 
			item.Quantidade, status))
		
		if item.Descricao != "" {
			desc := item.Descricao
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			mensagem.WriteString(fmt.Sprintf("   💬 %s\n", desc))
		}
		mensagem.WriteString("\n")
	}
	
	// Informações de paginação
	mensagem.WriteString(fmt.Sprintf("📄 **Página %d de %d** | **Total: %d itens**\n", 
		resultado.Pagina, resultado.TotalPags, resultado.Total))
	
	if resultado.HasProx || resultado.HasAnt {
		mensagem.WriteString("\n💡 *Use:*\n")
		if resultado.HasAnt {
			mensagem.WriteString(fmt.Sprintf("• `/listar pagina=%d` - Página anterior\n", resultado.Pagina-1))
		}
		if resultado.HasProx {
			mensagem.WriteString(fmt.Sprintf("• `/listar pagina=%d` - Próxima página\n", resultado.Pagina+1))
		}
	}
	
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// enviarResumoFormatado envia apenas nome e quantidade
func enviarResumoFormatado(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	var mensagem strings.Builder
	
	mensagem.WriteString("📋 **Resumo do Estoque**\n\n")
	
	for i, item := range resultado.Itens {
		mensagem.WriteString(fmt.Sprintf("**%d.** %s - **%d unid.**\n", 
			i+1, item.Nome, item.Quantidade))
	}
	
	mensagem.WriteString(fmt.Sprintf("\n📊 **Total: %d itens cadastrados**", resultado.Total))
	
	if resultado.HasProx {
		mensagem.WriteString(fmt.Sprintf("\n\n💡 *Use `/listar_resumo %d` para próxima página*", resultado.Pagina+1))
	}
	
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// enviarListagemDetalhada envia informações completas
func enviarListagemDetalhada(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	for _, item := range resultado.Itens {
		enviarItemDetalhado(bot, chatID, item.ID)
	}
	
	if len(resultado.Itens) > 1 {
		info := fmt.Sprintf("📄 **Página %d de %d** | **Total: %d itens**", 
			resultado.Pagina, resultado.TotalPags, resultado.Total)
		
		if resultado.HasProx {
			info += fmt.Sprintf("\n💡 *Use `/listar_detalhado %d` para próxima página*", resultado.Pagina+1)
		}
		
		msg := tgbotapi.NewMessage(chatID, info)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}

// enviarItemDetalhado envia informações completas de um item específico
func enviarItemDetalhado(bot *tgbotapi.BotAPI, chatID int64, itemID uint) {
	var item db.Item
	if err := db.DB.First(&item, itemID).Error; err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Item com ID %d não encontrado.", itemID))
		bot.Send(msg)
		return
	}
	
	var mensagem strings.Builder
	
	mensagem.WriteString(fmt.Sprintf("🔍 **Detalhes do Item #%d**\n\n", item.ID))
	mensagem.WriteString(fmt.Sprintf("📝 **Nome:** %s\n", item.Nome))
	mensagem.WriteString(fmt.Sprintf("📦 **Quantidade:** %d unidades\n", item.Quantidade))
	
	status := item.Status
	if status == "" {
		status = "Disponível"
	}
	mensagem.WriteString(fmt.Sprintf("📋 **Status:** %s\n", status))
	
	if item.Descricao != "" {
		mensagem.WriteString(fmt.Sprintf("💬 **Descrição:** %s\n", item.Descricao))
	}
	
	if item.Fornecedor != "" {
		mensagem.WriteString(fmt.Sprintf("🏢 **Fornecedor:** %s\n", item.Fornecedor))
	}
	
	if item.DataEnvio != nil {
		mensagem.WriteString(fmt.Sprintf("📅 **Data Envio:** %s\n", item.DataEnvio))
	}
	
	if item.FotoPath != "" {
		mensagem.WriteString("📸 **Foto:** Disponível\n")
	}
	
	mensagem.WriteString(fmt.Sprintf("\n⏰ **Criado em:** %s", time.Now().Format("02/01/2006 15:04")))
	
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// enviarStatusDisponiveis lista os status disponíveis
func enviarStatusDisponiveis(bot *tgbotapi.BotAPI, chatID int64) {
	var statuses []string
	db.DB.Model(&db.Item{}).Distinct("status").Pluck("status", &statuses)
	
	var mensagem strings.Builder
	mensagem.WriteString("📊 **Status Disponíveis:**\n\n")
	
	if len(statuses) == 0 {
		mensagem.WriteString("Nenhum status encontrado.\n")
	} else {
		for i, status := range statuses {
			if status == "" {
				status = "Disponível"
			}
			mensagem.WriteString(fmt.Sprintf("**%d.** %s\n", i+1, status))
		}
	}
	
	mensagem.WriteString("\n💡 *Use:* `/listar_status <nome_do_status>`")
	
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// enviarMensagemErro envia mensagem de erro formatada
func enviarMensagemErro(bot *tgbotapi.BotAPI, chatID int64, erro string) {
	msg := tgbotapi.NewMessage(chatID, "❌ **Erro:** "+erro)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
