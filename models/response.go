package models

// Account - структура счета или карты
type Account struct {
	ID      string `json:"id"`
	Balance struct {
		Amount   float64 `json:"amount"`
		Currency struct {
			Code string `json:"code"`
		} `json:"currency"`
	} `json:"balance"`
}

// BalanceResponse - ответ с балансами счета и карт
type BalanceResponse struct {
	Accounts []Account `json:"accounts"`
	Cards    []Account `json:"cards"`
}
