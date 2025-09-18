package modules

// IsUserInCadastroFlow verifica se um usuário está no fluxo de cadastro
func IsUserInCadastroFlow(userID int64) bool {
	_, exists := userStates[userID]
	return exists
}

// IsUserInAtualizacaoFlow verifica se um usuário está no fluxo de atualização
func IsUserInAtualizacaoFlow(userID int64) bool {
	_, exists := atualizaStates[userID]
	return exists
}

// IsUserInReparoFlow verifica se um usuário está no fluxo de reparo
func IsUserInReparoFlow(userID int64) bool {
	_, exists := reparoStates[userID]
	return exists
}
