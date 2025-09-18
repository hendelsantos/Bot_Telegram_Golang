#!/bin/bash

# Script avançado de reset da API do Telegram
# Resolve a maioria dos problemas de conflito

# Verifica se o token está definido
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ Erro: TELEGRAM_BOT_TOKEN não definida"
    echo "💡 Configure com: export TELEGRAM_BOT_TOKEN='seu_token_aqui'"
    exit 1
fi

BOT_TOKEN="$TELEGRAM_BOT_TOKEN"

echo "🔧 RESET AVANÇADO DA API TELEGRAM"
echo "=================================="

# 1. Matar todos os processos do bot
echo "1. Encerrando todos os processos do bot..."
pkill -9 -f "bot_estoque" 2>/dev/null || true
pkill -9 -f "botgo" 2>/dev/null || true
pkill -9 -f "main.go" 2>/dev/null || true

# 2. Limpar conexões de rede relacionadas
echo "2. Limpando conexões de rede..."
lsof -ti:8080 | xargs kill -9 2>/dev/null || true

# 3. Reset completo da API
echo "3. Fazendo reset completo da API..."
curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/deleteWebhook?drop_pending_updates=true"
sleep 2
curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/close"
sleep 2

# 4. Forçar um getUpdates com offset alto para limpar fila
echo "4. Limpando fila de updates..."
curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getUpdates?offset=999999999&timeout=1"
sleep 1

# 5. Verificar se API está respondendo
echo "5. Verificando status da API..."
RESPONSE=$(curl -s "https://api.telegram.org/bot$BOT_TOKEN/getMe")
if echo "$RESPONSE" | grep -q '"ok":true'; then
    echo "✅ API está respondendo normalmente"
else
    echo "❌ Problema na API: $RESPONSE"
    exit 1
fi

# 6. Aguardar um tempo de segurança
echo "6. Aguardando 10 segundos para estabilização..."
sleep 10

echo ""
echo "✅ Reset concluído com sucesso!"
echo "💡 Agora você pode iniciar o bot com:"
echo "   cd /home/hendel/Documentos/BOTS/BOT_GO"
echo "   ./bot_estoque"
echo ""
echo "Se ainda houver conflitos, considere:"
echo "1. Aguardar mais 30 minutos"
echo "2. Criar um novo bot no @BotFather"
echo "3. Usar webhooks em vez de polling"
