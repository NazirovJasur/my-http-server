package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"NazirovJasur/my-http-server.git/internal/models"
)

// Хранилище в памяти
var expenses []models.Expense
var nextID = 1

// Парсим все шаблоны из папки web/templates
var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

// GET / — список трат
func ListExpenses(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := tmpl.ExecuteTemplate(w, "index.html", expenses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /add — показать форму
func ShowAddForm(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "add.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /add — обработать форму
func AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	dateStr := r.FormValue("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	expenses = append(expenses, models.Expense{
		ID:          nextID,
		Amount:      amount,
		Description: description,
		Date:        date,
	})
	nextID++

	// Post/Redirect/Get — защита от повторной отправки формы
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
