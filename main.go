package main

import (
	"fmt"
	"log"
	"net/http"
	"test_project/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	// Роуты
	r.HandleFunc("/", handlers.GlavHandler).Methods("GET")
	r.HandleFunc("/account/balance", handlers.CheckBalance).Methods("GET")

	// CORS настройка
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8000"}, // Разрешаем фронтенду с этого адреса
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Оборачиваем маршруты в CORS мидлварь
	handler := c.Handler(r)
	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
