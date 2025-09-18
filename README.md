# 🤖 Bot de Gestão de Estoque em Go

Um bot profissional para Telegram desenvolvido em Go para gerenciamento completo de estoque com funcionalidades avançadas.

## 🚀 Funcionalidades

### 📦 Gestão de Estoque
- ✅ Cadastro de itens com fotos
- ✅ Busca avançada e filtros
- ✅ Atualização de informações
- ✅ Controle de quantidade

### 📋 Listagem Profissional
- ✅ Listagem paginada
- ✅ Múltiplos formatos de visualização
- ✅ Filtros por status
- ✅ Alertas de estoque baixo
- ✅ Navegação inteligente

### 🔧 Controle de Reparos
- ✅ Registro de envio para reparo
- ✅ Controle de retorno
- ✅ Histórico de movimentações

### 📊 Relatórios
- ✅ Exportação para CSV
- ✅ Histórico detalhado
- ✅ Estatísticas de uso

## 🛠️ Tecnologias

- **Go 1.21+** - Linguagem principal
- **GORM** - ORM para banco de dados
- **SQLite** - Banco de dados local
- **Telegram Bot API** - Interface do bot
- **Markdown** - Formatação de mensagens

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Git
- Token do bot do Telegram (via @BotFather)

## ⚡ Instalação Rápida

### 1. Clone o repositório
```bash
git clone https://github.com/hendelsantos/Bot_Telegram_Golang.git
cd Bot_Telegram_Golang
```

### 2. Configure as dependências
```bash
go mod tidy
```

### 3. Configure o token
```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite o arquivo .env e adicione seu token
nano .env
```

### 4. Execute o bot
```bash
# Defina a variável de ambiente
export TELEGRAM_BOT_TOKEN="seu_token_aqui"

# Compile e execute
go build -o bot_estoque ./cmd/main.go
./bot_estoque
```

## 🎯 Comandos Disponíveis

### 📋 Comandos Básicos
| Comando | Descrição |
|---------|-----------|
| `/start` | Inicia o bot |
| `/menu` | Exibe menu principal |
| `/novoitem` | Cadastra novo item |
| `/buscar <termo>` | Busca itens |
| `/atualizar <ID>` | Atualiza item |

### 📊 Comandos de Listagem
| Comando | Descrição |
|---------|-----------|
| `/listar` | Lista todos os itens |
| `/listar_resumo` | Lista resumida |
| `/listar_detalhado` | Lista detalhada |
| `/listar_status <status>` | Lista por status |
| `/listar_baixo_estoque` | Alerta estoque baixo |

### 🔧 Comandos de Reparo
| Comando | Descrição |
|---------|-----------|
| `/enviar_reparo <ID>` | Envia para reparo |
| `/retornar_reparo <ID>` | Registra retorno |

### 📈 Relatórios
| Comando | Descrição |
|---------|-----------|
| `/exportar_estoque` | Exporta para CSV |
| `/historico <ID>` | Histórico do item |

## 🏗️ Estrutura do Projeto

```
├── cmd/
│   └── main.go              # Ponto de entrada
├── internal/
│   ├── bot/
│   │   └── bot.go           # Lógica do bot
│   ├── db/
│   │   ├── db.go            # Conexão do banco
│   │   ├── models.go        # Modelos de dados
│   │   └── historico.go     # Histórico
│   └── modules/
│       ├── cadastro.go      # Cadastro de itens
│       ├── consulta.go      # Consultas
│       ├── listagem.go      # Listagem avançada
│       ├── atualizacao.go   # Atualizações
│       ├── reparo.go        # Controle de reparos
│       ├── exportar.go      # Exportação
│       ├── historico.go     # Histórico
│       ├── menu.go          # Menu principal
│       └── user_state.go    # Estados do usuário
├── scripts/
│   ├── advanced_reset.sh    # Reset do bot
│   ├── start_bot.sh         # Inicialização
│   └── demo_listagem.sh     # Demonstração
├── docs/
│   └── MODULO_LISTAGEM.md   # Documentação
├── go.mod                   # Dependências Go
└── README.md               # Este arquivo
```

## 🔧 Configuração Avançada

### Variáveis de Ambiente

```bash
# Token do bot (obrigatório)
TELEGRAM_BOT_TOKEN=seu_token_aqui

# Configurações opcionais
DB_PATH=estoque.db
PORT=8080
USE_WEBHOOK=false
DEBUG=true
```

### Webhooks (Produção)

Para uso em produção, configure webhooks:

```bash
export USE_WEBHOOK=true
export WEBHOOK_URL=https://seudominio.com/webhook
export PORT=8080
```

## 📊 Exemplos de Uso

### Cadastro de Item
1. `/novoitem`
2. Envie uma foto
3. Digite o nome
4. Digite a descrição
5. Digite a quantidade
6. Confirme com "sim"

### Listagem Avançada
```bash
# Lista básica
/listar

# Lista com paginação
/listar pagina=2 limite=5

# Lista por status
/listar_status Disponível

# Alerta de estoque baixo
/listar_baixo_estoque 10
```

### Busca Inteligente
```bash
# Busca geral
/buscar motor

# Busca por status
/buscar status Em Reparo

# Busca por fornecedor
/buscar fornecedor ACME
```

## 🛡️ Segurança

- ✅ Tokens nunca commitados
- ✅ Validação de entrada
- ✅ Sanitização de dados
- ✅ Proteção contra SQL injection
- ✅ Rate limiting interno

## 📚 Documentação

- [Módulo de Listagem](docs/MODULO_LISTAGEM.md)
- [Scripts de Demonstração](scripts/demo_listagem.sh)

## 🚀 Deploy

### Docker (Recomendado)
```dockerfile
# Em breve...
```

### Systemd Service
```bash
# Em breve...
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Crie um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👨‍💻 Autor

**Hendel Santos**
- GitHub: [@hendelsantos](https://github.com/hendelsantos)

## 🙏 Agradecimentos

- Comunidade Go
- Telegram Bot API
- GORM ORM
- Contribuidores do projeto

---

⭐ **Se este projeto foi útil, considere dar uma estrela!**
