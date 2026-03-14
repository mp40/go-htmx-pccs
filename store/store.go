package store

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Hash  string    `json:"hash"`
}

type Store struct {
	db *sql.DB
}

func NewStoreService(db *sql.DB) *Store {
	return &Store{db: db}
}

func (d *Store) CountUserByEmail(email string) (int, error) {
	return 0, nil
}

func (d *Store) GetUserByEmail(email string) (*User, error) {
	user := User{}
	return &user, nil
}

func (d *Store) CountUserByID(ID uuid.UUID) (int, error) {
	return 0, nil
}

func (d *Store) AddUser(email string, hash string) (id uuid.UUID, err error) {
	ID := uuid.New()
	return ID, nil
}
