package modules
package modules
import (
	"fmt" (
	"strconv"
	"strings"
	"time"gs"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botgo/internal/db"gram-bot-api/telegram-bot-api/v5"
)"botgo/internal/db"
)
// ListarConfig define configurações para listagem
type ListarConfig struct {figurações para listagem
	Limite     intig struct {
	Pagina     int
	Ordenacao  string
	Filtros    map[string]interface{}
}Filtros    map[string]interface{}
}
// ResultadoListagem representa o resultado de uma consulta
type ResultadoListagem struct { o resultado de uma consulta
	Itens      []db.Itemm struct {
	Total      int64Item
	Pagina     int64
	TotalPags  int
	HasProx    bool
	HasAnt     bool
}HasAnt     bool
}
// HandleListar gerencia todos os comandos de listagem
func HandleListar(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	args := strings.Fields(update.Message.Text)te tgbotapi.Update) {
	comando := args[0]elds(update.Message.Text)
	comando := args[0]
	switch comando {
	case "/listar":{
		handleListarTodos(bot, update, args[1:])
	case "/listar_resumo":, update, args[1:])
		handleListarResumo(bot, update, args[1:])
	case "/listar_status":t, update, args[1:])
		handleListarPorStatus(bot, update, args[1:])
	case "/listar_baixo_estoque":pdate, args[1:])
		handleListarBaixoEstoque(bot, update, args[1:])
	case "/listar_detalhado":(bot, update, args[1:])
		handleListarDetalhado(bot, update, args[1:])
	}handleListarDetalhado(bot, update, args[1:])
}}
}
// handleListarTodos lista todos os itens com paginação
func handleListarTodos(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.IDapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	// Configuração padrão
	config := ListarConfig{
		Limite:    10,rConfig{
		Pagina:    1,,
		Ordenacao: "id DESC",
		Filtros:   make(map[string]interface{}),
	}Filtros:   make(map[string]interface{}),
	}
	// Parse dos argumentos
	for _, arg := range args {
		switch {g := range args {
		case strings.HasPrefix(arg, "pagina="):
			if p, err := strconv.Atoi(strings.TrimPrefix(arg, "pagina=")); err == nil && p > 0 {
				config.Pagina = pnv.Atoi(strings.TrimPrefix(arg, "pagina=")); err == nil && p > 0 {
			}config.Pagina = p
		case strings.HasPrefix(arg, "limite="):
			if l, err := strconv.Atoi(strings.TrimPrefix(arg, "limite=")); err == nil && l > 0 && l <= 50 {
				config.Limite = lnv.Atoi(strings.TrimPrefix(arg, "limite=")); err == nil && l > 0 && l <= 50 {
			}config.Limite = l
		case arg == "ordenar_nome":
			config.Ordenacao = "nome ASC"
		case arg == "ordenar_qtd":ASC"
			config.Ordenacao = "quantidade ASC"
		case arg == "ordenar_data":dade ASC"
			config.Ordenacao = "created_at DESC"
		}config.Ordenacao = "created_at DESC"
	}}
	}
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {:= buscarItensComPaginacao(config)
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		returnMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
	}return
	}
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Nenhum item cadastrado no estoque.")
		bot.Send(msg)pi.NewMessage(chatID, "📭 Nenhum item cadastrado no estoque.")
		returnnd(msg)
	}return
	}
	enviarListagemFormatada(bot, chatID, resultado, "📋 **Lista Completa do Estoque**")
}enviarListagemFormatada(bot, chatID, resultado, "📋 **Lista Completa do Estoque**")
}
// handleListarResumo lista apenas informações básicas
func handleListarResumo(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.IDtapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	config := ListarConfig{
		Limite:    20,rConfig{
		Pagina:    1,,
		Ordenacao: "nome ASC",
		Filtros:   make(map[string]interface{}),
	}Filtros:   make(map[string]interface{}),
	}
	// Parse da página se fornecida
	if len(args) > 0 { se fornecida
		if p, err := strconv.Atoi(args[0]); err == nil && p > 0 {
			config.Pagina = pnv.Atoi(args[0]); err == nil && p > 0 {
		}config.Pagina = p
	}}
	}
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {:= buscarItensComPaginacao(config)
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		returnMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
	}return
	}
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Estoque vazio.")
		bot.Send(msg)pi.NewMessage(chatID, "📭 Estoque vazio.")
		returnnd(msg)
	}return
	}
	enviarResumoFormatado(bot, chatID, resultado)
}enviarResumoFormatado(bot, chatID, resultado)
}
// handleListarPorStatus lista itens por status específico
func handleListarPorStatus(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.IDgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	if len(args) == 0 {
		enviarStatusDisponiveis(bot, chatID)
		returnStatusDisponiveis(bot, chatID)
	}return
	}
	status := strings.Join(args, " ")
	status := strings.Join(args, " ")
	config := ListarConfig{
		Limite:    15,rConfig{
		Pagina:    1,,
		Ordenacao: "nome ASC",
		Filtros:   map[string]interface{}{"status": status},
	}Filtros:   map[string]interface{}{"status": status},
	}
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {:= buscarItensComPaginacao(config)
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		returnMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
	}return
	}
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📭 Nenhum item encontrado com status: **%s**", status))
		msg.ParseMode = "Markdown"(chatID, fmt.Sprintf("📭 Nenhum item encontrado com status: **%s**", status))
		bot.Send(msg) = "Markdown"
		returnnd(msg)
	}return
	}
	titulo := fmt.Sprintf("📊 **Itens com Status: %s**", status)
	enviarListagemFormatada(bot, chatID, resultado, titulo)atus)
}enviarListagemFormatada(bot, chatID, resultado, titulo)
}
// handleListarBaixoEstoque lista itens com quantidade baixa
func handleListarBaixoEstoque(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	limite := 5 // Quantidade mínima padrão
	if len(args) > 0 {ntidade mínima padrão
		if l, err := strconv.Atoi(args[0]); err == nil && l >= 0 {
			limite = l= strconv.Atoi(args[0]); err == nil && l >= 0 {
		}limite = l
	}}
	}
	config := ListarConfig{
		Limite:    50,rConfig{
		Pagina:    1,,
		Ordenacao: "quantidade ASC",
		Filtros:   map[string]interface{}{"quantidade_max": limite},
	}Filtros:   map[string]interface{}{"quantidade_max": limite},
	}
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {:= buscarItensComPaginacao(config)
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		returnMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
	}return
	}
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Não há itens com estoque baixo (≤ %d)", limite))
		bot.Send(msg)pi.NewMessage(chatID, fmt.Sprintf("✅ Não há itens com estoque baixo (≤ %d)", limite))
		returnnd(msg)
	}return
	}
	titulo := fmt.Sprintf("⚠️ **Alerta: Estoque Baixo (≤ %d)**", limite)
	enviarListagemFormatada(bot, chatID, resultado, titulo))**", limite)
}enviarListagemFormatada(bot, chatID, resultado, titulo)
}
// handleListarDetalhado lista itens com todas as informações
func handleListarDetalhado(bot *tgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.IDgbotapi.BotAPI, update tgbotapi.Update, args []string) {
	chatID := update.Message.Chat.ID
	// Suporte a ID específico
	if len(args) > 0 {pecífico
		if id, err := strconv.Atoi(args[0]); err == nil {
			enviarItemDetalhado(bot, chatID, uint(id)) nil {
			returnItemDetalhado(bot, chatID, uint(id))
		}return
	}}
	}
	config := ListarConfig{
		Limite:    5,arConfig{
		Pagina:    1,
		Ordenacao: "id DESC",
		Filtros:   make(map[string]interface{}),
	}Filtros:   make(map[string]interface{}),
	}
	// Parse da página
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil && p > 0 {
			config.Pagina = pnv.Atoi(args[0]); err == nil && p > 0 {
		}config.Pagina = p
	}}
	}
	resultado, err := buscarItensComPaginacao(config)
	if err != nil {:= buscarItensComPaginacao(config)
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		returnMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
	}return
	}
	if resultado.Total == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Nenhum item cadastrado.")
		bot.Send(msg)pi.NewMessage(chatID, "📭 Nenhum item cadastrado.")
		returnnd(msg)
	}return
	}
	enviarListagemDetalhada(bot, chatID, resultado)
}enviarListagemDetalhada(bot, chatID, resultado)
}
// buscarItensComPaginacao executa a consulta com base na configuração
func buscarItensComPaginacao(config ListarConfig) (*ResultadoListagem, error) {
	query := db.DB.Model(&db.Item{})ig ListarConfig) (*ResultadoListagem, error) {
	query := db.DB.Model(&db.Item{})
	// Aplicar filtros
	for campo, valor := range config.Filtros {
		switch campo {r := range config.Filtros {
		case "status":
			query = query.Where("status LIKE ?", "%"+valor.(string)+"%")
		case "quantidade_max":status LIKE ?", "%"+valor.(string)+"%")
			query = query.Where("quantidade <= ?", valor.(int))
		case "nome":ry.Where("quantidade <= ?", valor.(int))
			query = query.Where("nome LIKE ?", "%"+valor.(string)+"%")
		case "fornecedor":re("nome LIKE ?", "%"+valor.(string)+"%")
			query = query.Where("fornecedor LIKE ?", "%"+valor.(string)+"%")
		}query = query.Where("fornecedor LIKE ?", "%"+valor.(string)+"%")
	}}
	}
	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errCount(&total).Error; err != nil {
	}return nil, err
	}
	// Calcular offset
	offset := (config.Pagina - 1) * config.Limite
	offset := (config.Pagina - 1) * config.Limite
	// Buscar itens
	var itens []db.Item
	if err := query.Offset(offset).Limit(config.Limite).Order(config.Ordenacao).Find(&itens).Error; err != nil {
		return nil, errOffset(offset).Limit(config.Limite).Order(config.Ordenacao).Find(&itens).Error; err != nil {
	}return nil, err
	}
	// Calcular informações de paginação
	totalPags := int((total + int64(config.Limite) - 1) / int64(config.Limite))
	totalPags := int((total + int64(config.Limite) - 1) / int64(config.Limite))
	return &ResultadoListagem{
		Itens:     itens,istagem{
		Total:     total,
		Pagina:    config.Pagina,
		TotalPags: totalPags,ina,
		HasProx:   config.Pagina < totalPags,
		HasAnt:    config.Pagina > 1,talPags,
	}, nilt:    config.Pagina > 1,
}}, nil
}
// enviarListagemFormatada envia lista formatada com paginação
func enviarListagemFormatada(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem, titulo string) {
	var mensagem strings.Builderbot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem, titulo string) {
	var mensagem strings.Builder
	mensagem.WriteString(titulo + "\n\n")
	mensagem.WriteString(titulo + "\n\n")
	for i, item := range resultado.Itens {
		status := item.Statusesultado.Itens {
		if status == "" {atus
			status = "Disponível"
		}status = "Disponível"
		}
		mensagem.WriteString(fmt.Sprintf("**%d.** `%s` (ID: %d)\n", 
			i+1, item.Nome, item.ID))printf("**%d.** `%s` (ID: %d)\n", 
		mensagem.WriteString(fmt.Sprintf("   📦 Qtd: **%d**  |  📋 Status: **%s**\n", 
			item.Quantidade, status))printf("   📦 Qtd: **%d**  |  📋 Status: **%s**\n", 
			item.Quantidade, status))
		if item.Descricao != "" {
			desc := item.Descricao {
			if len(desc) > 50 {cao
				desc = desc[:47] + "..."
			}desc = desc[:47] + "..."
			mensagem.WriteString(fmt.Sprintf("   💬 %s\n", desc))
		}mensagem.WriteString(fmt.Sprintf("   💬 %s\n", desc))
		mensagem.WriteString("\n")
	}mensagem.WriteString("\n")
	}
	// Informações de paginação
	mensagem.WriteString(fmt.Sprintf("📄 **Página %d de %d** | **Total: %d itens**\n", 
		resultado.Pagina, resultado.TotalPags, resultado.Total))| **Total: %d itens**\n", 
		resultado.Pagina, resultado.TotalPags, resultado.Total))
	if resultado.HasProx || resultado.HasAnt {
		mensagem.WriteString("\n💡 *Use:*\n")nt {
		if resultado.HasAnt {"\n💡 *Use:*\n")
			mensagem.WriteString(fmt.Sprintf("• `/listar pagina=%d` - Página anterior\n", resultado.Pagina-1))
		}mensagem.WriteString(fmt.Sprintf("• `/listar pagina=%d` - Página anterior\n", resultado.Pagina-1))
		if resultado.HasProx {
			mensagem.WriteString(fmt.Sprintf("• `/listar pagina=%d` - Próxima página\n", resultado.Pagina+1))
		}mensagem.WriteString(fmt.Sprintf("• `/listar pagina=%d` - Próxima página\n", resultado.Pagina+1))
	}}
	}
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"(chatID, mensagem.String())
	bot.Send(msg) = "Markdown"
}bot.Send(msg)
}
// enviarResumoFormatado envia apenas nome e quantidade
func enviarResumoFormatado(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	var mensagem strings.Buildert *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	var mensagem strings.Builder
	mensagem.WriteString("📋 **Resumo do Estoque**\n\n")
	mensagem.WriteString("📋 **Resumo do Estoque**\n\n")
	for i, item := range resultado.Itens {
		mensagem.WriteString(fmt.Sprintf("**%d.** %s - **%d unid.**\n", 
			i+1, item.Nome, item.Quantidade))**%d.** %s - **%d unid.**\n", 
	}	i+1, item.Nome, item.Quantidade))
	}
	mensagem.WriteString(fmt.Sprintf("\n📊 **Total: %d itens cadastrados**", resultado.Total))
	mensagem.WriteString(fmt.Sprintf("\n📊 **Total: %d itens cadastrados**", resultado.Total))
	if resultado.HasProx {
		mensagem.WriteString(fmt.Sprintf("\n\n💡 *Use `/listar_resumo %d` para próxima página*", resultado.Pagina+1))
	}mensagem.WriteString(fmt.Sprintf("\n\n💡 *Use `/listar_resumo %d` para próxima página*", resultado.Pagina+1))
	}
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"(chatID, mensagem.String())
	bot.Send(msg) = "Markdown"
}bot.Send(msg)
}
// enviarListagemDetalhada envia informações completas
func enviarListagemDetalhada(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	for _, item := range resultado.Itens {api.BotAPI, chatID int64, resultado *ResultadoListagem) {
		enviarItemDetalhado(bot, chatID, item.ID)
	}enviarItemDetalhado(bot, chatID, item.ID)
	}
	if len(resultado.Itens) > 1 {
		info := fmt.Sprintf("📄 **Página %d de %d** | **Total: %d itens**", 
			resultado.Pagina, resultado.TotalPags, resultado.Total)d itens**", 
			resultado.Pagina, resultado.TotalPags, resultado.Total)
		if resultado.HasProx {
			info += fmt.Sprintf("\n💡 *Use `/listar_detalhado %d` para próxima página*", resultado.Pagina+1)
		}info += fmt.Sprintf("\n💡 *Use `/listar_detalhado %d` para próxima página*", resultado.Pagina+1)
		}
		msg := tgbotapi.NewMessage(chatID, info)
		msg.ParseMode = "Markdown"(chatID, info)
		bot.Send(msg) = "Markdown"
	}bot.Send(msg)
}}
}
// enviarItemDetalhado envia informações completas de um item específico
func enviarItemDetalhado(bot *tgbotapi.BotAPI, chatID int64, itemID uint) {
	var item db.Itemtalhado(bot *tgbotapi.BotAPI, chatID int64, itemID uint) {
	if err := db.DB.First(&item, itemID).Error; err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Item com ID %d não encontrado.", itemID))
		bot.Send(msg)pi.NewMessage(chatID, fmt.Sprintf("❌ Item com ID %d não encontrado.", itemID))
		returnnd(msg)
	}return
	}
	var mensagem strings.Builder
	var mensagem strings.Builder
	mensagem.WriteString(fmt.Sprintf("🔍 **Detalhes do Item #%d**\n\n", item.ID))
	mensagem.WriteString(fmt.Sprintf("📝 **Nome:** %s\n", item.Nome))", item.ID))
	mensagem.WriteString(fmt.Sprintf("📦 **Quantidade:** %d unidades\n", item.Quantidade))
	mensagem.WriteString(fmt.Sprintf("📦 **Quantidade:** %d unidades\n", item.Quantidade))
	status := item.Status
	if status == "" {atus
		status = "Disponível"
	}status = "Disponível"
	mensagem.WriteString(fmt.Sprintf("📋 **Status:** %s\n", status))
	mensagem.WriteString(fmt.Sprintf("📋 **Status:** %s\n", status))
	if item.Descricao != "" {
		mensagem.WriteString(fmt.Sprintf("💬 **Descrição:** %s\n", item.Descricao))
	}mensagem.WriteString(fmt.Sprintf("💬 **Descrição:** %s\n", item.Descricao))
	}
	if item.Fornecedor != "" {
		mensagem.WriteString(fmt.Sprintf("🏢 **Fornecedor:** %s\n", item.Fornecedor))
	}mensagem.WriteString(fmt.Sprintf("🏢 **Fornecedor:** %s\n", item.Fornecedor))
	}
	if item.DataEnvio != nil {
		mensagem.WriteString(fmt.Sprintf("📅 **Data Envio:** %s\n", item.DataEnvio))
	}mensagem.WriteString(fmt.Sprintf("📅 **Data Envio:** %s\n", item.DataEnvio))
	}
	if item.FotoPath != "" {
		mensagem.WriteString("📸 **Foto:** Disponível\n")
	}mensagem.WriteString("📸 **Foto:** Disponível\n")
	}
	mensagem.WriteString(fmt.Sprintf("\n⏰ **Criado em:** %s", time.Now().Format("02/01/2006 15:04")))
	mensagem.WriteString(fmt.Sprintf("\n⏰ **Criado em:** %s", time.Now().Format("02/01/2006 15:04")))
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"(chatID, mensagem.String())
	bot.Send(msg) = "Markdown"
}bot.Send(msg)
}
// enviarStatusDisponiveis lista os status disponíveis
func enviarStatusDisponiveis(bot *tgbotapi.BotAPI, chatID int64) {
	var statuses []stringniveis(bot *tgbotapi.BotAPI, chatID int64) {
	db.DB.Model(&db.Item{}).Distinct("status").Pluck("status", &statuses)
	db.DB.Model(&db.Item{}).Distinct("status").Pluck("status", &statuses)
	var mensagem strings.Builder
	mensagem.WriteString("📊 **Status Disponíveis:**\n\n")
	mensagem.WriteString("📊 **Status Disponíveis:**\n\n")
	if len(statuses) == 0 {
		mensagem.WriteString("Nenhum status encontrado.\n")
	} else {m.WriteString("Nenhum status encontrado.\n")
		for i, status := range statuses {
			if status == "" {ange statuses {
				status = "Disponível"
			}status = "Disponível"
			mensagem.WriteString(fmt.Sprintf("**%d.** %s\n", i+1, status))
		}mensagem.WriteString(fmt.Sprintf("**%d.** %s\n", i+1, status))
	}}
	}
	mensagem.WriteString("\n💡 *Use:* `/listar_status <nome_do_status>`")
	mensagem.WriteString("\n💡 *Use:* `/listar_status <nome_do_status>`")
	msg := tgbotapi.NewMessage(chatID, mensagem.String())
	msg.ParseMode = "Markdown"(chatID, mensagem.String())
	bot.Send(msg) = "Markdown"
}bot.Send(msg)
}
// enviarMensagemErro envia mensagem de erro formatada
func enviarMensagemErro(bot *tgbotapi.BotAPI, chatID int64, erro string) {
	msg := tgbotapi.NewMessage(chatID, "❌ **Erro:** "+erro)64, erro string) {
	msg.ParseMode = "Markdown"(chatID, "❌ **Erro:** "+erro)
	bot.Send(msg) = "Markdown"
}bot.Send(msg)
}
























}    bot.Send(msg)    msg.ParseMode = "Markdown"    msg := tgbotapi.NewMessage(chatID, mensagem)        mensagem += "• `/listar_detalhado 123` - Mostra detalhes do item com ID 123.\n"    mensagem += "• `/listar_baixo_estoque 10` - Lista itens com estoque baixo (≤ 10).\n"    mensagem += "• `/listar_status disponível` - Lista todos os itens disponíveis.\n"    mensagem += "• `/listar pagina=2 limite=10` - Lista a página 2 com 10 itens por página.\n"    mensagem += "🔍 *Exemplos:*\n"    mensagem += "• Combine filtros e ordenações conforme necessário.\n\n"    mensagem += "• Utilize `limite=<número>` para definir a quantidade de itens por página (até 50).\n"    mensagem += "• Utilize `pagina=<número>` para navegar entre páginas de resultados.\n"    mensagem += "📄 *Dicas de Uso:*\n"    mensagem += "5. `/listar_detalhado` - Lista itens com todas as informações detalhadas.\n\n"    mensagem += "4. `/listar_baixo_estoque` - Lista itens com estoque abaixo de um limite definido.\n"    mensagem += "3. `/listar_status` - Lista itens filtrados por status. Ex: `/listar_status disponível`\n"    mensagem += "2. `/listar_resumo` - Lista itens com informações resumidas (nome e quantidade).\n"    mensagem += "1. `/listar` - Lista todos os itens cadastrados.\n"    mensagem += "Aqui estão os comandos disponíveis para listagem de itens:\n\n"    mensagem := "🛠️ **Ajuda - Comandos de Listagem**\n\n"func HandleAjudaListagem(bot *tgbotapi.BotAPI, chatID int64) {// HandleAjudaListagem agora aceita chatID diretamente.