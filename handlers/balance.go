package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"test_project/utils"
)

// CheckBalance - обработчик запроса баланса
func CheckBalance(w http.ResponseWriter, r *http.Request) {
	// Получаем токен для авторизации
	accessToken, err := utils.GetAccessToken()
	if err != nil {
		log.Printf("Ошибка при получении токена: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем баланс
	balanceResponse, err := utils.FetchBalance(accessToken)
	if err != nil {
		log.Printf("Ошибка при получении баланса: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balanceResponse)
}
