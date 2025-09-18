#!/bin/bash

# Script para garantir que todas as instâncias do bot sejam encerradas corretamente

echo "Procurando por instâncias do bot..."

# Procura processos que estão usando a porta comum do Telegram API
for pid in $(lsof -i :443 | grep -v "LISTEN" | awk '{print $2}' | grep -v "PID"); do
  echo "Verificando processo $pid..."
  if ps -p $pid -o cmd= | grep -q "botgo"; then
    echo "Encontrou processo do bot: $pid"
    kill -9 $pid
    echo "Processo $pid encerrado"
  fi
done

# Procura por qualquer processo do bot diretamente
if pgrep -f "botgo" > /dev/null; then
  echo "Encontrou processos do bot, encerrando..."
  pkill -f "botgo"
  sleep 1
  
  # Verifica se algum processo persistiu
  if pgrep -f "botgo" > /dev/null; then
    echo "Alguns processos persistiram, usando SIGKILL..."
    pkill -9 -f "botgo"
  fi
  echo "Todos os processos do bot encerrados."
else
  echo "Nenhum processo do bot encontrado."
fi

# Limpa o cache do API Token, caso necessário
curl -s -X GET https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/deleteWebhook
curl -s -X GET https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/close

echo "Limpeza concluída. Agora você pode iniciar uma nova instância do bot."
