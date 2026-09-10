package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"tailor-management-cli/entity"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(name, email, password, phone string) error {
	query := `INSERT INTO users (name, email, password, role, phone) VALUES (?, ?, ?, 'customer', ?)`
	_, err := h.db.Exec(query, strings.TrimSpace(name), strings.TrimSpace(email), password, strings.TrimSpace(phone))
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("email sudah terdaftar, silakan gunakan email lain")
		}
		return fmt.Errorf("gagal mendaftarkan user: %w", err)
	}
	return nil
}

func (h *AuthHandler) SignIn(email, password string) (*entity.User, error) {
	email = strings.TrimSpace(email)
	var user entity.User

	query := `SELECT id, name, email, password, role, phone, created_at FROM users WHERE email = ? LIMIT 1`
	err := h.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Phone,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("email/password salah ")
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// Validasi kecocokan password
	if user.Password != password {
		return nil, errors.New("email/password salah")
	}

	return &user, nil
}
