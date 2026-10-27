package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"NazirovJasur/my-http-server.git/internal/models"
	"NazirovJasur/my-http-server.git/internal/repository"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

type ExpenseHandler struct {
	repo *repository.ExpenseRepository
}

func NewExpenseHandler(repo *repository.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{repo: repo}
}

// Вспомогательная функция для получения userID из контекста
func getUserID(r *http.Request) int64 {
	return r.Context().Value(UserIDKey).(int64)
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	categoryID, _ := strconv.ParseInt(r.URL.Query().Get("category"), 10, 64)
	dateFrom := r.URL.Query().Get("from")
	dateTo := r.URL.Query().Get("to")

	items, err := h.repo.GetAll(userID, categoryID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, _ := h.repo.GetAllCategories()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Page":       "list",
		"Items":      items,
		"Categories": categories,
		"FilterCat":  categoryID,
		"FilterFrom": dateFrom,
		"FilterTo":   dateTo,
	})
}

func (h *ExpenseHandler) Add(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.repo.GetAllCategories()

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
			"Page":       "add",
			"Categories": categories,
		})
		return
	}

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	dateStr := r.FormValue("date")
	catID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)

	if amount <= 0 {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	e := models.Expense{
		UserID:      getUserID(r),
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        dateStr,
		CategoryID:  catID,
	}
	if err := h.repo.Create(e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ExpenseHandler) Edit(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	categories, _ := h.repo.GetAllCategories()

	if r.Method == http.MethodGet {
		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		e, err := h.repo.GetByID(userID, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
			"Page":       "edit",
			"Item":       e,
			"Categories": categories,
		})
		return
	}

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	dateStr := r.FormValue("date")
	catID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)

	e := models.Expense{
		ID:          id,
		UserID:      userID,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        dateStr,
		CategoryID:  catID,
	}
	if err := h.repo.Update(e); err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err := h.repo.Delete(getUserID(r), id); err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
