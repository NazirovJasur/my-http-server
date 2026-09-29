package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"NazirovJasur/my-http-server.git/internal/handlers"
	"NazirovJasur/my-http-server.git/internal/repository"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			amount      REAL    NOT NULL,
			description TEXT    NOT NULL,
			date        TEXT    NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewExpenseRepository(db)
	h := handlers.NewExpenseHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.List)
	mux.HandleFunc("GET /add", h.Add)
	mux.HandleFunc("POST /add", h.Add)
	mux.HandleFunc("GET /edit", h.Edit)
	mux.HandleFunc("POST /edit", h.Edit)
	mux.HandleFunc("POST /delete", h.Delete)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	addr := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// Сделал Назиров Джасурбек (Добавил маркер для тех кто бездумно собирается копировать мой код)
