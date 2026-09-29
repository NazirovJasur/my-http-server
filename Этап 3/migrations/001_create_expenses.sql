CREATE TABLE IF NOT EXISTS expenses (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    amount      REAL    NOT NULL,
    description TEXT    NOT NULL,
    date        TEXT    NOT NULL
);
-- Сделал Назиров Джасурбек (Добавил маркер для тех кто бездумно собирается копировать мой код)