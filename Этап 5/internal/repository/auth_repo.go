package repository

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"NazirovJasur/my-http-server.git/internal/models"
)

type Session struct {
	UserID    int64
	ExpiresAt time.Time
}

type AuthRepository struct {
	db       *sql.DB
	sessions map[string]Session
	mu       sync.RWMutex
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db:       db,
		sessions: make(map[string]Session),
	}
}

// Создание пользователя
func (r *AuthRepository) CreateUser(email, passwordHash string) error {
	_, err := r.db.Exec(`INSERT INTO users (email, password_hash) VALUES (?, ?)`, email, passwordHash)
	return err
}

// Получение пользователя по email (для входа)
func (r *AuthRepository) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(`SELECT id, email, password_hash FROM users WHERE email = ?`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

// --- УПРАВЛЕНИЕ СЕССИЯМИ ---

func (r *AuthRepository) CreateSession(sessionID string, userID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[sessionID] = Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Сессия живет 24 часа
	}
}

func (r *AuthRepository) GetUserBySession(sessionID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sess, ok := r.sessions[sessionID]
	if !ok {
		return 0, errors.New("сессия не найдена")
	}
	if time.Now().After(sess.ExpiresAt) {
		delete(r.sessions, sessionID) // Удаляем просроченную
		return 0, errors.New("сессия истекла")
	}
	return sess.UserID, nil
}

func (r *AuthRepository) DeleteSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, sessionID)
}
