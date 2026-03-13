package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Hash  string    `json:"hash"`
}

type Data struct {
	db *sql.DB
}

func NewDataService(db *sql.DB) *Data {
	return &Data{db: db}
}

func (d *Data) CountUserByEmail(email string) (int, error) {
	return 0, nil
}

func (d *Data) GetUserByEmail(email string) (*User, error) {
	user := User{}
	return &user, nil
}

func (d *Data) CountUserByID(ID uuid.UUID) (int, error) {
	return 0, nil
}

func (d *Data) AddUser(email string, hash string) (id uuid.UUID, err error) {
	ID := uuid.New()
	return ID, nil
}
