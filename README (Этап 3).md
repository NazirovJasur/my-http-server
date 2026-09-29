# 💰 Учёт личных трат (этап 3)

Учебное веб-приложение на Go. Этап 3: постоянное хранение данных в SQLite через слой репозитория и полный CRUD трат.

## Возможности

- Список трат в таблице, сортировка от новых к старым.
- Добавление, редактирование и удаление трат (удаление с подтверждением).
- Данные хранятся в файле `expenses.db` и сохраняются между перезапусками.
- Post/Redirect/Get после отправки форм.
- Общий layout для всех страниц.

## Технологии

- Go 1.22 или новее
- `net/http`, `html/template`, `database/sql`
- SQLite, драйвер `modernc.org/sqlite` (чистый Go, без CGO)
- HTML и CSS без фреймворков

## Структура проекта

```
.
├── cmd/
│   └── server/
│       └── main.go                  # открытие БД, создание таблицы, маршруты
├── internal/
│   ├── handlers/
│   │   └── expense.go               # ExpenseHandler: List, Add, Edit, Delete (без SQL)
│   ├── models/
│   │   └── expense.go               # структура Expense
│   └── repository/
│       └── expense_repo.go          # ExpenseRepository: весь SQL-код
├── migrations/
│   └── 001_create_expenses.sql      # схема таблицы expenses
├── web/
│   ├── templates/
│   │   ├── layout.html
│   │   ├── list.html
│   │   ├── add.html
│   │   └── edit.html
│   └── static/
│       └── style.css
├── expenses.db                      # создаётся при первом запуске
├── go.mod
└── go.sum
```

## Запуск

Из корня проекта (где лежит `go.mod`):

```bash
go mod tidy
go run ./cmd/server
```

Откройте в браузере: <http://localhost:8080>

Чтобы начать с чистой базы, остановите сервер и удалите `expenses.db`. При следующем запуске таблица создастся заново.

## База данных

Таблица `expenses`:

| Поле          | Тип     | Описание                      |
|---------------|---------|-------------------------------|
| `id`          | INTEGER | первичный ключ, автоинкремент |
| `amount`      | REAL    | сумма                         |
| `description` | TEXT    | описание                      |
| `date`        | TEXT    | дата в формате `YYYY-MM-DD`   |

Схема описана в `migrations/001_create_expenses.sql`. Таблица создаётся при старте сервера запросом `CREATE TABLE IF NOT EXISTS`.

Все SQL-запросы находятся в `ExpenseRepository` и используют placeholder `?`, что защищает от SQL-инъекций. Методы: `Create`, `GetAll`, `GetByID`, `Update`, `Delete`.

## Маршруты

| Метод | Путь          | Описание                                 |
|-------|---------------|------------------------------------------|
| GET   | `/`           | Список трат                              |
| GET   | `/add`        | Форма добавления                         |
| POST  | `/add`        | Сохранение новой траты, redirect на `/`  |
| GET   | `/edit?id=N`  | Форма редактирования траты N             |
| POST  | `/edit`       | Сохранение изменений (id в скрытом поле) |
| POST  | `/delete`     | Удаление траты (id в форме)              |
| GET   | `/static/...` | Статические файлы (CSS)                  |

## Следующий этап

Категории трат и фильтрация по категории и датам (этап 4).

## Автор

NazirovJasur, группа Б-ИСиТ-22
