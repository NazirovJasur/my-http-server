package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Добро пожаловать на учебный HTTP-сервер на Go!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Учебный проект: HTTP-сервер на стандартной библиотеке net/http.")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/ping", pingHandler)

	addr := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
