package modules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"botgo/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

// HandleListagem gerencia todos os comandos de listagem
func HandleListagem(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := strings.Fields(message.Text)
	comando := args[0]

	// Remove o comando da lista de argumentos
	var params []string
	if len(args) > 1 {
		params = args[1:]
	}

	switch comando {
	case "/listar":
		handleListarTodos(bot, message.Chat.ID, params)
	case "/listar_resumo":
		handleListarResumo(bot, message.Chat.ID, params)
	case "/listar_status":
		handleListarPorStatus(bot, message.Chat.ID, params)
	case "/listar_baixo_estoque":
		handleListarBaixoEstoque(bot, message.Chat.ID, params)
	case "/listar_detalhado":
		handleListarDetalhado(bot, message.Chat.ID, params)
	}
}

// handleListarTodos lista todos os itens com paginação
func handleListarTodos(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	config := ListarConfig{
		Limite:    10,
		Pagina:    1,
		Ordenacao: "id DESC",
		Filtros:   make(map[string]interface{}),
	}
	parseArgumentos(&config, args)

	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens: "+err.Error())
		return
	}
	if resultado.Total == 0 {
		enviarMensagemSimples(bot, chatID, "📭 Nenhum item cadastrado no estoque.")
		return
	}
	enviarListagemFormatada(bot, chatID, resultado, "📋 **Lista Completa do Estoque**")
}

// handleListarResumo lista apenas informações básicas
func handleListarResumo(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	config := ListarConfig{
		Limite:    20,
		Pagina:    1,
		Ordenacao: "nome ASC",
		Filtros:   make(map[string]interface{}),
	}
	parseArgumentos(&config, args)

	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens.")
		return
	}
	if resultado.Total == 0 {
		enviarMensagemSimples(bot, chatID, "📭 Estoque vazio.")
		return
	}
	enviarResumoFormatado(bot, chatID, resultado)
}

// handleListarPorStatus lista itens por status específico
func handleListarPorStatus(bot *tgbotapi.BotAPI, chatID int64, args []string) {
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
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens.")
		return
	}
	if resultado.Total == 0 {
		enviarMensagemSimples(bot, chatID, fmt.Sprintf("📭 Nenhum item encontrado com status: **%s**", status))
		return
	}
	titulo := fmt.Sprintf("📊 **Itens com Status: %s**", status)
	enviarListagemFormatada(bot, chatID, resultado, titulo)
}

// handleListarBaixoEstoque lista itens com quantidade baixa
func handleListarBaixoEstoque(bot *tgbotapi.BotAPI, chatID int64, args []string) {
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
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens.")
		return
	}
	if resultado.Total == 0 {
		enviarMensagemSimples(bot, chatID, fmt.Sprintf("✅ Não há itens com estoque baixo (≤ %d)", limite))
		return
	}
	titulo := fmt.Sprintf("⚠️ **Alerta: Estoque Baixo (≤ %d)**", limite)
	enviarListagemFormatada(bot, chatID, resultado, titulo)
}

// handleListarDetalhado lista itens com todas as informações
func handleListarDetalhado(bot *tgbotapi.BotAPI, chatID int64, args []string) {
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
	parseArgumentos(&config, args)

	resultado, err := buscarItensComPaginacao(config)
	if err != nil {
		enviarMensagemErro(bot, chatID, "Erro ao buscar itens.")
		return
	}
	if resultado.Total == 0 {
		enviarMensagemSimples(bot, chatID, "📭 Nenhum item cadastrado.")
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
	}}

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
	var b strings.Builder
	b.WriteString(titulo + "\n\n")

	for i, item := range resultado.Itens {
		b.WriteString(fmt.Sprintf("**%d.** `%s` (ID: %d)\n", (resultado.Pagina-1)*resultado.Limite+i+1, item.Nome, item.ID))
		b.WriteString(fmt.Sprintf("   📦 Qtd: **%d**  |  📋 Status: **%s**\n\n", item.Quantidade, item.Status))
	}

	b.WriteString(fmt.Sprintf("📄 **Página %d de %d** | **Total: %d itens**\n", resultado.Pagina, resultado.TotalPags, resultado.Total))
	enviarMensagemSimples(bot, chatID, b.String())
}

// enviarResumoFormatado envia apenas nome e quantidade
func enviarResumoFormatado(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	var b strings.Builder
	b.WriteString("📋 **Resumo do Estoque**\n\n")
	for i, item := range resultado.Itens {
		b.WriteString(fmt.Sprintf("**%d.** %s - **%d unid.**\n", (resultado.Pagina-1)*resultado.Limite+i+1, item.Nome, item.Quantidade))
	}
	b.WriteString(fmt.Sprintf("\n📊 **Total: %d itens cadastrados**", resultado.Total))
	enviarMensagemSimples(bot, chatID, b.String())
}

