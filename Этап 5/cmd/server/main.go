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

	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	// Обновленная схема БД с пользователями
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			email         TEXT    NOT NULL UNIQUE,
			password_hash TEXT    NOT NULL
		);

		CREATE TABLE IF NOT EXISTS categories (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT    NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS expenses (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id     INTEGER NOT NULL,
			amount      REAL    NOT NULL,
			description TEXT    NOT NULL,
			date        TEXT    NOT NULL,
			category_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);

		INSERT OR IGNORE INTO categories (id, name) VALUES 
		(1, 'Еда'),
		(2, 'Транспорт'),
		(3, 'Жильё'),
		(4, 'Развлечения'),
		(5, 'Другое');
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Репозитории
	expenseRepo := repository.NewExpenseRepository(db)
	authRepo := repository.NewAuthRepository(db) // Новый репозиторий

	// Обработчики
	expenseHandler := handlers.NewExpenseHandler(expenseRepo)
	authHandler := handlers.NewAuthHandler(authRepo)

	// Middleware авторизации
	authMiddleware := handlers.AuthMiddleware(authRepo)

	mux := http.NewServeMux()

	// Публичные маршруты (без авторизации)
	mux.HandleFunc("GET /register", authHandler.Register)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("GET /login", authHandler.Login)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", authHandler.Logout)

	// Защищенные маршруты (оборачиваем в middleware)
	mux.Handle("GET /", authMiddleware(http.HandlerFunc(expenseHandler.List)))
	mux.Handle("GET /add", authMiddleware(http.HandlerFunc(expenseHandler.Add)))
	mux.Handle("POST /add", authMiddleware(http.HandlerFunc(expenseHandler.Add)))
	mux.Handle("GET /edit", authMiddleware(http.HandlerFunc(expenseHandler.Edit)))
	mux.Handle("POST /edit", authMiddleware(http.HandlerFunc(expenseHandler.Edit)))
	mux.Handle("POST /delete", authMiddleware(http.HandlerFunc(expenseHandler.Delete)))

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	addr := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
