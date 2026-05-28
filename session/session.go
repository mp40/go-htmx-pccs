package session

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt int64     `json:"created_at"` // unix timestamp
	ExpiresAt int64     `json:"expires_at"` // unix timestamp
}

type Store struct {
	db *sql.DB
}

func NewSessionService(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetSessionByID(ID uuid.UUID) (*Session, error) {
	session := Session{}
	var idStr string
	var userIdStr string
	err := s.db.QueryRow(
		"SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?", ID.String(),
	).Scan(&idStr, &userIdStr, &session.CreatedAt, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	session.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	session.UserID, err = uuid.Parse(userIdStr)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error) {
	ID := uuid.New()
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)",
		ID.String(), userID.String(), now, expiresAt.Unix(),
	)
	if err != nil {
		return uuid.Nil, err
	}
	return ID, nil
}
