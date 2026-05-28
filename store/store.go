package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Hash      string    `json:"hash"`
	CreatedAt int64     `json:"created_at"` // unix timestamp
	UpdatedAt int64     `json:"updated_at"` // unix timestamp
}

type Store struct {
	db *sql.DB
}

func NewStoreService(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CountUserByEmail(email string) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return count, err
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	user := User{}
	var idStr string
	err := s.db.QueryRow("SELECT id, email, hash, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&idStr, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CountUserByID(ID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", ID.String()).Scan(&count)
	return count, err
}

func (s *Store) AddUser(email string, hash string) (uuid.UUID, error) {
	ID := uuid.New()
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO users (id, email, hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		ID.String(), email, hash, now, now,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return ID, nil
}
