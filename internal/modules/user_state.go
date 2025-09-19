package modules

// Estados de usuário para diferentes fluxos
var userStateGlobal = make(map[int64]string)

func IsUserInCadastroFlow(chatID int64) bool {
    state, exists := userStateGlobal[chatID]
    return exists && state == "cadastro"
}

func IsUserInAtualizacaoFlow(chatID int64) bool {
    state, exists := userStateGlobal[chatID]
    return exists && (state == "atualizacao_nome" || state == "atualizacao_descricao" || state == "atualizacao_quantidade" || state == "atualizacao_foto")
}

func ClearUserState(chatID int64) {
    delete(userStateGlobal, chatID)
}

func SetUserState(chatID int64, state string) {
    userStateGlobal[chatID] = state
}

func GetUserState(chatID int64) string {
    return userStateGlobal[chatID]
}
