package modules

// Centraliza todos os estados de conversa do usuário.
var userStateCadastro = make(map[int64]*UserStateCadastro)
var userStateAtualizacao = make(map[int64]*UserStateAtualizacao)
var userStateReparo = make(map[int64]*UserStateReparo)

// Funções para verificar em qual fluxo o usuário está.
func IsUserInCadastroFlow(chatID int64) bool {
    _, exists := userStateCadastro[chatID]
    return exists
}

func IsUserInAtualizacaoFlow(chatID int64) bool {
    _, exists := userStateAtualizacao[chatID]
    return exists
}

func IsUserInReparoFlow(chatID int64) bool {
    _, exists := userStateReparo[chatID]
    return exists
}

// Funções para limpar os estados de fluxo.
func ClearCadastroState(chatID int64) {
    delete(userStateCadastro, chatID)
}

func ClearAtualizacaoState(chatID int64) {
    delete(userStateAtualizacao, chatID)
}

func ClearReparoState(chatID int64) {
    delete(userStateReparo, chatID)
}

// Função genérica para limpar todos os estados de um usuário.
func ClearUserState(chatID int64) {
    ClearCadastroState(chatID)
    ClearAtualizacaoState(chatID)
    ClearReparoState(chatID)
}
