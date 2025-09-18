#!/bin/bash

# Script de configuração inicial do Bot de Estoque
# Executa a configuração completa do projeto

echo "🤖 CONFIGURAÇÃO INICIAL DO BOT DE ESTOQUE"
echo "=========================================="
echo ""

# Verifica se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado. Instale Go 1.21+ primeiro."
    echo "📥 Download: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go encontrado: $(go version)"
echo ""

# Configura dependências
echo "📦 Configurando dependências..."
go mod tidy
if [ $? -eq 0 ]; then
    echo "✅ Dependências configuradas com sucesso!"
else
    echo "❌ Erro ao configurar dependências."
    exit 1
fi
echo ""

# Verifica se existe .env
if [ ! -f ".env" ]; then
    echo "📝 Criando arquivo de configuração..."
    cp .env.example .env
    echo "✅ Arquivo .env criado a partir do exemplo."
    echo ""
    echo "⚠️  IMPORTANTE: Configure seu token no arquivo .env"
    echo "📝 Use: nano .env"
    echo "🤖 Obtenha seu token em: https://t.me/BotFather"
else
    echo "✅ Arquivo .env já existe."
fi
echo ""

# Compila o projeto
echo "🔨 Compilando o projeto..."
go build -o bot_estoque ./cmd/main.go
if [ $? -eq 0 ]; then
    echo "✅ Compilação bem-sucedida!"
    echo "📁 Executável criado: ./bot_estoque"
else
    echo "❌ Erro na compilação."
    exit 1
fi
echo ""

# Verifica token
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "⚠️  Token do Telegram não definido."
    echo ""
    echo "📋 PRÓXIMOS PASSOS:"
    echo "1. Configure seu token:"
    echo "   export TELEGRAM_BOT_TOKEN='seu_token_aqui'"
    echo ""
    echo "2. OU edite o arquivo .env:"
    echo "   nano .env"
    echo ""
    echo "3. Execute o bot:"
    echo "   ./bot_estoque"
else
    echo "✅ Token do Telegram configurado!"
    echo ""
    echo "🚀 PRONTO PARA USAR!"
    echo "Execute: ./bot_estoque"
fi
echo ""

echo "📚 COMANDOS ÚTEIS:"
echo "• ./scripts/advanced_reset.sh  - Reset da API"
echo "• ./scripts/demo_listagem.sh   - Demonstração"
echo "• ./scripts/start_bot.sh       - Inicialização completa"
echo ""

echo "📖 DOCUMENTAÇÃO:"
echo "• README.md                    - Guia completo"
echo "• docs/MODULO_LISTAGEM.md      - Documentação técnica"
echo ""

echo "✨ Configuração concluída com sucesso!"
