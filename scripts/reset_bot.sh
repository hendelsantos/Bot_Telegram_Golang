#!/bin/bash

# Script para resetar o bot pelo API do Telegram

echo "Resetando o bot pelo API do Telegram..."

# Token do bot
BOT_TOKEN="8320983858:AAG8ID2e7ReRk81YFEF3hJiqvD9HCDvLNlU"

# Chamadas à API
echo "Deletando webhook..."
curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/deleteWebhook?drop_pending_updates=true"

echo -e "\nFechando sessão atual..."
curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/close"

echo -e "\nEnviando comando getMe para verificar acesso..."
curl -s -X GET "https://api.telegram.org/bot$BOT_TOKEN/getMe"

echo -e "\n\nReset completo. Aguarde 10 segundos antes de iniciar o bot novamente."
echo "Matando processos existentes..."
pkill -f "bot_estoque" || true

# Aguarda 10 segundos
sleep 10

echo "Agora você pode iniciar o bot novamente."
