package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"test_project/models"

	"github.com/joho/godotenv"
)

// GetAccessToken - функция для получения токена авторизации
func GetAccessToken() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки .env файла: %w", err)
	}

	login := os.Getenv("login")
	password := os.Getenv("password")

	if login == "" || password == "" {
		return "", fmt.Errorf("переменные окружения CLIENT_ID и CLIENT_SECRET не установлены")
	}

	authRequest := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     login,
		"client_secret": password,
	}

	var buf bytes.Buffer
	for k, v := range authRequest {
		buf.WriteString(url.QueryEscape(k))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v))
		buf.WriteString("&")
	}
	body := bytes.TrimSuffix(buf.Bytes(), []byte("&"))

	resp, err := http.Post(
		"https://auth.bankingapi.ru/auth/realms/kubernetes/protocol/openid-connect/token",
		"application/x-www-form-urlencoded",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса токена: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка авторизации: %d, тело ответа: %s", resp.StatusCode, string(bodyBytes))
	}

	var authResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return authResponse.AccessToken, nil
}

// FetchBalance - функция для получения данных о балансе
func FetchBalance(accessToken string) (models.BalanceResponse, error) {
	// Настройка клиента с игнорированием сертификатов
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Отключаем проверку сертификата
		},
	}
	client := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest("GET", "https://api.bankingapi.ru/api/rb/pmnt/acceptance/universal/hackathon/v1/products", nil)
	if err != nil {
		return models.BalanceResponse{}, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.BalanceResponse{}, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.BalanceResponse{}, fmt.Errorf("ошибка при чтении ответа от API: %w", err)
	}

	var balanceResponse models.BalanceResponse
	if err := json.Unmarshal(body, &balanceResponse); err != nil {
		return models.BalanceResponse{}, fmt.Errorf("ошибка распаковки JSON: %w", err)
	}

	return balanceResponse, nil
}
