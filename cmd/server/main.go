package main

import (
	"log"
	"net/http"

	"NazirovJasur/my-http-server.git/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Страницы (различать методы прямо в шаблоне пути)
	mux.HandleFunc("GET /", handlers.ListExpenses)
	mux.HandleFunc("GET /add", handlers.ShowAddForm)
	mux.HandleFunc("POST /add", handlers.AddExpense)

	// Статика: /static/style.css → web/static/style.css
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	addr := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