// enviarListagemDetalhada envia informações completas
func enviarListagemDetalhada(bot *tgbotapi.BotAPI, chatID int64, resultado *ResultadoListagem) {
	for _, item := range resultado.Itens {
		enviarItemDetalhado(bot, chatID, item.ID)
	}
	if len(resultado.Itens) > 1 {
		info := fmt.Sprintf("📄 **Página %d de %d** | **Total: %d itens**", resultado.Pagina, resultado.TotalPags, resultado.Total)
		enviarMensagemSimples(bot, chatID, info)
	}
}

// enviarItemDetalhado envia informações completas de um item específico
func enviarItemDetalhado(bot *tgbotapi.BotAPI, chatID int64, itemID uint) {
	var item db.Item
	if err := db.DB.First(&item, itemID).Error; err != nil {
		enviarMensagemErro(bot, chatID, fmt.Sprintf("Item com ID %d não encontrado.", itemID))
		return
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("🔍 **Detalhes do Item #%d**\n\n", item.ID))
	b.WriteString(fmt.Sprintf("📝 **Nome:** %s\n", item.Nome))
	b.WriteString(fmt.Sprintf("📦 **Quantidade:** %d unidades\n", item.Quantidade))
	b.WriteString(fmt.Sprintf("📋 **Status:** %s\n", item.Status))
	if item.Descricao != "" {
		b.WriteString(fmt.Sprintf("💬 **Descrição:** %s\n", item.Descricao))
	}
	if item.Fornecedor != "" {
		b.WriteString(fmt.Sprintf("🏢 **Fornecedor:** %s\n", item.Fornecedor))
	}
	if item.DataEnvio != nil {
		b.WriteString(fmt.Sprintf("📅 **Data Envio:** %s\n", item.DataEnvio.Format("02/01/2006")))
	}
	b.WriteString(fmt.Sprintf("\n⏰ **Criado em:** %s", item.CreatedAt.Format("02/01/2006 15:04")))
	enviarMensagemSimples(bot, chatID, b.String())
}

// enviarStatusDisponiveis lista os status disponíveis
func enviarStatusDisponiveis(bot *tgbotapi.BotAPI, chatID int64) {
	var statuses []string
	db.DB.Model(&db.Item{}).Distinct("status").Pluck("status", &statuses)
	var b strings.Builder
	b.WriteString("📊 **Status Disponíveis:**\n\n")
	if len(statuses) == 0 {
		b.WriteString("Nenhum status encontrado.\n")
	} else {
		for i, status := range statuses {
			if status != "" {
				b.WriteString(fmt.Sprintf("**%d.** %s\n", i+1, status))
			}
		}
	}
	b.WriteString("\n💡 *Use:* `/listar_status <nome_do_status>`")
	enviarMensagemSimples(bot, chatID, b.String())
}

func HandleAjudaListagem(bot *tgbotapi.BotAPI, chatID int64) {
	text := `🛠️ **Ajuda - Comandos de Listagem**

Aqui estão os comandos disponíveis para listagem de itens:

1.  **/listar** - Lista todos os itens cadastrados.
2.  **/listar_resumo** - Lista itens com informações resumidas.
3.  **/listar_status** - Lista itens filtrados por status.
4.  **/listar_baixo_estoque** - Lista itens com estoque baixo.
5.  **/listar_detalhado** - Lista itens com todas as informações.

📄 *Dicas de Uso:*
• Utilize ` + "`pagina=<número>`" + ` para navegar entre páginas.
• Utilize ` + "`limite=<número>`" + ` para definir a quantidade de itens por página.

🔍 *Exemplos:*
• ` + "`/listar pagina=2 limite=10`" + `
• ` + "`/listar_status disponível`" + `
• ` + "`/listar_baixo_estoque 10`" + `
• ` + "`/listar_detalhado 123`" + `
`
	enviarMensagemSimples(bot, chatID, text)
}

// Funções auxiliares
func enviarMensagemErro(bot *tgbotapi.BotAPI, chatID int64, erro string) {
	enviarMensagemSimples(bot, chatID, "❌ **Erro:** "+erro)
}

func enviarMensagemSimples(bot *tgbotapi.BotAPI, chatID int64, texto string) {
	msg := tgbotapi.NewMessage(chatID, texto)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func parseArgumentos(config *ListarConfig, args []string) {
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
		}
	}
}