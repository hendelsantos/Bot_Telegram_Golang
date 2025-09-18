#!/bin/bash

# Script para iniciar o bot de gestão de estoque
# Este script garante que não há outras instâncias rodando

echo "Iniciando bot de gestão de estoque..."

# Primeiro, encerra qualquer instância existente do bot
if pgrep -f "botgo" > /dev/null; then
  echo "Instâncias existentes encontradas, encerrando..."
  pkill -f "botgo"
  sleep 2
  
  # Verifica se alguma instância persistiu
  if pgrep -f "botgo" > /dev/null; then
    echo "Tentando encerramento forçado..."
    pkill -9 -f "botgo"
    sleep 1
  fi
fi

# Limpa o cache do API Token usando curl
if [ -n "$TELEGRAM_BOT_TOKEN" ]; then
  echo "Limpando cache da API Telegram..."
  curl -s -X GET "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/deleteWebhook"
  curl -s -X GET "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/close"
  sleep 1
else
  echo "❌ Erro: TELEGRAM_BOT_TOKEN não definida."
  echo "💡 Configure com: export TELEGRAM_BOT_TOKEN='seu_token_aqui'"
  exit 1
fi

echo "Compilando aplicação..."
cd /home/hendel/Documentos/BOTS/BOT_GO
go build -o bot_estoque cmd/main.go

echo "Iniciando aplicação..."
./bot_estoque
