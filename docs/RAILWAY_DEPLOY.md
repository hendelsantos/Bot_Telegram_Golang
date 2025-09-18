# Deploy no Railway

## 🚀 Configuração Rápida

### 1. Conectar Repositório
1. Acesse [Railway.app](https://railway.app)
2. Faça login com GitHub
3. Clique em "Deploy from GitHub repo"
4. Selecione `hendelsantos/Bot_Telegram_Golang`

### 2. Configurar Variáveis de Ambiente
No Railway Dashboard:
1. Vá em **Variables**
2. Adicione:
   ```
   TELEGRAM_BOT_TOKEN=seu_token_aqui
   DEBUG=false
   ```

### 3. Deploy Automático
- O Railway detectará automaticamente o projeto Go
- Usará o arquivo `nixpacks.toml` para build
- Build command: `go build -o bot_estoque ./cmd/main.go`
- Start command: `./bot_estoque`

## 🔧 Arquivos de Configuração

### `nixpacks.toml`
```toml
[phases.setup]
nixPkgs = ['go_1_21']

[phases.build]
cmds = ['go build -o bot_estoque ./cmd/main.go']

[start]
cmd = './bot_estoque'
```

### `railway.toml`
```toml
[build]
builder = "NIXPACKS"

[deploy]
startCommand = "./bot_estoque"
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10
```

## 📊 Monitoramento

### Logs
```bash
# Verificar logs no Railway Dashboard
# Ou usar Railway CLI:
railway logs
```

### Métricas
- CPU e Memória disponíveis no Dashboard
- Restart automático em caso de falha
- Health checks automáticos

## 🔒 Segurança

- ✅ Token armazenado como variável de ambiente
- ✅ Logs sem informações sensíveis
- ✅ Debug desabilitado em produção
- ✅ Database local (SQLite)

## 🚨 Troubleshooting

### Erro "no Go files in /app"
- ✅ Resolvido com `nixpacks.toml`
- ✅ Build command correta: `./cmd/main.go`

### Bot não responde
1. Verificar variável `TELEGRAM_BOT_TOKEN`
2. Verificar logs no Railway
3. Verificar se o bot está ativo no @BotFather

### Database Issues
- SQLite é criado automaticamente
- Dados persistem entre deploys
- Backup manual se necessário

## 📈 Performance

### Recursos Padrão
- 512MB RAM
- 1 vCPU
- Scaling automático disponível

### Otimizações
- Binary compilado com flags de otimização
- Imagem Alpine Linux (pequena)
- Restart policy configurada

## 💡 Dicas

1. **Primeiro Deploy**: Pode levar 2-3 minutos
2. **Updates**: Deploy automático a cada push no GitHub
3. **Downtime**: Zero downtime com Railway Pro
4. **Logs**: Mantenha `DEBUG=false` em produção
5. **Backup**: Considere backup periódico do SQLite
