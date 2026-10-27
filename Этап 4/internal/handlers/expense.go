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

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Получаем параметры фильтрации
	categoryID, _ := strconv.ParseInt(r.URL.Query().Get("category"), 10, 64)
	dateFrom := r.URL.Query().Get("from")
	dateTo := r.URL.Query().Get("to")

	items, err := h.repo.GetAll(categoryID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.repo.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Page":       "list",
		"Items":      items,
		"Categories": categories,
		"FilterCat":  categoryID,
		"FilterFrom": dateFrom,
		"FilterTo":   dateTo,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ExpenseHandler) Add(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.repo.GetAllCategories()

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
			"Page":       "add",
			"Categories": categories,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}
	dateStr := r.FormValue("date")
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Некорректный формат даты", http.StatusBadRequest)
		return
	}
	catID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	if catID == 0 {
		http.Error(w, "Выберите категорию", http.StatusBadRequest)
		return
	}

	e := models.Expense{
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
	categories, _ := h.repo.GetAllCategories()

	if r.Method == http.MethodGet {
		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		e, err := h.repo.GetByID(id)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
			"Page":       "edit",
			"Item":       e,
			"Categories": categories,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}
	dateStr := r.FormValue("date")
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Некорректный формат даты", http.StatusBadRequest)
		return
	}
	catID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	if catID == 0 {
		http.Error(w, "Выберите категорию", http.StatusBadRequest)
		return
	}

	e := models.Expense{
		ID:          id,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        dateStr,
		CategoryID:  catID,
	}
	if err := h.repo.Update(e); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
